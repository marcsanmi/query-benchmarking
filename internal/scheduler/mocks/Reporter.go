// Code generated by mockery v2.5.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Reporter is an autogenerated mock type for the Reporter type
type Reporter struct {
	mock.Mock
}

// PrintStats provides a mock function with given fields:
func (_m *Reporter) PrintStats() {
	_m.Called()
}

// UpdateStats provides a mock function with given fields: elapsed
func (_m *Reporter) UpdateStats(elapsed time.Duration) {
	_m.Called(elapsed)
}