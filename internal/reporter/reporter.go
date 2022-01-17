package reporter

import (
	"math"
	"sort"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Report defines the fields that I need to keep track in order to print the report and
// holds a mutex as well in order to update the stats in a thread safe way.
type Report struct {
	totalQueriesRun     int
	elapsedTimes        []time.Duration
	maxQueryTime        time.Duration
	minQueryTime        time.Duration
	totalProcessingTime time.Duration
	mu                  sync.Mutex
}

func New() *Report {
	return &Report{
		mu: sync.Mutex{},
	}
}

// UpdateStats updates the shared state in a safe way given an elapsed query time.
func (r *Report) UpdateStats(elapsed time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.totalQueriesRun++
	r.totalProcessingTime += elapsed

	if r.maxQueryTime < elapsed {
		r.maxQueryTime = elapsed
	}

	if r.minQueryTime == 0 || r.minQueryTime > elapsed {
		r.minQueryTime = elapsed
	}

	r.elapsedTimes = append(r.elapsedTimes, elapsed)
}

// PrintStats prints all the stats in the form of the specified requirements.
func (r *Report) PrintStats() {
	log.Info("> Stats summary:")

	r.mu.Lock()
	defer r.mu.Unlock()

	log.Infof("\t # of queries processed: %d", r.totalQueriesRun)
	log.Infof("\t total processing time: %s", r.totalProcessingTime)
	log.Infof("\t maximum query time: %s", r.maxQueryTime)
	log.Infof("\t minimum query time: %s", r.minQueryTime)
	log.Infof("\t average query time: %s", CalcAverage(r.totalProcessingTime, r.totalQueriesRun))
	log.Infof("\t median  query time: %s", CalcMedian(r.elapsedTimes))
}

// CalcAverage computes the average query time given the total queries processing time and the total queries run.
func CalcAverage(totalProcessingTime time.Duration, totalQueriesRun int) time.Duration {
	if totalQueriesRun == 0 {
		return 0
	}

	return time.Duration(math.Round(float64(totalProcessingTime) / float64(totalQueriesRun)))
}

// CalcMedian computes the median query time given the samples of all query elapsed times.
func CalcMedian(elapsedTimes []time.Duration) time.Duration {
	if len(elapsedTimes) == 0 {
		return 0
	}

	sort.Slice(elapsedTimes, func(i, j int) bool {
		return int64(elapsedTimes[i]) < int64(elapsedTimes[j])
	})

	midTime := len(elapsedTimes) / 2

	// Odd length.
	if len(elapsedTimes)%2 != 0 {
		return elapsedTimes[midTime]
	}

	// Even length.
	return (elapsedTimes[midTime-1] + elapsedTimes[midTime]) / 2
}
