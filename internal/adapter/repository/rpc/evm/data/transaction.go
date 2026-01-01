package data

import "github.com/LiquidCats/watcher/v2/internal/app/domain/entities"

type Transaction struct {
	BlockHash        entities.BlockHash `json:"blockHash"`
	BlockNumber      HexUint64          `json:"blockNumber"`
	From             entities.Address   `json:"from"`
	Gas              HexUint64          `json:"gas"`
	GasPrice         HexUint64          `json:"gasPrice"`
	Hash             entities.TxID      `json:"hash"`
	Input            string             `json:"input"`
	Nonce            HexUint64          `json:"nonce"`
	To               entities.Address   `json:"to"`
	TransactionIndex string             `json:"transactionIndex"`
	Value            HexUint64          `json:"value"`
	Type             HexUint64          `json:"type"`
	ChainID          HexUint64          `json:"chainId"`
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
