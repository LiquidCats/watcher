package entity

type TxId string
type Address string

type Transaction struct {
	TxId      TxId
	BlockHash BlockHash
	Outputs   []*TransactionOutput
	Inputs    []*TransactionInput
}

type TransactionOutput struct {
	Address  Address
	Value    uint64
	Decimals uint64
}

type TransactionInput struct {
	Address Address
}
