package runner_test

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/LiquidCats/watcher/v2/internal/adapter/runner"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Mock Handler ---.
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

func TestWorker_Run(t *testing.T) {
	type testCase struct {
		name         string
		channelSetup func() chan int
		handlerSetup func(*MockHandler[int])
		workers      uint
		cancelCtx    bool
		wantErr      error
		wantHandled  []int
	}

	tests := []testCase{
		{
			name: "all jobs handled, two workers, channel close",
			channelSetup: func() chan int {
				ch := make(chan int, 3)
				for i := range 3 {
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
		},
		{
			name: "multiple workers handle all jobs",
			channelSetup: func() chan int {
				ch := make(chan int, 10)
				for i := range 10 {
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
			handlerSetup: func(_ *MockHandler[int]) {},
			workers:      2,
			wantErr:      nil,
			wantHandled:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := tt.channelSetup()
			handler := new(MockHandler[int])
			tt.handlerSetup(handler)

			var buf bytes.Buffer
			logger := zerolog.New(&buf).With().Timestamp().Logger()

			ctx := logger.WithContext(t.Context())
			var cancel context.CancelFunc

			if tt.cancelCtx {
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			w := runner.NewWorker[int]("test", tt.workers, ch, handler)
			err := w.Run(ctx)

			// Assert error (if expected)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr, "unexpected error value")
			} else {
				require.NoError(t, err)
			}
			// Assert handled values (regardless of order, since workers may race)
			assert.ElementsMatch(t, tt.wantHandled, handler.Values)
			handler.AssertExpectations(t)
		})
	}
}
