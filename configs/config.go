package configs

type Config struct {
	App  App  `yaml:"app" envconfig:"APP"`
	Utxo Utxo `yaml:"utxo" envconfig:"UTXO"`
	Evm  Evm  `yaml:"evm" envconfig:"EVM"`

	DB DB `yaml:"db" envconfig:"DB"`

	Redis Redis `yaml:"redis" envconfig:"REDIS"`
}
