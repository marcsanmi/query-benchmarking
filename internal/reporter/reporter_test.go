package reporter_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/marcsanmi/query-benchmarking/internal/reporter"
	"github.com/stretchr/testify/assert"
)

func TestReporter_UpdateStats(t *testing.T) {
	report := reporter.New()

	tests := []struct {
		name    string
		elapsed time.Duration
	}{
		{
			name:    "Add 1st elapsed",
			elapsed: 10 * time.Millisecond,
		},
		{
			name:    "Insert 2nd elapsed",
			elapsed: 5 * time.Millisecond,
		},
		{
			name:    "Add 3rd elapsed",
			elapsed: 3 * time.Millisecond,
		},
	}

	t.Run("group", func(t *testing.T) {
		for _, tt := range tests {
			tt := tt // // capture range variable
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				report.UpdateStats(tt.elapsed)
			})
		}
	})

	// Kind of cheating to assert the report values by accessing unexported fields.
	rs := reflect.ValueOf(report).Elem()
	assert.Equal(t, int64(3), rs.Field(0).Int())        // totalQueriesRun
	assert.Equal(t, int64(10000000), rs.Field(2).Int()) // maxQueryTime
	assert.Equal(t, int64(3000000), rs.Field(3).Int())  // minQueryTime
	assert.Equal(t, int64(18000000), rs.Field(4).Int()) // totalProcessingTime
}

func TestReporter_CalcMedian(t *testing.T) {
	tests := []struct {
		name  string
		input []time.Duration
		want  time.Duration
	}{
		{
			name:  "Empty input",
			input: []time.Duration{},
			want:  0,
		},
		{
			name: "Odd input length",
			input: []time.Duration{
				3 * time.Millisecond,
				10 * time.Millisecond,
				7 * time.Millisecond,
			},
			want: 7 * time.Millisecond,
		},
		{
			name: "Even input length",
			input: []time.Duration{
				3 * time.Millisecond,
				10 * time.Millisecond,
				7 * time.Millisecond,
				5 * time.Millisecond,
			},
			want: 6 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reporter.CalcMedian(tt.input)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestReporter_CalcAverage(t *testing.T) {
	tests := []struct {
		name           string
		processingTime time.Duration
		queriesRun     int
		want           time.Duration
	}{
		{
			name:       "No queries run",
			queriesRun: 0,
			want:       0,
		},
		{
			name:           "Average simple calculation",
			processingTime: 10 * time.Millisecond,
			queriesRun:     5,
			want:           2 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reporter.CalcAverage(tt.processingTime, tt.queriesRun)
			assert.Equal(t, tt.want, result)
		})
	}
}
