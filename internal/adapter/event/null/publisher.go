package null

import (
	"context"
	"watcher/configs"
	"watcher/internal/app/domain/entity"
)

type Publisher struct {
}

func NewPublisher(_ configs.Config) *Publisher {
	return &Publisher{}
}

func (p *Publisher) NewBlock(_ context.Context, _ *entity.Block) error {
	return nil
}

func (p *Publisher) ConfirmBlock(_ context.Context, _ *entity.Block) error {
	return nil
}

func (p *Publisher) RejectBlock(_ context.Context, _ *entity.Block) error {
	return nil
}
