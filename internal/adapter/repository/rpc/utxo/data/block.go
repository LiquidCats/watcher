package data

import (
	"github.com/LiquidCats/watcher/v2/internal/app/kernel/domain/entities"
)

type Block struct {
	Hash              entities.BlockHash `json:"hash"`
	Height            uint32             `json:"height"`
	PreviousBlockHash entities.BlockHash `json:"previousblockhash"`
	Tx                []*Transaction     `json:"tx"`
}
