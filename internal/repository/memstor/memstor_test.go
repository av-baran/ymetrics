package memstor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.NotEmpty(t, New())
}

func TestMemStorage_SetGetGauge(t *testing.T) {
	type arguments struct {
		name  string
		value float64
	}
	tests := []struct {
		name string
		args arguments
	}{
		{
			name: "test gauge",
			args: arguments{
				name:  "newname",
				value: float64(0.001),
			},
		},
	}
	storage := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.SetGauge(tt.args.name, tt.args.value)
			got, err := storage.GetGauge(tt.args.name)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.value, got)
			got1, err := storage.GetGauge("unknown metric")
			assert.Error(t, err)
			assert.Equal(t, got1, float64(0))
		})
	}
}

func TestMemStorage_SetGetCounter(t *testing.T) {
	type arguments struct {
		name  string
		value int64
	}
	tests := []struct {
		name string
		args arguments
	}{
		{
			name: "test counter",
			args: arguments{
				name:  "newname",
				value: int64(10),
			},
		},
	}
	storage := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.AddCounter(tt.args.name, tt.args.value)
			got, err := storage.GetCounter(tt.args.name)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.value, got)
			got1, err := storage.GetCounter("unknown metric")
			assert.Error(t, err)
			assert.Equal(t, got1, int64(0))
			storage.AddCounter(tt.args.name, 1)
			got2, err := storage.GetCounter(tt.args.name)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.value+1, got2)
		})
	}
}
