package main

import (
	_ "go.uber.org/automaxprocs"

	"database/sql"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/pkg/graceful"
	_ "github.com/lib/pq"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

const appName = "watcher"

func main() {
	log := logger.NewLogger(appName)
	cfg, err := configs.Load(appName)
	if err != nil {
		log.Fatal("app: config", zap.Error(err))
	}

	conn, err := sql.Open(cfg.DB.Driver, cfg.DB.ToDSN())
	if err != nil {
		log.Fatal("app: database", zap.Error(err))
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			log.Fatal("app: database", zap.Error(err))
		}
	}(conn)

	rpc, err := blockchain.GetBlockchainRpcRepository(cfg)
	if err != nil {
		log.Fatal("app: rpc", zap.Error(err))
	}

	storage := postgres.NewStorageRepository(conn)
	publisher := redis.NewPublisher(appName, cfg)

	blockConfirmationUsecase := usecase.NewBlockConfirmationUsecase(cfg, rpc, storage, publisher, log.Named("confirmation"))
	blockHandlingUsecase := usecase.NewBlockHandlingUsecase(cfg, rpc, storage, publisher, log.Named("handling"))
	cleanupUsecase := usecase.NewCleanupUsecase(cfg, storage, log.Named("cleanup"))

	watcher := async.NewWatcher(log, cfg, blockConfirmationUsecase, blockHandlingUsecase, cleanupUsecase)

	if err := graceful.Wait(
		watcher.Watch,
		watcher.Confirm,
		watcher.Cleanup,
	); err != nil {
		log.Fatal("app: graceful", zap.Error(err))
	}

	log.Info("app: stopped")
}
