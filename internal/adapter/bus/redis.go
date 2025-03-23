package bus

import (
	"context"
	"fmt"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type RedisPublisher struct {
	cfg   configs.Redis
	redis *redis.Client
}

func NewRedisPublisher(cfg configs.Redis) *RedisPublisher {
	r := redis.NewClient(&redis.Options{
		Addr: fmt.Sprint(cfg.Host, ":", cfg.Port),
		DB:   cfg.DB,
	})
	return &RedisPublisher{
		cfg:   cfg,
		redis: r,
	}
}

func (p *RedisPublisher) pub(ctx context.Context, channel string, data any) error {
	dataBytes, err := sonic.ConfigDefault.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "marshal data")
	}

	if err = p.redis.Publish(ctx, channel, dataBytes).Err(); err != nil {
		return errors.Wrap(err, "publish message")
	}

	return nil
}

func (p *RedisPublisher) PublishTransaction(ctx context.Context, transaction entities.Transaction) error {
	if err := p.pub(ctx, p.cfg.TransactionChannel, transaction); err != nil {
		return errors.Wrap(err, "publish transaction")
	}

	return nil
}

func (p *RedisPublisher) PublishBlock(ctx context.Context, block entities.Block) error {
	if err := p.pub(ctx, p.cfg.BlockChannel, block); err != nil {
		return errors.Wrap(err, "publish block")
	}

	return nil
}
