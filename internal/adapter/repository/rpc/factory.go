package rpc

import (
	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/evm"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/port/rpc"
	"github.com/rotisserie/eris"
)

func Factory(cfg configs.ChainConfig) (rpc.Client, error) {
	switch cfg.Type {
	case entities.TypeEvm:
		return evm.NewClient(cfg.RPC), nil
	case entities.TypeUtxo:
		return utxo.NewClient(cfg.RPC), nil
	}

	return nil, eris.Errorf("factory: unknown type: %v", cfg.Type)
}
