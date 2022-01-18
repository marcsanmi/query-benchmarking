package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/marcsanmi/query-benchmarking/internal"
	log "github.com/sirupsen/logrus"
)

//go:generate mockery --name=PromscaleClient

// PromscaleClient defines the Promscale client used to query the Promscale server.
type PromscaleClient interface {
	EvaluateExpressionQueryOverRange(ctx context.Context, query internal.Query) (time.Duration, error)
}

//go:generate mockery --name=Reporter

// Reporter defines the reporter in charge of updating the stats.
type Reporter interface {
	PrintStats()
	UpdateStats(elapsed time.Duration)
}

// Scheduler is the orchestrator, the piece that glues up everything. In charge of creating the workers, calling
// the promscale client and printing the final stats.
type Scheduler struct {
	workers         int
	queryChan       chan internal.Query
	doneChan        chan struct{}
	promscaleClient PromscaleClient
	reporter        Reporter
}

func New(
	workers int,
	queryChan chan internal.Query,
	doneChan chan struct{},
	promscaleClient PromscaleClient,
	reporter Reporter,
) Scheduler {
	return Scheduler{
		workers:         workers,
		queryChan:       queryChan,
		doneChan:        doneChan,
		promscaleClient: promscaleClient,
		reporter:        reporter,
	}
}

// Run spawns a goroutine for each worker, reading on the query channel to then make the corresponding http calls and
// reporting. It also listens to done channel to exit from the goroutine.
func (s *Scheduler) Run() {
	defer func() {
		close(s.doneChan)
		close(s.queryChan)
	}()

	wg := sync.WaitGroup{}
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case q := <-s.queryChan:
					elapsed, err := s.promscaleClient.EvaluateExpressionQueryOverRange(context.Background(), q)
					if err != nil {
						log.Errorf("error evaluating range expression query: %v", err.Error())
						continue
					}
					s.reporter.UpdateStats(elapsed)
				case <-s.doneChan:
					return
				}
			}
		}()
	}

	wg.Wait()
	log.Info("All queries have been processed! Printing stats...")
	s.reporter.PrintStats()
}
