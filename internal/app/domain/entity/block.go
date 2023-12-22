package entity

type BlockHash string
type BlockHeight uint32

type Block struct {
	Height   BlockHeight
	Hash     BlockHash
	Previous BlockHash
}
