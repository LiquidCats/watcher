package runner_test

import (
	"context"
	"testing"
	"time"

	"github.com/LiquidCats/watcher/v2/internal/adapter/runner"
	"github.com/LiquidCats/watcher/v2/test/mocks"
	"github.com/rotisserie/eris"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProcessor_Run(t *testing.T) {
	var tests = []struct {
		name    string
		perform func(t *testing.T)
	}{
		{
			name: "JobExecutedAndContextCancelled",
			perform: func(t *testing.T) {
				// Arrange
				job := mocks.NewMockJob(t)
				job.On("Handle", mock.Anything).Return(nil)

				ctx := t.Context()

				processor := runner.NewTicker("test", 10*time.Millisecond, job)

				// Cancel context after short time
				ctx, cancel := context.WithTimeout(ctx, 35*time.Millisecond)
				defer cancel()

				// Act
				err := processor.Run(ctx)

				// Assert
				assert.ErrorIs(t, err, context.DeadlineExceeded)
			},
		},
		{
			name: "JobErrorLogged",
			perform: func(t *testing.T) {
				// Arrange
				job := mocks.NewMockJob(t)
				testErr := eris.New("job error")
				job.On("Handle", mock.Anything).Return(testErr)

				ctx := t.Context()

				processor := runner.NewTicker("errjob", 10*time.Millisecond, job)

				// Cancel context after one interval
				ctx, cancel := context.WithTimeout(ctx, 15*time.Millisecond)
				defer cancel()

				// Act
				err := processor.Run(ctx)
				require.Error(t, err)
			},
		},
		{
			name: "StopsOnContextCancelledImmediately",
			perform: func(t *testing.T) {
				// Arrange
				job := mocks.NewMockJob(t)
				job.On("Handle", mock.Anything).Return(nil)

				ctx := t.Context()

				processor := runner.NewTicker("cancelled", 100*time.Millisecond, job)

				// Cancel context immediately
				ctx, cancel := context.WithCancel(ctx)
				cancel()

				// Act
				err := processor.Run(ctx)

				// Assert
				assert.ErrorIs(t, err, context.Canceled)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.perform(t)
		})
	}
}
