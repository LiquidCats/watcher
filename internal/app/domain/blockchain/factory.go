package blockchain

import (
	"fmt"
	"github.com/ybbus/jsonrpc/v3"
	"watcher/configs"
	"watcher/internal/adapter/repository/rpc/bitcoin"
	"watcher/internal/adapter/repository/rpc/ethereum"
	"watcher/internal/app/domain/entity"
	"watcher/internal/port"
)

func GetBlockchainRpcRepository(cfg configs.Config) (port.RpcRepository, error) {
	client := jsonrpc.NewClient(cfg.NodeUrl)

	switch cfg.Blockchain {
	case entity.Ethereum:
		return ethereum.NewRpcRepository(client), nil
	case entity.Bitcoin:
		return bitcoin.NewRpcRepository(client), nil
	}

	return nil, fmt.Errorf("unsupported blockchain")
}
