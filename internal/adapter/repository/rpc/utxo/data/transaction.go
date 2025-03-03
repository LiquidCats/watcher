package data

import (
	entities2 "github.com/LiquidCats/watcher/v2/internal/app/kernel/domain/entities"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	Txid          entities2.TxID       `json:"txid"`
	Vin           []TransactionVin     `json:"vin"`
	Vout          []TransactionVout    `json:"vout"`
	Fee           decimal.Decimal      `json:"fee"`
	Confirmations uint16               `json:"confimations"`
	Blockhash     *entities2.BlockHash `json:"blockhash,omitempty"`
}

type TransactionVin struct {
	Txid        entities2.TxID          `json:"txid"`
	Vout        uint32                  `json:"vout"`
	ScriptSig   TransactionVinScriptSig `json:"scriptSig"`
	Txinwitness []string                `json:"txinwitness"`
	Sequence    uint32                  `json:"sequence"`
}

type TransactionVinScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type TransactionVout struct {
	Value        decimal.Decimal             `json:"value"`
	N            int                         `json:"n"`
	ScriptPubKey TransactionVoutScriptPubKey `json:"scriptPubKey"`
}

type TransactionVoutScriptPubKey struct {
	Address entities2.Address `json:"address"`
}
