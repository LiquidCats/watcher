package data

import (
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
)

type Block struct {
	Hash              entities.BlockHash   `json:"hash"`
	Confirmations     int                  `json:"confirmations"`
	Height            entities.BlockHeight `json:"height"`
	Version           int                  `json:"version"`
	VersionHex        string               `json:"versionHex"`
	Merkleroot        string               `json:"merkleroot"`
	Time              int                  `json:"time"`
	Mediantime        int                  `json:"mediantime"`
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
