package kafka

import (
	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Publisher struct {
	appName  string
	cfg      configs.Config
	producer *kafka.Producer
}

func NewPublisher(appName string, cfg configs.Config) (*Publisher, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.Kafka.Host})
	if err != nil {
		return nil, err
	}

	return &Publisher{
		appName:  appName,
		cfg:      cfg,
		producer: producer,
	}, nil
}
