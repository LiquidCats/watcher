package entity

type (
	BlockHash   string
	BlockHeight uint32
	BlockStatus string
)

type Block struct {
	Height   BlockHeight
	Hash     BlockHash
	Previous BlockHash
}

const BlocksTopic = "blocks"
