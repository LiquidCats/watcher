package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"watcher/configs"
	"watcher/internal/app/domain/entity"
	"watcher/internal/app/domain/utils"
)

type Publisher struct {
	appName  string
	cfg      configs.Config
	producer *redis.Client
}

func NewPublisher(appName string, cfg configs.Config) *Publisher {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return &Publisher{
		appName:  appName,
		cfg:      cfg,
		producer: rdb,
	}
}

func (p *Publisher) NewBlock(ctx context.Context, blockchain entity.Blockchain, block *entity.Block) error {
	topic := utils.MakeBlocksTopic(p.appName, blockchain, entity.BlockStatusNew)
	b, err := json.Marshal(block)
	if err != nil {
		return err
	}

	cmd := p.producer.Publish(ctx, topic, b)

	return cmd.Err()
}

func (p *Publisher) ConfirmBlock(ctx context.Context, blockchain entity.Blockchain, block *entity.Block) error {
	topic := utils.MakeBlocksTopic(p.appName, blockchain, entity.BlockStatusConfirmed)
	b, err := json.Marshal(block)
	if err != nil {
		return err
	}

	cmd := p.producer.Publish(ctx, topic, b)

	return cmd.Err()
}

func (p *Publisher) RejectBlock(ctx context.Context, blockchain entity.Blockchain, block *entity.Block) error {
	topic := utils.MakeBlocksTopic(p.appName, blockchain, entity.BlockStatusRejected)
	b, err := json.Marshal(block)
	if err != nil {
		return err
	}

	cmd := p.producer.Publish(ctx, topic, b)

	return cmd.Err()
}
