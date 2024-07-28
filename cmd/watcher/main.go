package main

import (
	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/event/redis"
	"github.com/LiquidCats/watcher/v2/internal/adapter/logger"
	"go.uber.org/zap"
)

const appName = "watcher"

func main() {
	log := logger.NewLogger(appName)
	cfg, err := configs.Load(appName)
	if err != nil {
		log.Fatal("app: config", zap.Error(err))
	}

	_ = redis.NewPublisher(appName, cfg)

	log.Info("app: stopped")
}
