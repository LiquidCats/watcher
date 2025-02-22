package data

import (
	"github.com/LiquidCats/watcher/v2/internal/app/kernel/entities"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	Txid          entities.TxID       `json:"txid"`
	Vin           []TransactionVin    `json:"vin"`
	Vout          []TransactionVout   `json:"vout"`
	Fee           decimal.Decimal     `json:"fee"`
	Confirmations uint16              `json:"confimations"`
	Blockhash     *entities.BlockHash `json:"blockhash"`
	_             struct{}
}

type TransactionVin struct {
	Txid        entities.TxID           `json:"txid"`
	Vout        uint32                  `json:"vout"`
	ScriptSig   TransactionVinScriptSig `json:"scriptSig"`
	Txinwitness []string                `json:"txinwitness"`
	Sequence    uint32                  `json:"sequence"`
	_           struct{}
}

type TransactionVinScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
	_   struct{}
}

type TransactionVout struct {
	Value        decimal.Decimal             `json:"value"`
	N            int                         `json:"n"`
	ScriptPubKey TransactionVoutScriptPubKey `json:"scriptPubKey"`
	_            struct{}
}

type TransactionVoutScriptPubKey struct {
	Address entities.Address `json:"address"`
	_       struct{}
}
