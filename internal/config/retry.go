package config

import "time"

const (
	retryDefaultBackoffLimit  = 3
	retryDefaultSleepTime     = 1 * time.Second
	retryDefaultTimeIncrement = 2 * time.Second
)

type RetryConfig struct {
	BackoffLimit  int
	SleepTime     time.Duration
	TimeIncrement time.Duration
}
