package entities

type (
	TxID           string
	RawTransaction string
)

type Transaction interface {
	GetTxID() TxID
	GetBlockHash() BlockHash
}
