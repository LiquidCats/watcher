package ethereum

import (
	"strconv"
	"strings"
	"watcher/internal/app/domain/entity"
)

type block struct {
	Hash       string `json:"hash"`
	Height     string `json:"height"`
	ParentHash string `json:"parent_hash"`
}

func (r *block) toEntity() *entity.Block {
	height := strings.ToLower(r.Height)
	height = strings.Replace(height, "0x", "", -1)
	output, _ := strconv.ParseUint(height, 16, 64)

	return &entity.Block{
		Height:   entity.BlockHeight(output),
		Hash:     entity.BlockHash(r.Hash),
		Previous: entity.BlockHash(r.ParentHash),
	}
}
