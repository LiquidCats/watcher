package null

import (
	"github.com/LiquidCats/watcher/v2/configs"
)

type Publisher struct {
}

func NewPublisher(_ configs.Config) *Publisher {
	return &Publisher{}
}
