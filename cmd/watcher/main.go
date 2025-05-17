package main

import (
	"context"
	"os"

	"github.com/LiquidCats/graceful"
	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/bus"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/database"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc"
	"github.com/LiquidCats/watcher/v2/internal/adapter/runner"
	"github.com/LiquidCats/watcher/v2/internal/adapter/state"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/usecase"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	_ "github.com/lib/pq"
	_ "go.uber.org/automaxprocs"
)

const app = "watcher"

func main() {
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()

	ctx := logger.WithContext(context.Background())

	cfg, err := configs.Load(app)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("failed to load config")
	}

	conn, err := pgx.Connect(ctx, cfg.DB.ToDSN())
	defer func() {
		if err = conn.Close(ctx); err != nil {
			logger.Fatal().Stack().Err(err).Msg("close connection")
		}
	}()
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("connect to database")
	}

	if err = database.Migrate(conn); err != nil {
		logger.Fatal().Stack().Err(err).Msg("migrate")
	}

	rpcAdapter, err := rpc.Factory(cfg.App.Type, cfg)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("factory")
	}
	databaseAdapter := database.New(conn)
	publisherAdapter := bus.NewRedisPublisher(cfg.Redis)

	blocksState := state.NewPersister[entities.BlockHash](databaseAdapter)
	transactionState := state.NewPersister[entities.TxID](databaseAdapter)

	transactionChan := make(chan entities.Transaction, 10000)
	defer close(transactionChan)

	blockChan := make(chan entities.Block, 1000)
	defer close(blockChan)

	txIdHashChan := make(chan entities.TxID, 10000)
	defer close(txIdHashChan)

	blocksUseCase := usecase.NewBlocksJob(cfg.App, blocksState, rpcAdapter, blockChan)
	mempoolUseCase := usecase.NewMempoolJob(cfg.App, transactionState, rpcAdapter, txIdHashChan)

	txIDHandler := usecase.NewTxIDHandler(rpcAdapter, transactionChan)
	blockHandler := usecase.NewBlockHandler(cfg.App, publisherAdapter, blocksState, transactionChan)
	transactionHandler := usecase.NewTransactionHandler(publisherAdapter)

	txIDWorker := runner.NewWorker(cfg.App.TxIDWorkerCount, txIdHashChan, txIDHandler)
	transactionWorker := runner.NewWorker(cfg.App.TransactionWorkerCount, transactionChan, transactionHandler)
	blockWorker := runner.NewWorker(cfg.App.BlockWorkerCount, blockChan, blockHandler)

	blockProcessor := runner.NewProcessor(cfg.App.Key("blocks"), cfg.App.ScanInterval, blocksUseCase)
	mempoolProcessor := runner.NewProcessor(cfg.App.Key("mempool"), cfg.App.ScanInterval, mempoolUseCase)

	runners := []graceful.Runner{
		graceful.Signals,
		txIDWorker.Run,
		transactionWorker.Run,
		blockWorker.Run,
		blockProcessor.Run,
		mempoolProcessor.Run,
	}

	logger.Info().Msg("starting application")
	if err = graceful.WaitContext(ctx, runners...); err != nil {
		logger.Fatal().Err(err).Msg("application terminated")
	}

	logger.Info().Msg("shutting down")
}
