package configs

type Config struct {
	App App `yaml:"app" envconfig:"APP"`

	Metrics any `envconfig:"METRICS"`

	Chains ChainsConfig `yaml:"chains"`

	DB DB `yaml:"db" envconfig:"DB"`

	Redis Redis `yaml:"redis" envconfig:"REDIS"`
}
