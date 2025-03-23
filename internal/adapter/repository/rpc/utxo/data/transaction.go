package data

import (
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	TxID          entities.TxID      `json:"txid"`
	Vin           []TransactionVin   `json:"vin"`
	Vout          []TransactionVout  `json:"vout"`
	Fee           decimal.Decimal    `json:"fee"`
	Confirmations uint16             `json:"confirmations"`
	BlockHash     entities.BlockHash `json:"blockhash,omitempty"`
}

func (t *Transaction) GetTxID() entities.TxID {
	return t.TxID
}

func (t *Transaction) GetBlockHash() entities.BlockHash {
	return t.BlockHash
}

type TransactionVin struct {
	Txid        entities.TxID           `json:"txid"`
	Vout        uint32                  `json:"vout"`
	ScriptSig   TransactionVinScriptSig `json:"scriptSig"`
	TxInWitness []string                `json:"txinwitness"`
	Sequence    uint32                  `json:"sequence"`
}

type TransactionVinScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type TransactionVout struct {
	Value        decimal.Decimal             `json:"value"`
	N            uint64                      `json:"n"`
	ScriptPubKey TransactionVoutScriptPubKey `json:"scriptPubKey"`
}

type TransactionVoutScriptPubKey struct {
	Address entities.Address `json:"address"`
}
