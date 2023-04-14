package httphandlers

import (
	"errors"
	"testing"

	"github.com/av-baran/ymetrics/internal/delivery/httphandlers/mocks"
	"github.com/av-baran/ymetrics/internal/entity/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
	"github.com/stretchr/testify/assert"
)

func Test_getMetric(t *testing.T) {
	type Result struct {
		s string
		e error
	}
	type ResultGaugeMock struct {
		v float64
		e error
	}
	type ResultCounterMock struct {
		v int64
		e error
	}
	gaugeTests := []struct {
		name       string
		wantErr    bool
		arg        *metric.Rawdata
		mockFn     string
		mockArg    string
		mockResult ResultGaugeMock
		want       Result
	}{
		{
			name:    "test",
			wantErr: false,
			arg: &metric.Rawdata{
				Name: "name",
				Type: metric.GaugeType,
			},
			mockFn:  "GetGauge",
			mockArg: "name",
			mockResult: ResultGaugeMock{
				v: float64(0.01),
				e: nil,
			},
			want: Result{
				s: "0.01",
				e: nil,
			},
		},
		{
			name:    "test",
			wantErr: true,
			arg: &metric.Rawdata{
				Name: "name",
				Type: metric.GaugeType,
			},
			mockFn:  "GetGauge",
			mockArg: "name",
			mockResult: ResultGaugeMock{
				v: 0.0,
				e: errors.New(interrors.ErrStorageInternalError),
			},
			want: Result{
				s: "",
				e: errors.New(interrors.ErrStorageInternalError),
			},
		},
	}
	counterTests := []struct {
		name       string
		wantErr    bool
		arg        *metric.Rawdata
		mockFn     string
		mockArg    string
		mockResult ResultCounterMock
		want       Result
	}{
		{
			name:    "test",
			wantErr: false,
			arg: &metric.Rawdata{
				Name: "name",
				Type: metric.CounterType,
			},
			mockFn:  "GetCounter",
			mockArg: "name",
			mockResult: ResultCounterMock{
				v: int64(10),
				e: nil,
			},
			want: Result{
				s: "10",
				e: nil,
			},
		},
		{
			name:    "test",
			wantErr: true,
			arg: &metric.Rawdata{
				Name: "name",
				Type: metric.CounterType,
			},
			mockFn:  "GetCounter",
			mockArg: "name",
			mockResult: ResultCounterMock{
				v: 0,
				e: errors.New(interrors.ErrStorageInternalError),
			},
			want: Result{
				s: "",
				e: errors.New(interrors.ErrStorageInternalError),
			},
		},
	}

	for _, tt := range gaugeTests {
		t.Run(tt.name, func(t *testing.T) {
			g := mocks.NewMetricGetter(t)
			g.On(tt.mockFn, tt.mockArg).Once().Return(tt.mockResult.v, tt.mockResult.e)
			gotres, goterr := getMetric(g, tt.arg)
			assert.Equal(t, gotres, tt.want.s)
			if tt.wantErr {
				assert.Error(t, goterr)
			} else {
				assert.NoError(t, goterr)
			}
		})
	}
	for _, tt := range counterTests {
		t.Run(tt.name, func(t *testing.T) {
			g := mocks.NewMetricGetter(t)
			g.On(tt.mockFn, tt.mockArg).Once().Return(tt.mockResult.v, tt.mockResult.e)
			gotres, goterr := getMetric(g, tt.arg)
			assert.Equal(t, gotres, tt.want.s)
			if tt.wantErr {
				assert.Error(t, goterr)
			} else {
				assert.NoError(t, goterr)
			}
		})
	}
}
