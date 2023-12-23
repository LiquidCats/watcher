package configs

import "github.com/kelseyhightower/envconfig"

func Load(prefix string) (Config, error) {
	var cfg Config

	err := envconfig.Process(prefix, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
