package entities

type (
	BlockHeight uint64
	BlockHash   string
)

type Block interface {
	GetHeight() BlockHeight
	GetHash() BlockHash
	GetPrevHash() BlockHash
	GetTransactions() []Transaction
}
