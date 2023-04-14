package httphandlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/av-baran/ymetrics/pkg/interrors"
)

func Test_getErrorCode(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		wantStatusCode int
	}{
		{
			name:           "Invalid metric type",
			err:            errors.New(interrors.ErrInvalidMetricType),
			wantStatusCode: http.StatusNotImplemented,
		},
		{
			name:           "Invalid metric value",
			err:            errors.New(interrors.ErrInvalidMetricValue),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Metric with another type",
			err:            errors.New(interrors.ErrMetricExistsWithAnotherType),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Metric not found",
			err:            errors.New(interrors.ErrMetricNotFound),
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "Unknown error",
			err:            errors.New("unknown error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotStatusCode := getErrorCode(tt.err); gotStatusCode != tt.wantStatusCode {
				t.Errorf("getErrorCode() = %v, want %v", gotStatusCode, tt.wantStatusCode)
			}
		})
	}
}
