package entities

type TxID string

type Transaction interface {
	GetTxID() TxID
	GetBlockHash() BlockHash
}
