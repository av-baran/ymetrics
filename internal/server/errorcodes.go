package server

import (
	"net/http"

	"github.com/av-baran/ymetrics/pkg/interrors"
)

func getErrorCode(e error) (statusCode int) {
	switch e.Error() {
	case interrors.ErrInvalidMetricType:
		statusCode = http.StatusNotImplemented
	case interrors.ErrInvalidMetricValue:
		statusCode = http.StatusBadRequest
	case interrors.ErrMetricExistsWithAnotherType:
		statusCode = http.StatusBadRequest
	case interrors.ErrMetricNotFound:
		statusCode = http.StatusNotFound
	default:
		statusCode = http.StatusInternalServerError
	}
	return
}
