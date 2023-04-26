package distributor

import (
	"github.com/hibiken/asynq"
	"github.com/hiennguyen9874/go-boilerplate/config"
	"github.com/hiennguyen9874/go-boilerplate/pkg/logger"
)

type RedisTaskDistributor struct {
	RedisClient *asynq.Client
	Cfg         *config.Config
	Logger      logger.Logger
}

func NewRedisClient(cfg *config.Config) *asynq.Client {
	redisOpt := asynq.RedisClientOpt{
		Addr: cfg.TaskRedis.Addr,
		DB:   cfg.TaskRedis.Db,
	}

	return asynq.NewClient(redisOpt)
}

func NewRedisTaskDistributor(redisClient *asynq.Client, cfg *config.Config, loggger logger.Logger) RedisTaskDistributor {
	return RedisTaskDistributor{
		RedisClient: redisClient,
		Cfg:         cfg,
		Logger:      loggger,
	}
}
