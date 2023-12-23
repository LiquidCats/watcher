package blockchain

import (
	"fmt"
	"watcher/configs"
	"watcher/internal/adapter/repository/rpc/bitcoin"
	"watcher/internal/adapter/repository/rpc/ethereum"
	"watcher/internal/app/domain/entity"
	"watcher/internal/port"
)

func GetBlockchainRpcRepository(cfg configs.Config) (port.RpcRepository, error) {
	switch cfg.Blockchain {
	case entity.Ethereum:
		return ethereum.NewRpcRepository(cfg), nil
	case entity.Bitcoin:
		return bitcoin.NewRpcRepository(cfg), nil
	}

	return nil, fmt.Errorf("unsupported blockchain")
}
