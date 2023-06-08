package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/pkg/interrors"
)

func sendError(w http.ResponseWriter, msg string, e error) {
	errorMsg := fmt.Sprintf("%s: %s", msg, e)
	logger.Error(errorMsg)
	http.Error(w, errorMsg, getErrorCode(e))
}

func getErrorCode(e error) (statusCode int) {
	switch {
	case errors.Is(e, interrors.ErrInvalidMetricType):
		statusCode = http.StatusNotImplemented

	case
		errors.Is(e, interrors.ErrInvalidMetricValue),
		errors.Is(e, interrors.ErrMetricExistsWithAnotherType),
		errors.Is(e, strconv.ErrSyntax),
		errors.Is(e, strconv.ErrRange),
		errors.Is(e, interrors.ErrInvalidSign):

		statusCode = http.StatusBadRequest

	case
		errors.Is(e, interrors.ErrMetricNotFound):

		statusCode = http.StatusNotFound

	case
		errors.Is(e, interrors.ErrPingDB):

		statusCode = http.StatusInternalServerError

	default:
		logger.Error("unknown internal error happens: %s", e)
		statusCode = http.StatusInternalServerError
	}
	return
}
