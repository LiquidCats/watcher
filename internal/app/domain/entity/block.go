package entity

type BlockHash string
type BlockHeight uint32
type BlockStatus string

type Block struct {
	Height   BlockHeight
	Hash     BlockHash
	Previous BlockHash
}

const BlocksTopic = "blocks"

const (
	BlockStatusNew       BlockStatus = "new"
	BlockStatusConfirmed BlockStatus = "confirmed"
	BlockStatusRejected  BlockStatus = "rejected"
)
