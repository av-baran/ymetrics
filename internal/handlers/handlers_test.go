package handlers

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/av-baran/ymetrics/internal/entities/metric"
	"github.com/av-baran/ymetrics/internal/httperror"
	"github.com/stretchr/testify/assert"
)

//
// func TestUpdateMetrics(t *testing.T) {
// 	type args struct {
// 		s storage
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want http.HandlerFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := UpdateMetrics(tt.args.s); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("UpdateMetrics() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
func Test_parseURL(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		path    string
		want    *metric.Rawdata
		want1   *httperror.Error
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
			want1: &httperror.Error{
				Msg:  "any",
				Code: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseURL(tt.path)
			if tt.wantErr {
				assert.Nil(t, got)
				assert.NotNil(t, got1)
				assert.Equal(t, tt.want1.Code, got1.Code)
			} else {
				if !reflect.DeepEqual(*got, *tt.want) {
					t.Errorf("parseURL() got = %v, want %v", got, tt.want)
				}
				assert.Nil(t, got1)
			}
		})
	}
}
