package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"watcher/configs"
	"watcher/database"
	"watcher/internal/adapter/event/null"
	"watcher/internal/adapter/logger"
	"watcher/internal/adapter/process/async"
	"watcher/internal/adapter/repository/storage/postgres"
	"watcher/internal/app/domain/blockchain"
	"watcher/internal/app/usecase"
	"watcher/pkg/graceful"
)

const AppName = "watcher"

func main() {
	log := logger.NewLogger(AppName)
	cfg, err := configs.Load(AppName)
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

	if err := database.Migrate(conn, cfg); err != nil {
		log.Fatal("app: migrate", zap.Error(err))
	}

	rpc, err := blockchain.GetBlockchainRpcRepository(cfg)
	if err != nil {
		log.Fatal("app: rpc", zap.Error(err))
	}

	storage := postgres.NewStorageRepository(conn)
	publisher := null.NewPublisher(cfg)

	blockConfirmationUsecase := usecase.NewBlockConfirmationUsecase(cfg, rpc, storage, publisher, log)
	tickUsecase := usecase.NewBlockHandlingUsecase(cfg, rpc, storage, publisher, log)

	watcher := async.NewWatcher(log, cfg, blockConfirmationUsecase, tickUsecase)

	if err := graceful.Wait(
		watcher.Watch,
		watcher.Confirm,
	); err != nil {
		log.Fatal("app: graceful", zap.Error(err))
	}

	log.Info("app: stopped")
}
