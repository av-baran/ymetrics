package httphandlers

// FIXME
// import (
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"reflect"
// 	"testing"
//
// 	"github.com/av-baran/ymetrics/internal/entities/metric"
// 	"github.com/av-baran/ymetrics/internal/interrors"
// 	"github.com/av-baran/ymetrics/internal/storage/memstor"
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestUpdateMetricHandler(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		request      string
// 		method       string
// 		expectedCode int
// 		expectedBody string
// 	}{
// 		// {
// 		// 	name:         "GET",
// 		// 	request:      "/update/gauge/name/1.0",
// 		// 	method:       http.MethodGet,
// 		// 	expectedCode: http.StatusMethodNotAllowed,
// 		// 	expectedBody: "",
// 		// },
// 		// {
// 		// 	name:         "PUT",
// 		// 	request:      "/update/gauge/name/1.0",
// 		// 	method:       http.MethodPut,
// 		// 	expectedCode: http.StatusMethodNotAllowed,
// 		// 	expectedBody: "",
// 		// },
// 		// {
// 		// 	name:         "DELETE",
// 		// 	request:      "/update/gauge/name/1.0",
// 		// 	method:       http.MethodDelete,
// 		// 	expectedCode: http.StatusMethodNotAllowed,
// 		// 	expectedBody: "",
// 		// },
// 		{
// 			name:         "POST - Bad Type",
// 			request:      "/update/unknowntype/name/1.0",
// 			method:       http.MethodPost,
// 			expectedCode: http.StatusNotImplemented,
// 			expectedBody: "",
// 		},
// 		{
// 			name:         "POST - Post Bad value",
// 			request:      "/update/gauge/name/fuuuu",
// 			method:       http.MethodPost,
// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: "",
// 		},
// 		{
// 			name:         "POST - OK",
// 			request:      "/update/gauge/name/1.0",
// 			method:       http.MethodPost,
// 			expectedCode: http.StatusOK,
// 			expectedBody: "",
// 		},
// 		{
// 			//FIXME depends on previous test value. Should be same metric name
// 			name:         "POST - already exists",
// 			request:      "/update/counter/name/1.0",
// 			method:       http.MethodPost,
// 			expectedCode: http.StatusBadRequest,
// 			expectedBody: "",
// 		},
// 	}
//
// 	s := memstor.New()
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			request := httptest.NewRequest(tt.method, tt.request, nil)
// 			w := httptest.NewRecorder()
// 			h := http.HandlerFunc(UpdateMetricHandler(s))
// 			h(w, request)
//
// 			result := w.Result()
// 			defer result.Body.Close()
//
// 			assert.Equal(t, tt.expectedCode, result.StatusCode)
// 		})
// 	}
//
// }
//
// func Test_parseURL(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		wantErr bool
// 		path    string
// 		want    *metric.Rawdata
// 		want1   error
// 	}{
// 		{
// 			name:    "ok",
// 			wantErr: false,
// 			path:    "/update/gauge/name/1.0",
// 			want: &metric.Rawdata{
// 				Name:  "name",
// 				Type:  "gauge",
// 				Value: "1.0",
// 			},
// 			want1: nil,
// 		},
// 		{
// 			name:    "bad url",
// 			wantErr: true,
// 			path:    "/update/unknown/type/name/1.0",
// 			want:    nil,
// 			want1:   errors.New(interrors.ErrBadURL),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, got1 := parseURL(tt.path)
// 			if tt.wantErr {
// 				assert.Nil(t, got)
// 				assert.NotNil(t, got1)
// 				assert.Equal(t, tt.want1.Error(), got1.Error())
// 			} else {
// 				if !reflect.DeepEqual(*got, *tt.want) {
// 					t.Errorf("parseURL() got = %v, want %v", got, tt.want)
// 				}
// 				assert.Nil(t, got1)
// 			}
// 		})
// 	}
// }
