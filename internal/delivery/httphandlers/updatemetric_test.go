package httphandlers

import (
	"errors"
	"testing"

	"github.com/av-baran/ymetrics/internal/delivery/httphandlers/mocks"
	"github.com/av-baran/ymetrics/internal/entity/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
	"github.com/stretchr/testify/assert"
)

func Test_storeMetric(t *testing.T) {
	counterTests := []struct {
		name       string
		wantErr    bool
		Arg        *metric.Rawdata
		mockFn     string
		mockArg    *metric.Counter
		mockResult error
		want       error
	}{
		{
			name:    "test",
			wantErr: false,
			Arg: &metric.Rawdata{
				Name:  "name",
				Type:  metric.CounterType,
				Value: "100",
			},
			mockFn: "UpdateCounter",
			mockArg: &metric.Counter{
				Name:  "name",
				Value: int64(100),
			},
			mockResult: nil,
			want:       nil,
		},
		{
			name:    "test",
			wantErr: false,
			Arg: &metric.Rawdata{
				Name:  "name",
				Type:  metric.CounterType,
				Value: "0",
			},
			mockFn: "UpdateCounter",
			mockArg: &metric.Counter{
				Name:  "name",
				Value: int64(0),
			},
			mockResult: nil,
			want:       nil,
		},
		{
			name:    "test",
			wantErr: true,
			Arg: &metric.Rawdata{
				Name:  "name",
				Type:  metric.CounterType,
				Value: "nope",
			},
			mockFn: "UpdateCounter",
			mockArg: &metric.Counter{
				Name:  "name",
				Value: int64(1000),
			},
			mockResult: errors.New(interrors.ErrInvalidMetricValue),
			want:       errors.New(interrors.ErrInvalidMetricValue),
		},
	}
	gaugeTests := []struct {
		name       string
		wantErr    bool
		Arg        *metric.Rawdata
		mockFn     string
		mockArg    *metric.Gauge
		mockResult error
		want       error
	}{
		{
			name:    "test",
			wantErr: false,
			Arg: &metric.Rawdata{
				Name:  "name",
				Type:  metric.GaugeType,
				Value: "0.0001",
			},
			mockFn: "UpdateGauge",
			mockArg: &metric.Gauge{
				Name:  "name",
				Value: float64(0.0001),
			},
			mockResult: nil,
			want:       nil,
		},
		{
			name:    "test",
			wantErr: false,
			Arg: &metric.Rawdata{
				Name:  "name",
				Type:  metric.GaugeType,
				Value: "0.0",
			},
			mockFn: "UpdateGauge",
			mockArg: &metric.Gauge{
				Name:  "name",
				Value: float64(0.0),
			},
			mockResult: nil,
			want:       nil,
		},
		{
			name:    "test",
			wantErr: true,
			Arg: &metric.Rawdata{
				Name:  "name",
				Type:  metric.GaugeType,
				Value: "nope",
			},
			mockFn: "UpdateGauge",
			mockArg: &metric.Gauge{
				Name:  "name",
				Value: float64(0.0),
			},
			mockResult: errors.New(interrors.ErrInvalidMetricValue),
			want:       errors.New(interrors.ErrInvalidMetricValue),
		},
	}
	for _, tt := range gaugeTests {
		t.Run(tt.name, func(t *testing.T) {
			u := mocks.NewMetricUpdater(t)
			if !tt.wantErr {
				u.On(tt.mockFn, tt.mockArg).Once().Return(tt.mockResult)
				got := storeMetric(u, tt.Arg)
				assert.Equal(t, got, tt.want)
			} else {
				u.AssertNotCalled(t, tt.mockFn)
				got := storeMetric(u, tt.Arg)
				assert.Equal(t, got, tt.want)
			}
		})
	}

	for _, tt := range counterTests {
		t.Run(tt.name, func(t *testing.T) {
			u := mocks.NewMetricUpdater(t)
			if !tt.wantErr {
				u.On(tt.mockFn, tt.mockArg).Once().Return(tt.mockResult)
				got := storeMetric(u, tt.Arg)
				assert.Equal(t, got, tt.want)
			} else {
				u.AssertNotCalled(t, tt.mockFn)
				got := storeMetric(u, tt.Arg)
				assert.Equal(t, got, tt.want)
			}
		})
	}
}
