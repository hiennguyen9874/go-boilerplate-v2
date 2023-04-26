package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/hibiken/asynq"
	_ "github.com/hiennguyen9874/go-boilerplate/docs" // docs is generated by Swag CLI, you have to import it.
	httpSwagger "github.com/hiennguyen9874/go-boilerplate/pkg/http-swagger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/hiennguyen9874/go-boilerplate/config"
	authHttp "github.com/hiennguyen9874/go-boilerplate/internal/auth/delivery/http"
	itemHttp "github.com/hiennguyen9874/go-boilerplate/internal/items/delivery/http"
	itemRepository "github.com/hiennguyen9874/go-boilerplate/internal/items/repository"
	itemUseCase "github.com/hiennguyen9874/go-boilerplate/internal/items/usecase"
	apiMiddleware "github.com/hiennguyen9874/go-boilerplate/internal/middleware"
	userHttp "github.com/hiennguyen9874/go-boilerplate/internal/users/delivery/http"
	userDistributor "github.com/hiennguyen9874/go-boilerplate/internal/users/distributor"
	userRepository "github.com/hiennguyen9874/go-boilerplate/internal/users/repository"
	userUseCase "github.com/hiennguyen9874/go-boilerplate/internal/users/usecase"
	"github.com/hiennguyen9874/go-boilerplate/pkg/logger"
)

// @title Go boilerplate
// @version 1.0

// @BasePath /api
// @securitydefinitions.oauth2.password	OAuth2Password
// @tokenUrl /api/auth/login
func New(db *gorm.DB, redisClient *redis.Client, taskRedisClient *asynq.Client, cfg *config.Config, logger logger.Logger) (*chi.Mux, error) {
	r := chi.NewRouter()

	// Repository
	userPgRepo := userRepository.CreateUserPgRepository(db)
	userRedisRepo := userRepository.CreateUserRedisRepository(redisClient)
	itemPgRepo := itemRepository.CreateItemPgRepository(db)

	// Distributor
	userRedisTaskDistributor := userDistributor.NewUserRedisTaskDistributor(taskRedisClient, cfg, logger)

	// UseCase
	userUC := userUseCase.CreateUserUseCaseI(userPgRepo, userRedisRepo, userRedisTaskDistributor, cfg, logger)
	itemUC := itemUseCase.CreateItemUseCaseI(itemPgRepo, cfg, logger)

	// Handler
	userHandler := userHttp.CreateUserHandler(userUC, cfg, logger)
	authHandler := authHttp.CreateAuthHandler(userUC, cfg, logger)
	itemHandler := itemHttp.CreateItemHandler(itemUC, cfg, logger)

	// middleware
	mw := apiMiddleware.CreateMiddlewareManager(cfg, logger, userUC)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(time.Duration(cfg.Server.ProcessTimeout) * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.Handler(mw.Cors()))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
	))

	apiRouter := chi.NewRouter()
	r.Mount("/api", apiRouter)

	apiRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, "pong")
	})

	authHttp.MapAuthRoute(apiRouter, authHandler, mw)
	userHttp.MapUserRoute(apiRouter, userHandler, mw)
	itemHttp.MapItemRoute(apiRouter, itemHandler, mw)

	return r, nil
}