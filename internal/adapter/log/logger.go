package log

import (
	"github.com/LiquidCats/watcher/v2/internal/port/adapter/log"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}
type SyncFunc func() error

func NewNopLogger() log.Logger[zap.Field] {
	return &Logger{zap.NewNop()}
}

func NewLogger(appName string) (log.Logger[zap.Field], SyncFunc, error) {
	l, err := zap.NewProduction()
	if err != nil {
		return nil, nil, err
	}

	return &Logger{l.Named(appName)}, l.Sync, nil
}

func (l *Logger) Named(name string) *Logger {
	named := l.Logger.Named(name)
	return &Logger{named}
}
