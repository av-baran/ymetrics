package interrors

import (
	"time"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/logger"
)

func RetryOnErr(cfg config.RetryConfig, f func() error) error {
	var resErr error

	sleepTime := cfg.SleepTime
	for try := 0; try < cfg.BackoffLimit; try++ {
		if resErr = f(); resErr == nil {
			return resErr
		}
		logger.Errorf("got retryable error: %s, attempt: %s/%s", resErr, try, cfg.BackoffLimit)
		time.Sleep(sleepTime)
		sleepTime += cfg.TimeIncrement
	}
	return resErr
}
