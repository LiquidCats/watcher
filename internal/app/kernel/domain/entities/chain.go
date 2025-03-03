package entities

type Driver string

const (
	DriverRPC Driver = "rpc"
	DriverP2P Driver = "p2p"
)

type Type string

const (
	TypeEvm  Type = "evm"
	TypeUtxo Type = "utxo"
)

type Chain string

const (
	ChainBitcoin  Chain = "bitcoin"
	ChainEthereum Chain = "ethereum"
)
