package configs

import (
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entity"
)

type Config struct {
	Blockchain entity.Blockchain
	NodeUrl    string      `envconfig:"NODE_URL"`
	Kafka      KafkaConfig `envconfig:"KAFKA"`
	Redis      RedisConfig `envconfig:"REDIS"`
}

type KafkaConfig struct {
	Host string
}

type RedisConfig struct {
	Host     string
	Password string
	DB       int
}
