package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/pkg/interrors"
)

func getErrorCode(e error) (statusCode int) {
	switch {
	case errors.Is(e, interrors.ErrInvalidMetricType):
		statusCode = http.StatusNotImplemented

	case
		errors.Is(e, interrors.ErrInvalidMetricValue),
		errors.Is(e, interrors.ErrMetricExistsWithAnotherType),
		errors.Is(e, strconv.ErrSyntax),
		errors.Is(e, strconv.ErrRange):

		statusCode = http.StatusBadRequest

	case errors.Is(e, interrors.ErrMetricNotFound):
		statusCode = http.StatusNotFound

	default:
		logger.Log.Sugar().Debugln("unknown internal error happens: %s", e)
		statusCode = http.StatusInternalServerError
	}
	return
}
