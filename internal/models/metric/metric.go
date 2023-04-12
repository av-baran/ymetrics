package metric

type Type string

const (
	GaugeType   = Type("gauge")
	CounterType = Type("counter")
)

type Rawdata struct {
	Type  Type
	Name  string
	Value string
}

type Gauge struct {
	Value float64
}
