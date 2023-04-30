package server

import (
	"errors"
	"net/http"

	"github.com/av-baran/ymetrics/pkg/interrors"
)

func getErrorCode(e error) (statusCode int) {
	switch {
	case errors.Is(e, interrors.ErrInvalidMetricType):
		statusCode = http.StatusNotImplemented
	case errors.Is(e, interrors.ErrInvalidMetricValue):
		statusCode = http.StatusBadRequest
	case errors.Is(e, interrors.ErrMetricExistsWithAnotherType):
		statusCode = http.StatusBadRequest
	case errors.Is(e, interrors.ErrMetricNotFound):
		statusCode = http.StatusNotFound
	default:
		statusCode = http.StatusInternalServerError
	}
	return
}
