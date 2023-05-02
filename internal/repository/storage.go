package repository

type Storager interface {
	SetGauge(string, float64)
	AddCounter(string, int64)

	GetGauge(string) (float64, error)
	GetCounter(string) (int64, error)

	GetAllGauge() map[string]float64
	GetAllCounter() map[string]int64
}
