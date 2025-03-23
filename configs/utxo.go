package configs

import "time"

type Utxo struct {
	RPC  UtxoRPC  `yaml:"rpc" envconfig:"RPC"`
	Peer UtxoPeer `yaml:"peer" envconfig:"PEER"`
}

type UtxoRPC struct {
	URL string `yaml:"url" envconfig:"URL"`
}

type UtxoPeer struct {
	URL               string        `yaml:"url" envconfig:"URL"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout" envconfig:"CONNECTION_TIMEOUT"`
}
