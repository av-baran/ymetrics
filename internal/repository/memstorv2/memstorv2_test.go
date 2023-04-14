package memstorv2

import (
	"testing"

	"github.com/av-baran/ymetrics/internal/entity/metric"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.NotEmpty(t, New())
}

func TestMemStorage_StoreMetric(t *testing.T) {
	tests := []struct {
		name    string
		arg     interface{}
		wantErr bool
	}{
		{
			name: "test gauge",
			arg: &metric.Gauge{
				Name:  "newname",
				Value: float64(0.001),
			},
			wantErr: false,
		},
		{
			name: "test counter",
			arg: &metric.Counter{
				Name:  "newname",
				Value: int64(1),
			},
			wantErr: false,
		},
		{
			name:    "test unknown",
			arg:     metric.UnknownType,
			wantErr: true,
		},
	}
	storage := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.StoreMetric(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.StoreMetric() error = %v, wantErr %v", err, tt.wantErr)
			}

			switch m := tt.arg.(type) {
			case *metric.Gauge:
				assert.Equal(t, storage.gaugeStor[m.Name], m.Value)
			case *metric.Counter:
				assert.Equal(t, storage.counterStor[m.Name], m.Value)
			default:
				assert.Error(t, err)

			}
		})
	}
}

func TestMemStorage_GetMetricType(t *testing.T) {
	tests := []struct {
		name    string
		storage *MemStorage
		arg     string
		wantErr bool
		wantRes metric.Type
	}{
		{
			name: "test gauge",
			storage: &MemStorage{
				gaugeStor:   map[string]float64{"name": 1.01},
				counterStor: map[string]int64{"another_name": 10},
			},
			arg:     "name",
			wantRes: metric.GaugeType,
			wantErr: false,
		},
		{
			name: "test counter",
			storage: &MemStorage{
				gaugeStor:   map[string]float64{"name": 1.01},
				counterStor: map[string]int64{"another_name": 10},
			},
			arg:     "another_name",
			wantRes: metric.CounterType,
			wantErr: false,
		},
		{
			name: "test unknown",
			storage: &MemStorage{
				gaugeStor:   map[string]float64{"name": 1.01},
				counterStor: map[string]int64{"another_name": 10},
			},
			arg:     "unknown_name",
			wantRes: metric.UnknownType,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.storage.GetMetricType(tt.arg)
			if got1 != !tt.wantErr {
				t.Errorf("MemStorage.GetMetric() ok = %v, wantErr %v", got1, tt.wantErr)
			}
			assert.Equal(t, got, tt.wantRes)
		})
	}
}

func TestMemStorage_GetMetric(t *testing.T) {
	tests := []struct {
		name    string
		storage *MemStorage
		arg     string
		wantErr bool
		wantRes interface{}
	}{
		{
			name: "test gauge",
			storage: &MemStorage{
				gaugeStor:   map[string]float64{"name": 1.01},
				counterStor: map[string]int64{"another_name": 10},
			},
			arg:     "name",
			wantRes: float64(1.01),
			wantErr: false,
		},
		{
			name: "test counter",
			storage: &MemStorage{
				gaugeStor:   map[string]float64{"name": 1.01},
				counterStor: map[string]int64{"another_name": 10},
			},
			arg:     "another_name",
			wantRes: int64(10),
			wantErr: false,
		},
		{
			name: "test unknown",
			storage: &MemStorage{
				gaugeStor:   map[string]float64{"name": 1.01},
				counterStor: map[string]int64{"another_name": 10},
			},
			arg:     "unknown_name",
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.storage.GetMetric(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.GetMetric() ok = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, got, tt.wantRes)
		})
	}
}
