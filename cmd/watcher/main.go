package main

import (
	"context"
	"os"

	"github.com/LiquidCats/graceful"
	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/bus"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/database"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/evm"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo"
	"github.com/LiquidCats/watcher/v2/internal/adapter/state"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/usecase"
	"github.com/LiquidCats/watcher/v2/pkg/backbround"
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

	repository := database.New(conn)
	blocksState := state.NewPersisterState[entities.BlockHash](repository)
	transactionState := state.NewPersisterState[entities.TxID](repository)
	publisher := bus.NewRedisPublisher(cfg.Redis)

	runners := append(
		[]graceful.Runner{},
		graceful.Signals,
	)

	var blocksUseCase *usecase.BlocksProcessor
	var mempoolUseCase *usecase.MempoolProcessor

	switch cfg.App.Type {
	case entities.TypeUtxo:
		client := utxo.NewClient(cfg.Utxo.RPC)

		blocksUseCase = usecase.NewBlocksProcessor(cfg.App, blocksState, client, publisher, publisher)
		mempoolUseCase = usecase.NewMempoolProcessor(cfg.App, transactionState, client, publisher)
	case entities.TypeEvm:
		client := evm.NewClient(cfg.Evm.RPC)

		blocksUseCase = usecase.NewBlocksProcessor(cfg.App, blocksState, client, publisher, publisher)
	}

	if blocksUseCase != nil {
		runners = append(
			runners,
			backbround.BackgroundRunner(
				cfg.App.Key("blocks"),
				cfg.App.ScanInterval,
				blocksUseCase,
			),
		)
	}

	if mempoolUseCase != nil {
		runners = append(
			runners,
			backbround.BackgroundRunner(
				cfg.App.Key("mempool"),
				cfg.App.ScanInterval,
				mempoolUseCase,
			),
		)
	}

	logger.Info().Msg("starting application")
	if err = graceful.WaitContext(ctx, runners...); err != nil {
		logger.Fatal().Err(err).Msg("application terminated")
	}

	logger.Info().Msg("shutting down")
}
