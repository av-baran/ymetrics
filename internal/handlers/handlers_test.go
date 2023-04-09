package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/av-baran/ymetrics/internal/entities/metric"
	"github.com/av-baran/ymetrics/internal/interrors"
	memstorage "github.com/av-baran/ymetrics/internal/storage/mem"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{}
	s := memstorage.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(UpdateMetricHandler(s))
			h(w, request)

			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

//
func Test_parseURL(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		path    string
		want    *metric.Rawdata
		want1   error
	}{
		{
			name:    "ok",
			wantErr: false,
			path:    "/update/gauge/name/1.0",
			want: &metric.Rawdata{
				Name:  "name",
				Type:  "gauge",
				Value: "1.0",
			},
			want1: nil,
		},
		{
			name:    "bad url",
			wantErr: true,
			path:    "/update/unknown/type/name/1.0",
			want:    nil,
			want1:   errors.New(interrors.ErrBadUrl),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseURL(tt.path)
			if tt.wantErr {
				assert.Nil(t, got)
				assert.NotNil(t, got1)
				assert.Equal(t, tt.want1.Error(), got1.Error())
			} else {
				if !reflect.DeepEqual(*got, *tt.want) {
					t.Errorf("parseURL() got = %v, want %v", got, tt.want)
				}
				assert.Nil(t, got1)
			}
		})
	}
}
