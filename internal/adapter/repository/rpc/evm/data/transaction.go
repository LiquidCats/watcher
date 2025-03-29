package data

import "github.com/LiquidCats/watcher/v2/internal/app/domain/entities"

type Transaction struct {
	BlockHash        entities.BlockHash `json:"blockHash"`
	BlockNumber      HexUint64          `json:"blockNumber"`
	From             string             `json:"from"`
	Gas              string             `json:"gas"`
	GasPrice         string             `json:"gasPrice"`
	Hash             entities.TxID      `json:"hash"`
	Input            string             `json:"input"`
	Nonce            string             `json:"nonce"`
	To               string             `json:"to"`
	TransactionIndex string             `json:"transactionIndex"`
	Value            string             `json:"value"`
	Type             string             `json:"type"`
	ChainID          string             `json:"chainId"`
	V                string             `json:"v"`
	R                string             `json:"r"`
	S                string             `json:"s"`
}

func (t Transaction) GetTxID() entities.TxID {
	return t.Hash
}

func (t Transaction) GetBlockHash() entities.BlockHash {
	return t.BlockHash
}
