package entities

type BlockHeight uint64
type BlockHash string

type UtxoBlock struct {
	Height       BlockHeight        `json:"height"`
	Hash         BlockHash          `json:"hash"`
	PrevHash     BlockHash          `json:"prevHash"`
	Transactions []*UtxoTransaction `json:"transactions"`
}
