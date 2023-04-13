package memstorv2

// import (
// 	"reflect"
// 	"testing"
//
// 	"github.com/av-baran/ymetrics/internal/entities/metric"
// 	"github.com/av-baran/ymetrics/internal/httperror"
// )

// func TestNew(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want *MemStorage
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := New(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("New() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestMemStorage_UpdateMetric(t *testing.T) {
// 	type fields struct {
// 		metrics     map[string]metric.Type
// 		gaugeStor   map[string]float64
// 		counterStor map[string]int64
// 	}
// 	type args struct {
// 		m *metric.Rawdata
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   *httperror.Error
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &MemStorage{
// 				metrics:     tt.fields.metrics,
// 				gaugeStor:   tt.fields.gaugeStor,
// 				counterStor: tt.fields.counterStor,
// 			}
// 			if got := s.UpdateMetric(tt.args.m); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("MemStorage.UpdateMetric() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestMemStorage_updateGauge(t *testing.T) {
// 	type fields struct {
// 		metrics     map[string]metric.Type
// 		gaugeStor   map[string]float64
// 		counterStor map[string]int64
// 	}
// 	type args struct {
// 		m *metric.Rawdata
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   *httperror.Error
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &MemStorage{
// 				metrics:     tt.fields.metrics,
// 				gaugeStor:   tt.fields.gaugeStor,
// 				counterStor: tt.fields.counterStor,
// 			}
// 			if got := s.updateGauge(tt.args.m); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("MemStorage.updateGauge() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestMemStorage_updateCounter(t *testing.T) {
// 	type fields struct {
// 		metrics     map[string]metric.Type
// 		gaugeStor   map[string]float64
// 		counterStor map[string]int64
// 	}
// 	type args struct {
// 		m *metric.Rawdata
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   *httperror.Error
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &MemStorage{
// 				metrics:     tt.fields.metrics,
// 				gaugeStor:   tt.fields.gaugeStor,
// 				counterStor: tt.fields.counterStor,
// 			}
// 			if got := s.updateCounter(tt.args.m); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("MemStorage.updateCounter() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestMemStorage_addMetric(t *testing.T) {
// 	type fields struct {
// 		metrics     map[string]metric.Type
// 		gaugeStor   map[string]float64
// 		counterStor map[string]int64
// 	}
// 	type args struct {
// 		m *metric.Rawdata
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   *httperror.Error
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &MemStorage{
// 				metrics:     tt.fields.metrics,
// 				gaugeStor:   tt.fields.gaugeStor,
// 				counterStor: tt.fields.counterStor,
// 			}
// 			if got := s.addMetric(tt.args.m); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("MemStorage.addMetric() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
