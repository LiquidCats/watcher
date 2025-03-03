package data

import (
	kernel "github.com/LiquidCats/watcher/v2/internal/app/kernel/domain/entities"
	"github.com/LiquidCats/watcher/v2/internal/app/utxo/domain/entities"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	Txid          kernel.TxID       `json:"txid"`
	Vin           []TransactionVin  `json:"vin"`
	Vout          []TransactionVout `json:"vout"`
	Fee           decimal.Decimal   `json:"fee"`
	Confirmations uint16            `json:"confimations"`
	Blockhash     kernel.BlockHash  `json:"blockhash,omitempty"`
}

func (t *Transaction) ToEntity() *entities.Transaction {
	ent := &entities.Transaction{
		TxID:      t.Txid,
		BlockHash: t.Blockhash,
		Inputs:    make([]*entities.Input, len(t.Vin)),
		Outputs:   make([]*entities.Output, len(t.Vout)),
	}

	for _, in := range t.Vin {
		ent.Inputs = append(ent.Inputs, in.ToEntity())
	}

	for _, out := range t.Vout {
		ent.Outputs = append(ent.Outputs, out.ToEntity())
	}

	return ent
}

type TransactionVin struct {
	Txid        kernel.TxID             `json:"txid"`
	Vout        uint32                  `json:"vout"`
	ScriptSig   TransactionVinScriptSig `json:"scriptSig"`
	TxInWitness []string                `json:"txinwitness"`
	Sequence    uint32                  `json:"sequence"`
}

func (ti *TransactionVin) ToEntity() *entities.Input {
	return &entities.Input{
		TxID: ti.Txid,
		Vout: ti.Vout,
	}
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

func (to *TransactionVout) ToEntity() *entities.Output {
	return &entities.Output{
		Address: to.ScriptPubKey.Address,
		Value:   to.Value.String(),
		N:       to.N,
	}
}

type TransactionVoutScriptPubKey struct {
	Address kernel.Address `json:"address"`
}
