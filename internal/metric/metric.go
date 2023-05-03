package metric

type Type string

const (
	GaugeType   = "gauge"
	CounterType = "counter"
	UnknownType = "unknown"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
