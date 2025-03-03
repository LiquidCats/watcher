package configs

import "time"

type Utxo struct {
	Rpc  UtxoRpc  `yaml:"rpc" envconfig:"RPC"`
	Peer UtxoPeer `yaml:"peer" envconfig:"PEER"`
}

type UtxoRpc struct {
	URL string `yaml:"url" envconfig:"URL"`
}

type UtxoPeer struct {
	Url               string        `yaml:"url" envconfig:"URL"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout" envconfig:"CONNECTION_TIMEOUT"`
}
