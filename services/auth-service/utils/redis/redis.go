package redis

import (
	"context"
	"fmt"

	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/logger"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewRedisClient),
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

// NewRedisClient creates a new Redis client.
func NewRedisClient(lc fx.Lifecycle, config *configs.Config, logger *logger.Logger) (*redis.Client, error) {
	cfg := Config{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := client.Ping(ctx).Result()
			if err != nil {
				return fmt.Errorf("failed to connect to Redis: %w", err)
			}
			logger.Info("Connected to Redis successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}
