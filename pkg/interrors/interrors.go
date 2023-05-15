package interrors

import "errors"

var (
	ErrInvalidMetricType           = errors.New("invalid metric type")
	ErrInvalidMetricValue          = errors.New("invalid value")
	ErrMetricExistsWithAnotherType = errors.New("metric with same name and different type already exists")
	ErrBadURL                      = errors.New("bad request")
	ErrMetricNotFound              = errors.New("metric not found")
	ErrStorageInternalError        = errors.New("metric storage error")
	ErrPingDB                      = errors.New("DB connection test failed")
)
