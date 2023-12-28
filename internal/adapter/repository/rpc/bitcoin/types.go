package bitcoin

import "watcher/internal/app/domain/entity"

type blockHeader struct {
	Hash              string `json:"hash"`
	Number            int    `json:"number"`
	PreviousBlockHash string `json:"previousblockhash"`
}

func (h *blockHeader) toEntity() *entity.Block {
	return &entity.Block{
		Height:   entity.BlockHeight(h.Number),
		Hash:     entity.BlockHash(h.Hash),
		Previous: entity.BlockHash(h.PreviousBlockHash),
	}
}

type blockStats struct {
	Height int    `json:"height"`
	Hash   string `json:"hash"`
}
