package configs

type Config struct {
	App  App  `yaml:"app" envconfig:"APP"`
	Utxo Utxo `yaml:"utxo" envconfig:"UTXO"`

	DB DB `yaml:"db" envconfig:"DB"`
}
