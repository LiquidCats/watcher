package data

import (
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
)

type Block[T any] struct {
	Hash              entities.BlockHash   `json:"hash"`
	Confirmations     int                  `json:"confirmations"`
	Height            entities.BlockHeight `json:"height"`
	Version           int                  `json:"version"`
	VersionHex        string               `json:"versionHex"`
	MerkleRoot        string               `json:"merkleroot"`
	Time              int                  `json:"time"`
	MedianTime        int                  `json:"mediantime"`
	Nonce             int                  `json:"nonce"`
	Bits              string               `json:"bits"`
	Target            string               `json:"target"`
	Difficulty        float64              `json:"difficulty"`
	Chainwork         string               `json:"chainwork"`
	NTx               int                  `json:"nTx"`
	PreviousBlockHash entities.BlockHash   `json:"previousblockhash"`
	NextBlockHash     entities.BlockHash   `json:"nextblockhash"`
	StrippedSize      int                  `json:"strippedsize"`
	Size              int                  `json:"size"`
	Weight            int                  `json:"weight"`
	Tx                []T                  `json:"tx"`
}

func (b *Block[T]) GetHeight() entities.BlockHeight {
	return b.Height
}

func (b *Block[T]) GetHash() entities.BlockHash {
	return b.Hash
}

func (b *Block[T]) GetPrevHash() entities.BlockHash {
	return b.PreviousBlockHash
}

func (b *Block[T]) GetTransactions() []entities.Transaction {
	txs := make([]entities.Transaction, len(b.Tx))
	for i, tx := range b.Tx {
		transac, ok := any(tx).(*Transaction)
		if !ok {
			continue
		}
		txs[i] = transac
	}

	return txs
}
