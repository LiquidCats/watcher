package main

import (
	"context"
	"os"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/bus"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/database"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo"
	"github.com/LiquidCats/watcher/v2/internal/adapter/state"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/usecase"
	"github.com/LiquidCats/watcher/v2/pkg/backbround"
	"github.com/LiquidCats/watcher/v2/pkg/graceful"
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
		if err := conn.Close(ctx); err != nil {
			logger.Fatal().Stack().Err(err).Msg("close connection")
		}
	}()
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("connect to database")
	}

	if err := database.Migrate(conn); err != nil {
		logger.Fatal().Stack().Err(err).Msg("migrate")
	}

	repository := database.New(conn)
	blocksState := state.NewPersisterState[entities.BlockHash](repository)
	transactionState := state.NewPersisterState[entities.TxID](repository)
	publisher := bus.NewRedisPublisher(cfg.Redis)

	runners := append(
		[]graceful.Runner{},
		graceful.Signals,
	)

	switch cfg.App.Type {
	case entities.TypeUtxo:
		client := utxo.NewClient(cfg.Utxo.Rpc)

		mempoolUseCase := usecase.NewMempoolProcessor(cfg.App, transactionState, client, publisher)
		blocksUseCase := usecase.NewBlocksProcessor(cfg.App, blocksState, client, publisher, publisher)

		runners = append(runners,
			backbround.BackgroundRunner("utxo.rpc.blocks", cfg.App.ScanInterval, blocksUseCase),
			backbround.BackgroundRunner("utxo.rpc.mempool", cfg.App.ScanInterval, mempoolUseCase),
		)
	case entities.TypeEvm:
		logger.Fatal().Msg("not implemented yet")
	}

	logger.Info().Msg("starting application")
	if err := graceful.WaitContext(ctx, runners...); err != nil {
		logger.Fatal().Err(err).Msg("application terminated")
	}

	logger.Info().Msg("shutting down")
}
