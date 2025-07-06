package configs

import "github.com/LiquidCats/graceful"

type Config struct {
	App App `yaml:"app" envconfig:"APP"`

	Metrics graceful.HttpConfig `envconfig:"METRICS"`

	Chains ChainsConfig `yaml:"chains"`

	DB DB `yaml:"db" envconfig:"DB"`

	Redis Redis `yaml:"redis" envconfig:"REDIS"`
}
