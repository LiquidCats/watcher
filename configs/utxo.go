package configs

import "time"

type Utxo struct {
	Node UtxoNode `yaml:"node" envconfig:"NODE"`
	Peer UtxoPeer `yaml:"peer" envconfig:"PEER"`
}

type UtxoNode struct {
	URL string `yaml:"url" envconfig:"URL"`
}

type UtxoPeer struct {
	Url               string        `yaml:"url" envconfig:"URL"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout" envconfig:"CONNECTION_TIMEOUT"`
}
