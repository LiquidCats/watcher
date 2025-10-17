package usecase

import (
	"context"

	"github.com/LiquidCats/watcher/v2/configs"
	database2 "github.com/LiquidCats/watcher/v2/internal/adapter/repository/database"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/port/database"
	"github.com/LiquidCats/watcher/v2/internal/app/port/state"
	"github.com/bytedance/sonic"
	"github.com/rotisserie/eris"
)

type BlocksPersisterJob struct {
	cfg   configs.ChainConfig
	db    database.StateDB
	state state.State[entities.BlockHash]
}

func NewBlocksPersisterJob(
	cfg configs.ChainConfig,
	db database.StateDB,
	state state.State[entities.BlockHash],
) *BlocksPersisterJob {
	return &BlocksPersisterJob{
		cfg:   cfg,
		db:    db,
		state: state,
	}
}

func (uc *BlocksPersisterJob) Handle(ctx context.Context) error {
	blockHashes := uc.state.Get()
	blockHashesBytes, err := sonic.ConfigFastest.Marshal(blockHashes)
	if err != nil {
		return eris.Wrap(err, "marshal block hashes")
	}
	err = uc.db.SetState(ctx, database2.SetStateParams{
		Key:   uc.cfg.Key("blocks"),
		Value: blockHashesBytes,
	})
	if err != nil {
		return eris.Wrap(err, "set state")
	}

	return nil
}
