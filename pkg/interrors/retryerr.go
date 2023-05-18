package interrors

import (
	"time"

	"github.com/av-baran/ymetrics/internal/logger"
)

func RetryOnErr(f func() error) error {
	var resErr error
	backoffLimit := 3
	sleepTime := 1 * time.Second
	timeInc := 2 * time.Second

	for try := 0; try < backoffLimit; try++ {
		if resErr = f(); resErr == nil {
			return resErr
		}
		logger.Errorf("got error: %s, retrying", resErr)
		time.Sleep(sleepTime)
		sleepTime += timeInc
	}
	return resErr
}
