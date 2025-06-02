package main

import (
	"context"
	"os"

	"github.com/LiquidCats/graceful"
	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/bus/redis"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/database"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc"
	"github.com/LiquidCats/watcher/v2/internal/adapter/runner"
	"github.com/LiquidCats/watcher/v2/internal/adapter/state"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"

	_ "github.com/lib/pq"
	_ "go.uber.org/automaxprocs"
)

const app = "watcher"
const (
	TransactionChannelCap = 10000
	BlocksChannelCap      = 128
)

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

	zerolog.DefaultContextLogger = &logger // nolint:reassign
	zerolog.SetGlobalLevel(cfg.App.LogLevel)

	poolConfig, err := pgxpool.ParseConfig(cfg.DB.ToDSN())
	if err != nil {
		logger.Fatal().Err(err).Msg("parse db config")
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("database config")
	}
	defer pool.Close()
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("connect to database")
	}

	migrationConn, err := pool.Acquire(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("acquire pool connection")
	}

	if err = database.Migrate(migrationConn.Conn()); err != nil {
		logger.Fatal().Stack().Err(err).Msg("migrate")
	}

	redisClient := redis2.NewClient(cfg.Redis.ToConfig(app))

	dbRepository := database.New(pool)

	runners := []graceful.Runner{
		graceful.Signals,
	}

	for _, chainConfig := range cfg.Chains {
		blockChan := make(chan entities.Block, BlocksChannelCap)
		txIDChan := make(chan entities.TxID, TransactionChannelCap)

		blocksState := state.NewPersister[entities.BlockHash](chainConfig.Persist, dbRepository)
		blocksPublisher := redis.NewPublisher[entities.Block](redisClient)
		transactionsPublisher := redis.NewPublisher[entities.Transaction](redisClient)

		rpcRepository, chainErr := rpc.Factory(chainConfig)
		if chainErr != nil {
			logger.Fatal().Any("err", eris.ToString(err, true)).Msg("rpc adapter")
			return
		}

		oldMempool, mempoolErr := rpcRepository.GetMempool(ctx)
		if mempoolErr != nil {
			logger.Fatal().Any("err", eris.ToString(mempoolErr, true)).Msg("old mempool")
			return
		}

		blocksJob := usecase.NewBlocksJob(chainConfig, blocksState, rpcRepository, blockChan, blocksPublisher)
		mempoolJob := usecase.NewMempoolJob(chainConfig, rpcRepository, txIDChan, oldMempool)

		blockTicker := runner.NewTicker(chainConfig.Key("blocks"), chainConfig.Scan.Interval, blocksJob)
		mempoolTicker := runner.NewTicker(chainConfig.Key("mempool"), chainConfig.Scan.Interval, mempoolJob)

		txIDHandler := usecase.NewTxIDHandler(chainConfig, rpcRepository, transactionsPublisher)
		blockTransactionsHandler := usecase.NewBlockTransactionsHandler(chainConfig, rpcRepository, transactionsPublisher)

		txIDWorker := runner.NewWorker(
			chainConfig.Key("txid"),
			chainConfig.Workers.TxIDWorkerCount,
			txIDChan,
			txIDHandler,
		)
		blockTransactionWorker := runner.NewWorker(
			chainConfig.Key("block_transactions"),
			chainConfig.Workers.BlockTransactionsWorkerCount,
			blockChan,
			blockTransactionsHandler,
		)

		runners = append(
			runners,
			//
			blockTicker.Run,
			mempoolTicker.Run,
			//
			txIDWorker.Run,
			blockTransactionWorker.Run,
		)
	}

	logger.Info().Msg("starting application")
	if err = graceful.WaitContext(ctx, runners...); err != nil {
		logger.Fatal().Err(err).Msg("application terminated")
	}

	logger.Info().Msg("shutting down")
}
