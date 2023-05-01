package metric

type Type string

const (
	GaugeType   = Type("gauge")
	CounterType = Type("counter")
	UnknownType = Type("unknown")
)

type Metric struct {
	Name  string
	Value interface{}
	Type  Type
}

type Gauge struct {
	Name  string
	Value float64
}

type Counter struct {
	Name  string
	Value int64
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
