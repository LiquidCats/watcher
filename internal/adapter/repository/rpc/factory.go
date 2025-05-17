package rpc

import (
	"fmt"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/evm"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/port/rpc"
)

func Factory(t entities.Type, cfg configs.Config) (rpc.Client, error) {
	switch t {
	case entities.TypeEvm:
		return evm.NewClient(cfg.Evm.RPC), nil
	case entities.TypeUtxo:
		return utxo.NewClient(cfg.Utxo.RPC), nil
	}

	return nil, fmt.Errorf("factory: unknown type: %v", t)
}
