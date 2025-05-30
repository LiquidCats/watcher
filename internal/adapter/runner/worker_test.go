package runner_test

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/LiquidCats/watcher/v2/internal/adapter/runner"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock Handler for testing.
type MockHandler[T any] struct {
	mock.Mock
	mu       sync.Mutex
	Values   []T
	CallTime []time.Time
}

func (m *MockHandler[T]) Handle(ctx context.Context, v T) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Values = append(m.Values, v)
	m.CallTime = append(m.CallTime, time.Now())
	args := m.Called(ctx, v)
	return args.Error(0)
}

func (m *MockHandler[T]) GetValues() []T {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]T{}, m.Values...)
}

func TestWorker_Run_BasicSuccess(t *testing.T) {
	ch := make(chan int, 3)
	handler := new(MockHandler[int])

	// Setup test data
	testValues := []int{1, 2, 3}
	for _, v := range testValues {
		ch <- v
	}

	// Setup mock expectations
	handler.On("Handle", mock.Anything, mock.AnythingOfType("int")).Return(nil)

	// Create worker with logger context
	var buf bytes.Buffer
	logger := zerolog.New(&buf).Level(zerolog.InfoLevel)
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	ctx = logger.WithContext(ctx)

	worker := runner.NewWorker("test-worker", 2, ch, handler)

	go func() {
		time.Sleep(150 * time.Millisecond)
		cancel()
	}()

	_ = worker.Run(ctx)
	receivedValues := handler.GetValues()
	assert.ElementsMatch(t, testValues, receivedValues)
	handler.AssertExpectations(t)
}

func TestWorker_Run_HandlerError(t *testing.T) {
	ch := make(chan string, 2)
	handler := new(MockHandler[string])

	// Setup test data
	ch <- "success"
	ch <- "error"

	// Setup mock to return error for second message
	handler.On("Handle", mock.Anything, "success").Return(nil)
	handler.On("Handle", mock.Anything, "error").Return(errors.New("handler error"))

	var buf bytes.Buffer
	logger := zerolog.New(&buf).Level(zerolog.ErrorLevel)
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	ctx = logger.WithContext(ctx)

	worker := runner.NewWorker("test-worker", 1, ch, handler)

	_ = worker.Run(ctx)

	receivedValues := handler.GetValues()
	assert.ElementsMatch(t, []string{"success", "error"}, receivedValues)
	handler.AssertExpectations(t)

	// Check that error was logged
	logOutput := buf.String()
	assert.Contains(t, logOutput, "handler error")
}

func TestWorker_Run_ContextCancellation(t *testing.T) {
	ch := make(chan int, 10)
	handler := new(MockHandler[int])

	handler.On("Handle", mock.Anything, mock.AnythingOfType("int")).Return(nil)

	// Fill channel with data
	for i := range 5 {
		ch <- i
	}

	// Create cancellable context
	var buf bytes.Buffer
	logger := zerolog.New(&buf).Level(zerolog.DebugLevel)
	ctx, cancel := context.WithCancel(logger.WithContext(t.Context()))

	worker := runner.NewWorker("test-worker", 2, ch, handler)

	// Cancel context after short delay
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	err := worker.Run(ctx)

	// Should return context error
	require.ErrorIs(t, err, context.Canceled)
}

