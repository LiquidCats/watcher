package configs

type Redis struct {
	Host string `yaml:"host" envconfig:"HOST"`
	Port int    `yaml:"port" envconfig:"PORT"`
	DB   int    `yaml:"db" envconfig:"DB"`

	BlockChannel       string `yaml:"block_channel" envconfig:"BLOCK_CHANNEL"`
	TransactionChannel string `yaml:"transaction_channel" envconfig:"TX_CHANNEL"`
}
