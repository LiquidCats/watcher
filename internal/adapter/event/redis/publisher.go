package redis

import (
	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/redis/go-redis/v9"
)

type Publisher struct {
	appName  string
	cfg      configs.Config
	producer *redis.Client
}

func NewPublisher(appName string, cfg configs.Config) *Publisher {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return &Publisher{
		appName:  appName,
		cfg:      cfg,
		producer: rdb,
	}
}
