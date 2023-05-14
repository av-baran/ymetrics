package memstor

import (
	"testing"

	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/stretchr/testify/assert"
)

const (
	gaugeVal   = float64(0.1)
	counterVal = int64(1)
	unknownVal = uint32(1)
)

//FIXME как проще получить указатель на значение для поля структуры
func getFloatPtr(v float64) *float64 {
	return &v
}

func getIntPtr(v int64) *int64 {
	return &v
}

func TestNew(t *testing.T) {
	assert.NotEmpty(t, New())
}

func TestMemStorage(t *testing.T) {
	tSetMetric := []struct {
		name    string
		metric  metric.Metric
		wantErr bool
	}{
		{
			name: "set gauge",
			metric: metric.Metric{
				ID:    "someGauge",
				MType: metric.GaugeType,
				Value: getFloatPtr(gaugeVal),
				Delta: nil,
			},
			wantErr: false,
		},
		{
			name: "set counter",
			metric: metric.Metric{
				ID:    "someCounter",
				MType: metric.CounterType,
				Value: nil,
				Delta: getIntPtr(counterVal),
			},
			wantErr: false,
		},
		{
			name: "set unknown type",
			metric: metric.Metric{
				ID:    "unknown",
				MType: metric.UnknownType,
				Value: getFloatPtr(gaugeVal),
				Delta: nil,
			},
			wantErr: true,
		},
	}

	tGetMetric := []struct {
		name      string
		metric    metric.Metric
		wantErr   bool
		wantValue *float64
		wantDelta *int64
	}{
		{
			name: "get gauge",
			metric: metric.Metric{
				ID:    "someGauge",
				MType: metric.GaugeType,
				Value: nil,
				Delta: nil,
			},
			wantErr:   false,
			wantValue: getFloatPtr(gaugeVal),
			wantDelta: nil,
		},
		{
			name: "get counter",
			metric: metric.Metric{
				ID:    "someCounter",
				MType: metric.CounterType,
				Delta: nil,
				Value: nil,
			},
			wantErr:   false,
			wantValue: nil,
			wantDelta: getIntPtr(counterVal),
		},
		{
			name: "get unknown type",
			metric: metric.Metric{
				ID:    "someGauge",
				MType: metric.UnknownType,
				Value: nil,
				Delta: nil,
			},
			wantErr:   true,
			wantValue: nil,
			wantDelta: nil,
		},
		{
			name: "get unknown name",
			metric: metric.Metric{
				ID:    "unknown",
				MType: metric.CounterType,
				Value: nil,
				Delta: nil,
			},
			wantErr:   true,
			wantValue: nil,
			wantDelta: nil,
		},
	}

	tAllMetric := make([]metric.Metric, 0)
	for _, v := range tSetMetric {
		if !v.wantErr {
			tAllMetric = append(tAllMetric, v.metric)
		}
	}

	s := New()
	for _, tt := range tSetMetric {
		t.Run(tt.name, func(t *testing.T) {
			err := s.SetMetric(tt.metric)
			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}

	for _, tt := range tGetMetric {
		t.Run(tt.name, func(t *testing.T) {
			res, err := s.GetMetric(tt.metric.ID, tt.metric.MType)
			if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, res.Delta, tt.wantDelta)
				assert.Equal(t, res.Value, tt.wantValue)
			} else {
				assert.Error(t, err)
			}
		})
	}

	gotMetrics := s.GetAllMetrics()
	assert.ObjectsAreEqual(tAllMetric, gotMetrics)
}
