package service

import (
	"reflect"
	"testing"

	"github.com/av-baran/ymetrics/internal/entity/metric"
	"github.com/av-baran/ymetrics/internal/repository/memstorv2"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	s := memstorv2.New()
	assert.NotEmpty(t, New(s))
}

func TestService_UpdateGauge(t *testing.T) {
	testStorage := &memstorv2.MemStorage{
		GaugeStor: map[string]float64{
			"name1": 0.01,
			"name2": 0.02,
			"name3": 0.03,
		},
		CounterStor: map[string]int64{
			"name4": 1,
			"name5": 2,
			"name6": 3,
		},
	}

	tests := []struct {
		name    string
		arg     *metric.Gauge
		wantErr bool
	}{
		{
			name: "test1",
			arg: &metric.Gauge{
				Name:  "name7",
				Value: float64(0.04),
			},
			wantErr: false,
		},
		{
			name: "test2",
			arg: &metric.Gauge{
				Name:  "name6",
				Value: float64(0.04),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Storage: testStorage,
			}
			if err := s.UpdateGauge(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_UpdateCounter(t *testing.T) {
	testStorage := &memstorv2.MemStorage{
		GaugeStor: map[string]float64{
			"name1": 0.01,
			"name2": 0.02,
			"name3": 0.03,
		},
		CounterStor: map[string]int64{
			"name4": 1,
			"name5": 2,
			"name6": 3,
		},
	}

	tests := []struct {
		name    string
		arg     *metric.Counter
		wantErr bool
	}{
		{
			name: "test1",
			arg: &metric.Counter{
				Name:  "name7",
				Value: int64(4),
			},
			wantErr: false,
		},
		{
			name: "test2",
			arg: &metric.Counter{
				Name:  "name7",
				Value: int64(4),
			},
			wantErr: false,
		},
		{
			name: "test3",
			arg: &metric.Counter{
				Name:  "name1",
				Value: int64(4),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Storage: testStorage,
			}
			if err := s.UpdateCounter(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetGauge(t *testing.T) {
	testStorage := &memstorv2.MemStorage{
		GaugeStor: map[string]float64{
			"name1": 0.01,
			"name2": 0.02,
			"name3": 0.03,
		},
		CounterStor: map[string]int64{
			"name4": 1,
			"name5": 2,
			"name6": 3,
		},
	}

	tests := []struct {
		name    string
		arg     string
		want    float64
		wantErr bool
	}{
		{
			name:    "test1",
			arg:     "name1",
			want:    float64(0.01),
			wantErr: false,
		},
		{
			name:    "test2",
			arg:     "unknown_name",
			want:    float64(0.0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Storage: testStorage,
			}
			got, err := s.GetGauge(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.GetGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetCounter(t *testing.T) {
	testStorage := &memstorv2.MemStorage{
		GaugeStor: map[string]float64{
			"name1": 0.01,
			"name2": 0.02,
			"name3": 0.03,
		},
		CounterStor: map[string]int64{
			"name4": 1,
			"name5": 2,
			"name6": 3,
		},
	}

	tests := []struct {
		name    string
		arg     string
		want    int64
		wantErr bool
	}{
		{
			name:    "test1",
			arg:     "name4",
			want:    int64(1),
			wantErr: false,
		},
		{
			name:    "test2",
			arg:     "unknown_name",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s := &Service{
			Storage: testStorage,
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetCounter(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.GetCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetAllMetrics(t *testing.T) {
	tests := []struct {
		name    string
		storage Storager
		want    map[string]interface{}
	}{
		{
			name: "test1",
			storage: &memstorv2.MemStorage{
				GaugeStor: map[string]float64{
					"name1": 0.01,
					"name2": 0.02,
					"name3": 0.03,
				},
				CounterStor: map[string]int64{
					"name4": 1,
					"name5": 2,
					"name6": 3,
				},
			},
			want: map[string]interface{}{
				"name1": float64(0.01),
				"name2": float64(0.02),
				"name3": float64(0.03),
				"name4": int64(1),
				"name5": int64(2),
				"name6": int64(3),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Storage: tt.storage,
			}
			if got := s.GetAllMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetAllMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
