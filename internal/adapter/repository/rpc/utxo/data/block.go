package data

import (
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
)

type Block struct {
	Hash              entities.BlockHash   `json:"hash"`
	Height            entities.BlockHeight `json:"height"`
	PreviousBlockHash entities.BlockHash   `json:"previousblockhash"`
	Tx                []*Transaction       `json:"tx"`
}

func (b *Block) GetHeight() entities.BlockHeight {
	return b.Height
}

func (b *Block) GetHash() entities.BlockHash {
	return b.Hash
}

func (b *Block) GetPrevHash() entities.BlockHash {
	return b.PreviousBlockHash
}

func (b *Block) GetTransactions() []entities.Transaction {
	txs := make([]entities.Transaction, len(b.Tx))
	for i, tx := range b.Tx {
		txs[i] = tx
	}

	return txs
}