func TestWorker_Run_ContextTimeout(t *testing.T) {
	ch := make(chan int, 1)
	handler := new(MockHandler[int])

	// Create context with very short timeout
	var buf bytes.Buffer
	logger := zerolog.New(&buf).Level(zerolog.DebugLevel)
	ctx, cancel := context.WithTimeout(logger.WithContext(t.Context()), 10*time.Millisecond)
	defer cancel()

	worker := runner.NewWorker("test-worker", 1, ch, handler)

	err := worker.Run(ctx)

	// Should return deadline exceeded error
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestWorker_Run_EmptyChannel(t *testing.T) {
	ch := make(chan int)
	handler := new(MockHandler[int])

	var buf bytes.Buffer
	logger := zerolog.New(&buf).Level(zerolog.ErrorLevel)
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	ctx = logger.WithContext(ctx)

	worker := runner.NewWorker("test-worker", 2, ch, handler)

	// Channel is empty, worker should handle this gracefully
	_ = worker.Run(ctx)
	assert.Empty(t, handler.GetValues())
	handler.AssertNotCalled(t, "Handle")
}

func TestWorker_Run_MultipleWorkers(t *testing.T) {
	ch := make(chan int, 20)
	handler := new(MockHandler[int])

	// Setup test data
	testValues := make([]int, 20)
	for i := range 20 {
		testValues[i] = i
		ch <- i
	}

	// Setup mock expectations
	handler.On("Handle", mock.Anything, mock.AnythingOfType("int")).Return(nil)

	var buf bytes.Buffer
	logger := zerolog.New(&buf).Level(zerolog.ErrorLevel)
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	ctx = logger.WithContext(ctx)

	worker := runner.NewWorker("test-worker", 4, ch, handler)

	_ = worker.Run(ctx)

	receivedValues := handler.GetValues()
	assert.ElementsMatch(t, testValues, receivedValues)
	assert.Len(t, receivedValues, 20)
	handler.AssertExpectations(t)
}

func TestWorker_Run_ZeroWorkers(t *testing.T) {
	ch := make(chan int, 1)
	handler := new(MockHandler[int])

	ch <- 42

	var buf bytes.Buffer
	logger := zerolog.New(&buf).Level(zerolog.InfoLevel)
	ctx := logger.WithContext(t.Context())

	worker := runner.NewWorker("test-worker", 0, ch, handler)

	err := worker.Run(ctx)

	// Should complete successfully but handle no messages
	require.NoError(t, err)
	assert.Empty(t, handler.GetValues())
	handler.AssertNotCalled(t, "Handle")
}

func TestWorker_Run_ConcurrentHandling(t *testing.T) {
	ch := make(chan int, 100)
	handler := new(MockHandler[int])

	// Fill channel with sequential numbers
	for i := range 100 {
		ch <- i
	}

	// Add small delay to simulate work
	handler.On("Handle", mock.Anything, mock.AnythingOfType("int")).
		Return(nil).
		Run(func(_ mock.Arguments) {
			time.Sleep(1 * time.Millisecond)
		})

	var buf bytes.Buffer
	logger := zerolog.New(&buf).Level(zerolog.ErrorLevel)
	ctx, cancel := context.WithTimeout(t.Context(), 150*time.Millisecond)
	defer cancel()
	ctx = logger.WithContext(ctx)

	worker := runner.NewWorker("concurrent-worker", 5, ch, handler)

	start := time.Now()
	_ = worker.Run(ctx)
	duration := time.Since(start)

	receivedValues := handler.GetValues()
	assert.Len(t, receivedValues, 100)

	// With 5 workers, should be faster than sequential processing
	assert.Less(t, duration, 200*time.Millisecond, "Concurrent processing should be reasonably fast")
}

// Benchmark tests.
func BenchmarkWorker_Run_SingleWorker(b *testing.B) {
	ch := make(chan int, b.N)
	handler := new(MockHandler[int])

	for i := range b.N {
		ch <- i
	}

	handler.On("Handle", mock.Anything, mock.AnythingOfType("int")).Return(nil)

	logger := zerolog.Nop()
	ctx := logger.WithContext(b.Context())
	worker := runner.NewWorker("bench-worker", 1, ch, handler)

	b.ResetTimer()
	err := worker.Run(ctx)
	require.NoError(b, err)
}

func BenchmarkWorker_Run_MultipleWorkers(b *testing.B) {
	ch := make(chan int, b.N)
	handler := new(MockHandler[int])

	for i := range b.N {
		ch <- i
	}

	handler.On("Handle", mock.Anything, mock.AnythingOfType("int")).Return(nil)

	logger := zerolog.Nop()
	ctx := logger.WithContext(b.Context())
	worker := runner.NewWorker("bench-worker", 4, ch, handler)

	b.ResetTimer()
	err := worker.Run(ctx)
	require.NoError(b, err)
}
