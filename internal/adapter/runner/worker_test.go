package runner_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/LiquidCats/watcher/v2/internal/adapter/runner"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock Handler ---
type MockHandler[T any] struct {
	mock.Mock
	mu     sync.Mutex
	Values []T
}

func (m *MockHandler[T]) Handle(ctx context.Context, v T) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Values = append(m.Values, v)
	args := m.Called(ctx, v)
	return args.Error(0)
}

// --- Helper to create logger in context and capture output ---
func loggerToBuf() (*zerolog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf).With().Timestamp().Logger()
	return &logger, &buf
}

func TestWorker_Run(t *testing.T) {
	type testCase struct {
		name         string
		channelSetup func() chan int
		handlerSetup func(*MockHandler[int])
		workers      uint
		cancelCtx    bool
		wantErr      error
		wantHandled  []int
		wantInLog    []string
	}

	tests := []testCase{
		{
			name: "all jobs handled, two workers, channel close",
			channelSetup: func() chan int {
				ch := make(chan int, 3)
				for i := 0; i < 3; i++ {
					ch <- i
				}
				close(ch)
				return ch
			},
			handlerSetup: func(m *MockHandler[int]) {
				m.On("Handle", mock.Anything, mock.AnythingOfType("int")).Return(nil)
			},
			workers:     2,
			wantErr:     nil,
			wantHandled: []int{0, 1, 2},
			wantInLog:   []string{"runner channel closed"},
		},
		{
			name: "handler returns error, error logged",
			channelSetup: func() chan int {
				ch := make(chan int, 1)
				ch <- 42
				close(ch)
				return ch
			},
			handlerSetup: func(m *MockHandler[int]) {
				m.On("Handle", mock.Anything, 42).Return(errors.New("handle failed"))
			},
			workers:     1,
			wantErr:     nil,
			wantHandled: []int{42},
			wantInLog:   []string{"runner handler error", "handle failed"},
		},
		{
			name: "multiple workers handle all jobs",
			channelSetup: func() chan int {
				ch := make(chan int, 10)
				for i := 0; i < 10; i++ {
					ch <- i
				}
				close(ch)
				return ch
			},
			handlerSetup: func(m *MockHandler[int]) {
				m.On("Handle", mock.Anything, mock.AnythingOfType("int")).Return(nil)
			},
			workers:     4,
			wantErr:     nil,
			wantHandled: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "closed channel with no jobs",
			channelSetup: func() chan int {
				ch := make(chan int)
				close(ch)
				return ch
			},
			handlerSetup: func(m *MockHandler[int]) {},
			workers:      2,
			wantErr:      nil,
			wantHandled:  []int{},
			wantInLog:    []string{"runner channel closed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := tt.channelSetup()
			handler := new(MockHandler[int])
			tt.handlerSetup(handler)

			logger, buf := loggerToBuf()
			ctx := logger.WithContext(context.Background())
			var cancel context.CancelFunc

			if tt.cancelCtx {
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			w := runner.NewWorker[int](tt.workers, ch, handler)
			err := w.Run(ctx)

			// Assert error (if expected)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr, "unexpected error value")
			} else {
				assert.NoError(t, err)
			}
			// Assert handled values (regardless of order, since workers may race)
			assert.ElementsMatch(t, tt.wantHandled, handler.Values)
			handler.AssertExpectations(t)
			// Assert log messages if needed
			if len(tt.wantInLog) > 0 {
				logStr := buf.String()
				for _, substr := range tt.wantInLog {
					assert.Contains(t, logStr, substr, fmt.Sprintf("expected log to contain: %q", substr))
				}
			}
		})
	}
}
