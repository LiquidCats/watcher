package entity

type (
	TxID    string
	Address string
	Asset   string
)

const (
	AssetBitcoin  Asset = "BTC"
	AssetEthereum Asset = "ETH"
)

type Transaction struct {
	TxID      TxID                 `json:"tx_id"`
	Contract  *Contract            `json:"contract,omitempty"`
	BlockHash BlockHash            `json:"block_hash"`
	Fee       *Fee                 `json:"fee,omitempty"`
	Outputs   []*TransactionOutput `json:"outputs"`
	Inputs    []*TransactionInput  `json:"inputs"`
}

type Fee struct {
	Asset    Asset  `json:"asset"`
	Value    uint64 `json:"value"`
	Decimals uint64 `json:"decimals"`
}

type Contract struct {
	Asset   Asset   `json:"ticker"`
	Detail  string  `json:"detail"`
	Address Address `json:"address"`
}

type TransactionOutput struct {
	Address  Address `json:"address"`
	Asset    Asset   `json:"asset"`
	Value    uint64  `json:"value"`
	Decimals uint64  `json:"decimals"`
}

type TransactionInput struct {
	Address Address `json:"address"`
}

const TransactionsTopic = "transactions"
