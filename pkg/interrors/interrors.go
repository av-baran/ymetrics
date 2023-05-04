package interrors

const (
	ErrInvalidMetricType           = "invalid metric type"
	ErrInvalidMetricValue          = "invalid value"
	ErrMetricExistsWithAnotherType = "metric with same name and different type already exists"
	ErrBadURL                      = "bad request"
	ErrMetricNotFound              = "metric not found"
	ErrStorageInternalError        = "metric storage error"
)
