package data

import "github.com/LiquidCats/watcher/v2/internal/app/domain/entities"

type TransactionReceipts []TransactionReceipt

type TransactionReceipt struct {
	BlockHash         entities.BlockHash   `json:"blockHash"`
	BlockNumber       entities.BlockHeight `json:"blockNumber"`
	ContractAddress   entities.Address     `json:"contractAddress"`
	CumulativeGasUsed HexUint64            `json:"cumulativeGasUsed"`
	EffectiveGasPrice HexUint64            `json:"effectiveGasPrice"`
	From              entities.Address     `json:"from"`
	GasUsed           HexUint64            `json:"gasUsed"`
	Logs              []TransactionLog     `json:"logs"`
	LogsBloom         string               `json:"logsBloom"`
	Status            HexUint64            `json:"status"`
	To                entities.Address     `json:"to"`
	TransactionHash   entities.TxID        `json:"transactionHash"`
	TransactionIndex  HexUint64            `json:"transactionIndex"`
	Type              HexUint64            `json:"type"`
}

type TransactionLog struct {
	Address          entities.Address   `json:"address"`
	Topics           []string           `json:"topics"`
	Data             string             `json:"data"`
	BlockNumber      HexUint64          `json:"blockNumber"`
	TransactionHash  entities.TxID      `json:"transactionHash"`
	TransactionIndex HexUint64          `json:"transactionIndex"`
	BlockHash        entities.BlockHash `json:"blockHash"`
	LogIndex         HexUint64          `json:"logIndex"`
	Removed          bool               `json:"removed"`
}
