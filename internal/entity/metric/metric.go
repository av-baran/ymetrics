package metric

type Type string

const (
	GaugeType   = Type("gauge")
	CounterType = Type("counter")
	UnknownType = Type("unknown")
)

type Rawdata struct {
	Type  Type
	Name  string
	Value string
}

type Gauge struct {
	Name  string
	Value float64
}

type Counter struct {
	Name  string
	Value int64
}
