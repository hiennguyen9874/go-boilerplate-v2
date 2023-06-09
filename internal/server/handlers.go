package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/hibiken/asynq"
	_ "github.com/hiennguyen9874/go-boilerplate-v2/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/hiennguyen9874/go-boilerplate-v2/ent"
	httpSwagger "github.com/hiennguyen9874/go-boilerplate-v2/pkg/http-swagger"
	"github.com/redis/go-redis/v9"

	"github.com/hiennguyen9874/go-boilerplate-v2/config"
	authHttp "github.com/hiennguyen9874/go-boilerplate-v2/internal/auth/delivery/http"
	itemHttp "github.com/hiennguyen9874/go-boilerplate-v2/internal/items/delivery/http"
	itemRepository "github.com/hiennguyen9874/go-boilerplate-v2/internal/items/repository"
	itemUseCase "github.com/hiennguyen9874/go-boilerplate-v2/internal/items/usecase"
	apiMiddleware "github.com/hiennguyen9874/go-boilerplate-v2/internal/middleware"
	userHttp "github.com/hiennguyen9874/go-boilerplate-v2/internal/users/delivery/http"
	userDistributor "github.com/hiennguyen9874/go-boilerplate-v2/internal/users/distributor"
	userRepository "github.com/hiennguyen9874/go-boilerplate-v2/internal/users/repository"
	userUseCase "github.com/hiennguyen9874/go-boilerplate-v2/internal/users/usecase"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/logger"
)

// @title Go boilerplate
// @version 1.0

// @BasePath /api
// @securitydefinitions.oauth2.password	OAuth2Password
// @tokenUrl /api/auth/login
func New(client *ent.Client, redisClient *redis.Client, taskRedisClient *asynq.Client, cfg *config.Config, logger logger.Logger) (*chi.Mux, error) {
	r := chi.NewRouter()

	// Repository
	userPgRepo := userRepository.CreateUserPgRepository(client)
	userRedisRepo := userRepository.CreateUserRedisRepository(redisClient)
	itemPgRepo := itemRepository.CreateItemPgRepository(client)

	// Distributor
	userRedisTaskDistributor := userDistributor.NewUserRedisTaskDistributor(taskRedisClient, cfg, logger)

	// UseCase
	userUC := userUseCase.CreateUserUseCase(userPgRepo, userRedisRepo, userRedisTaskDistributor, cfg, logger)
	itemUC := itemUseCase.CreateItemUseCase(itemPgRepo, cfg, logger)

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
