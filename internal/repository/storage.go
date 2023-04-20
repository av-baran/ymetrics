package repository

type Storager interface {
	SetGauge(string, float64)
	AddCounter(string, int64)

	GetGauge(string) string
	GetCounter(string) string

	GetAllGauge() map[string]float64
	GetAllCounter() map[string]int64
}
