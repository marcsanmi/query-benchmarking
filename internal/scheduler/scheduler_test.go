package scheduler_test

import (
	"testing"
	"time"

	"github.com/marcsanmi/query-benchmarking/internal"
	"github.com/marcsanmi/query-benchmarking/internal/scheduler"
	schedulerMocks "github.com/marcsanmi/query-benchmarking/internal/scheduler/mocks"
	"github.com/stretchr/testify/mock"
)

func TestScheduler_Run(t *testing.T) {
	numQueries := 3
	workers := 2
	queryChan := make(chan internal.Query)
	doneChan := make(chan struct{})

	promscaleClientMock := schedulerMocks.PromscaleClient{}
	promscaleClientMock.On("EvaluateExpressionQueryOverRange", mock.Anything, mock.Anything).
		Return(time.Duration(0), nil).Times(numQueries)

	reporterMock := schedulerMocks.Reporter{}
	reporterMock.On("PrintStats").Once()
	reporterMock.On("UpdateStats", mock.AnythingOfType("time.Duration")).Times(numQueries)

	go func() {
		for i := 0; i < numQueries; i++ {
			queryChan <- internal.Query{}
		}

		// Send done signal to the workers.
		for i := 0; i < workers; i++ {
			doneChan <- struct{}{}
		}
	}()

	s := scheduler.New(workers, queryChan, doneChan, &promscaleClientMock, &reporterMock)
	s.Run()

	promscaleClientMock.AssertExpectations(t)
	reporterMock.AssertExpectations(t)
}
