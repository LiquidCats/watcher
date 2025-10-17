package main

import (
	"context"
	"os"

	"github.com/LiquidCats/graceful"
	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/bus/redis"
	"github.com/LiquidCats/watcher/v2/internal/adapter/metrics/prometheus"
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

func main() { //nolint:funlen
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

	requestsToNodeMetric := prometheus.NewRequestsToNodeCount(app)

	redisClient := redis2.NewClient(cfg.Redis.ToConfig(app))

	dbRepository := database.New(pool)

	runners := []graceful.Runner{
		graceful.Signals,
		graceful.ServerRunner(prometheus.GinHandler(), cfg.Metrics),
	}

	for _, chainConfig := range cfg.Chains {
		blockChan := make(chan entities.Block, chainConfig.Workers.BlockTransactionsWorkerCount)
		txIDChan := make(chan entities.TxID, chainConfig.Workers.TxIDWorkerCount)

		blocksState := state.NewMemoryState[entities.BlockHash](chainConfig.Persist)
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

		blocksJob := usecase.NewBlocksJob(
			chainConfig,
			blocksState,
			rpcRepository,
			blockChan,
			blocksPublisher,
			usecase.BlocksJobMetrics{RequestToNodeCounter: requestsToNodeMetric},
		)
		mempoolJob := usecase.NewMempoolJob(
			chainConfig,
			rpcRepository,
			txIDChan,
			oldMempool,
			usecase.MempoolJobMetrics{RequestToNodeCounter: requestsToNodeMetric},
		)
		statePersisterJob := usecase.NewBlocksPersisterJob(
			chainConfig,
			dbRepository,
			blocksState,
		)

		blockTicker := runner.NewTicker(chainConfig.Key("blocks"), chainConfig.Scan.Interval, blocksJob)
		mempoolTicker := runner.NewTicker(chainConfig.Key("mempool"), chainConfig.Scan.Interval, mempoolJob)
		statePersisterTicker := runner.NewTicker(chainConfig.Key("persister"), chainConfig.Scan.Interval, statePersisterJob)

		txIDHandler := usecase.NewTxIDHandler(
			chainConfig,
			rpcRepository,
			transactionsPublisher,
			usecase.TxIDHandlerMetrics{RequestToNodeCounter: requestsToNodeMetric},
		)
		blockTransactionsHandler := usecase.NewBlockTransactionsHandler(
			chainConfig,
			rpcRepository,
			transactionsPublisher,
			usecase.BlockTransactionsHandlerMetrics{RequestToNodeCounter: requestsToNodeMetric},
		)

		txIDWorker := runner.NewWorker[entities.TxID](
			chainConfig.Key("txid"),
			chainConfig.Workers.TxIDWorkerCount,
			txIDChan,
			txIDHandler,
		)
		blockTransactionWorker := runner.NewWorker[entities.Block](
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
			statePersisterTicker.Run,
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
