package bitcoin

import (
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entity"
	"github.com/btcsuite/btcd/chaincfg"
)

func GetPrams(network entity.Network) *chaincfg.Params {
	switch network {
	case entity.NetworkMainNet:
		return &chaincfg.MainNetParams
	case entity.NetworkTestNet:
		return &chaincfg.TestNet3Params
	case entity.NetworkRegTest:
		return &chaincfg.RegressionNetParams
	}

	return nil
}
