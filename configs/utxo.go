package configs

import "time"

type Evm struct {
}

type Utxo struct {
	NodeUrl string   `yaml:"node_url" envconfig:"NODE_URL"`
	Peer    UtxoPeer `yaml:"peer" envconfig:"PEER"`
}

type UtxoPeer struct {
	Url               string        `yaml:"url" envconfig:"URL"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout" envconfig:"CONNECTION_TIMEOUT"`
}
