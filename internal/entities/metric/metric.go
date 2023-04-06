package metric

type Type string

const (
	Gauge   = Type("gauge")
	Counter = Type("counter")
)

type Rawdata struct {
	Type  Type
	Name  string
	Value string
}
