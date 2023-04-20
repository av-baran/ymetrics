package repository

type Storager interface {
	SetGauge(string, float64)
	AddCounter(string, int64)

	GetGauge(string) (string, error)
	GetCounter(string) (string, error)

	GetAllGauge() map[string]float64
	GetAllCounter() map[string]int64
}
