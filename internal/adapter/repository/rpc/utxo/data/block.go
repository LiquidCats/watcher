package data

import (
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	kernel "github.com/LiquidCats/watcher/v2/internal/app/utxo/domain/entities"
)

type Block struct {
	Hash              entities.BlockHash `json:"hash"`
	Height            uint32             `json:"height"`
	PreviousBlockHash entities.BlockHash `json:"previousblockhash"`
	Tx                []*Transaction     `json:"tx"`
}

func (b *Block) ToEntity() *kernel.UtxoBlock {
	ent := &kernel.UtxoBlock{
		Hash:         b.Hash,
		PrevHash:     b.PreviousBlockHash,
		Height:       entities.BlockHeight(b.Height),
		Transactions: make([]*kernel.UtxoTransaction, len(b.Tx)),
	}

	for i, tx := range b.Tx {
		ent.Transactions[i] = tx.ToEntity()
	}

	return ent
}
