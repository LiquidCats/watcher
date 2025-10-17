package usecase_test

import (
	"testing"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/database"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/usecase"
	"github.com/LiquidCats/watcher/v2/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestBlocksPersisterJob_Handle(t *testing.T) {
	cfg := configs.ChainConfig{
		Driver:  "rpc",
		Type:    "evm",
		Chain:   "ethereum",
		Persist: configs.PersistConfig{},
		Scan:    configs.ScanConfig{},
		Workers: configs.WorkersConfig{},
		RPC:     configs.RPCConfig{},
		Topics:  configs.TopicsConfig{},
	}

	db := mocks.NewMockStateDB(t)
	db.On("SetState", t.Context(), database.SetStateParams{
		Key:   "evm.rpc.ethereum.blocks",
		Value: []byte(`["test_value"]`),
	}).Return(nil)

	st := mocks.NewMockState[entities.BlockHash](t)
	st.On("Get").Return([]entities.BlockHash{"test_value"})

	uc := usecase.NewBlocksPersisterJob(cfg, db, st)

	err := uc.Handle(t.Context())
	require.NoError(t, err)
}
