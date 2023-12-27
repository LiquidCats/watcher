package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"strings"
	"watcher/configs"
	"watcher/internal/app/domain/entity"
)

type Publisher struct {
	appName  string
	cfg      configs.Config
	producer *kafka.Producer
}

func NewPublisher(appName string, cfg configs.Config) (*Publisher, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.Broadcaster.Host})
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

func (p *Publisher) makeBlocksTopic(s ...string) string {
	str := strings.Join(s, "-")
	return fmt.Sprint(p.appName, ".blocks-", str)
}

func (p *Publisher) NewBlock(_ context.Context, blockchain entity.Blockchain, block *entity.Block) error {
	topicName := p.makeBlocksTopic(string(blockchain), "new")

	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	return p.sendMessage(topicName, data)
}

func (p *Publisher) ConfirmBlock(_ context.Context, blockchain entity.Blockchain, block *entity.Block) error {
	topicName := p.makeBlocksTopic(string(blockchain), "confirm")

	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	return p.sendMessage(topicName, data)
}

func (p *Publisher) RejectBlock(_ context.Context, blockchain entity.Blockchain, block *entity.Block) error {
	topicName := p.makeBlocksTopic(string(blockchain), "reject")

	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	return p.sendMessage(topicName, data)
}
