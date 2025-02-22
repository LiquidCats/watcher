package entities

type BlockHeight uint64
type BlockHash string

type Block[T any] struct {
	Height       BlockHeight `json:"height"`
	Hash         BlockHash   `json:"hash"`
	PrevHash     BlockHash   `json:"prevHash"`
	Transactions []T         `json:"transactions"`
}
