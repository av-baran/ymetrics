package server

import (
	"net/http"
	"testing"

	"github.com/av-baran/ymetrics/pkg/interrors"
)

func Test_getErrorCode(t *testing.T) {
	tests := []struct {
		name           string
		e              error
		wantStatusCode int
	}{
		{
			name:           "invalid metric type",
			e:              interrors.ErrInvalidMetricType,
			wantStatusCode: http.StatusNotImplemented,
		},
		{
			name:           "invalid metric value",
			e:              interrors.ErrInvalidMetricValue,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "metric exists with another type",
			e:              interrors.ErrMetricExistsWithAnotherType,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "metric not found",
			e:              interrors.ErrMetricNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "unknow error",
			e:              interrors.ErrStorageInternalError,
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotStatusCode := getErrorCode(tt.e); gotStatusCode != tt.wantStatusCode {
				t.Errorf("getErrorCode() = %v, want %v", gotStatusCode, tt.wantStatusCode)
			}
		})
	}
}
