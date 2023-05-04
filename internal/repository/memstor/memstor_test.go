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
		metric  metric.Metrics
		wantErr bool
	}{
		{
			name: "set gauge",
			metric: metric.Metrics{
				ID:    "someGauge",
				MType: metric.GaugeType,
				Value: getFloatPtr(gaugeVal),
				Delta: nil,
			},
			wantErr: false,
		},
		{
			name: "set counter",
			metric: metric.Metrics{
				ID:    "someCounter",
				MType: metric.CounterType,
				Value: nil,
				Delta: getIntPtr(counterVal),
			},
			wantErr: false,
		},
		{
			name: "set unknown type",
			metric: metric.Metrics{
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
		metric    metric.Metrics
		wantErr   bool
		wantValue *float64
		wantDelta *int64
	}{
		{
			name: "get gauge",
			metric: metric.Metrics{
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
			metric: metric.Metrics{
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
			metric: metric.Metrics{
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
			metric: metric.Metrics{
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

	tAllMetric := make([]metric.Metrics, 0)
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
			err := s.GetMetric(&tt.metric)
			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.metric.Delta, tt.wantDelta)
			assert.Equal(t, tt.metric.Value, tt.wantValue)
		})
	}

	gotMetrics := s.GetAllMetrics()
	assert.Equal(t, tAllMetric, gotMetrics)
}

func TestMemStorage_AddCounter(t *testing.T) {
	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "first counter",
			args: args{
				name:  "c",
				value: 1,
			},
			want: 1,
		},
		{
			name: "second counter",
			args: args{
				name:  "c",
				value: 2,
			},
			want: 3,
		},
		{
			name: "third counter",
			args: args{
				name:  "c",
				value: 4,
			},
			want: 7,
		},
	}
	s := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, s.AddCounter(tt.args.name, tt.args.value))
		})
	}
}
