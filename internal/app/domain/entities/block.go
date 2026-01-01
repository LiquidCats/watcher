package entities

type (
	BlockHeight uint64
	BlockHash   string
	RawBlock    string
)

type Block struct {
	Height       BlockHeight
	Hash         BlockHash
	Transactions []TxID
}

type BlockWithTransactions[TxIn any] struct {
	Height       BlockHeight
	Hash         BlockHash
	Transactions []Transaction[TxIn]
}
