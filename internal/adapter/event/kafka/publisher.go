package kafka

import (
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"watcher/configs"
	"watcher/internal/app/domain/entity"
	"watcher/internal/app/domain/utils"
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

func (p *Publisher) sendMessage(topic string, data []byte) error {
	tp := kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}
	message := &kafka.Message{
		TopicPartition: tp,
		Value:          data,
	}
	if err := p.producer.Produce(message, nil); err != nil {
		return err
	}

	p.producer.Flush(300)

	return nil
}

func (p *Publisher) NewBlock(_ context.Context, blockchain entity.Blockchain, block *entity.Block) error {
	topicName := utils.MakeBlocksTopic(p.appName, blockchain, entity.BlockStatusNew)

	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	return p.sendMessage(topicName, data)
}

func (p *Publisher) ConfirmBlock(_ context.Context, blockchain entity.Blockchain, block *entity.Block) error {
	topicName := utils.MakeBlocksTopic(p.appName, blockchain, entity.BlockStatusConfirmed)

	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	return p.sendMessage(topicName, data)
}

func (p *Publisher) RejectBlock(_ context.Context, blockchain entity.Blockchain, block *entity.Block) error {
	topicName := utils.MakeBlocksTopic(p.appName, blockchain, entity.BlockStatusRejected)

	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	return p.sendMessage(topicName, data)
}
