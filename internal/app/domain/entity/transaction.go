package entity

type TransactionId string
type TransactionHash string
type TransactionBlockHash BlockHash

type Transaction struct {
	ID          TransactionId
	Hash        TransactionHash
	BlockHash   BlockHash
	BlockHeight BlockHeight
}
