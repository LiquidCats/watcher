package entities

type TxID string

type UtxoTransaction struct {
	TxID      TxID      `json:"tx_id"`
	BlockHash BlockHash `json:"block_hash"`
	Inputs    []*Input  `json:"inputs,omitempty"`
	Outputs   []*Output `json:"outputs,omitempty"`
}

type Input struct {
	TxID TxID   `json:"tx_id"`
	Vout uint32 `json:"vout"`
}

type Output struct {
	Address Address
	Value   string
	N       uint64
}
