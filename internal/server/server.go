package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/hiennguyen9874/go-boilerplate-v2/config"
	"github.com/hiennguyen9874/go-boilerplate-v2/ent"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/logger"
	"github.com/redis/go-redis/v9"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

// Server provides an http.Server.
type Server struct {
	server *http.Server
	cfg    *config.Config
	client *ent.Client
	logger logger.Logger
}

// NewServer creates and configures an APIServer serving all application routes.
func NewServer(cfg *config.Config, client *ent.Client, redisClient *redis.Client, taskRedisClient *asynq.Client, logger logger.Logger) (*Server, error) {
	logger.Info("configuring server...")

	api, err := New(client, redisClient, taskRedisClient, cfg, logger)
	if err != nil {
		return nil, err
	}

	var addr string
	if strings.Contains(cfg.Server.Port, ":") {
		addr = cfg.Server.Port
	} else {
		addr = ":" + cfg.Server.Port
	}

	return &Server{
		server: &http.Server{
			Addr:           addr,
			Handler:        api,
			ReadTimeout:    time.Second * time.Duration(cfg.Server.ReadTimeout),
			WriteTimeout:   time.Second * time.Duration(cfg.Server.WriteTimeout),
			MaxHeaderBytes: maxHeaderBytes,
		},
		cfg:    cfg,
		client: client,
		logger: logger,
	}, nil
}

// Start runs ListenAndServe on the http.Server with graceful shutdown.
func (srv *Server) Start() {
	srv.logger.Info("starting server...")

	go func() {
		srv.logger.Infof("Listening on %s\n", srv.server.Addr)
		if err := srv.server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit

	srv.logger.Infof("Shutting down server... Reason: %s", sig)

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	if err := srv.server.Shutdown(ctx); err != nil {
		panic(err)
	}
	srv.logger.Info("Server gracefully stopped")
}
