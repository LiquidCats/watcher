package data

import (
	"strconv"
	"strings"

	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
)

type HexUint64 string

func (x HexUint64) Uint64() uint64 {
	i, _ := strconv.ParseUint(
		strings.TrimPrefix(string(x), "0x"),
		16,
		64,
	)

	return i
}

func (x HexUint64) String() string {
	return string(x)
}

type Block[T any] struct {
	Difficulty      HexUint64          `json:"difficulty"`
	ExtraData       string             `json:"extraData"`
	GasLimit        string             `json:"gasLimit"`
	GasUsed         string             `json:"gasUsed"`
	Hash            entities.BlockHash `json:"hash"`
	LogsBloom       string             `json:"logsBloom"`
	Miner           string             `json:"miner"`
	MixHash         string             `json:"mixHash"`
	Nonce           HexUint64          `json:"nonce"`
	Number          HexUint64          `json:"number"`
	ParentHash      entities.BlockHash `json:"parentHash"`
	ReceiptsRoot    string             `json:"receiptsRoot"`
	Sha3Uncles      string             `json:"sha3Uncles"`
	Size            HexUint64          `json:"size"`
	StateRoot       string             `json:"stateRoot"`
	Timestamp       string             `json:"timestamp"`
	TotalDifficulty string             `json:"totalDifficulty"`
	Transactions    []T                `json:"transactions"`
}

type LatestBlock struct {
	Hash entities.BlockHash `json:"hash"`
}
