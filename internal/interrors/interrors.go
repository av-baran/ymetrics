package interrors

import "fmt"

const (
	ErrInvalidMetricType   = "invalid metric type"
	ErrInvalidMetricValue  = "invalid value"
	ErrMetricAlreadyExists = "metric with same name and different type already exists"
	ErrBadURL              = "bad request"
	ErrMetricNotFound      = "metric not found"
)

func New(msg string, code int) *Error {
	return &Error{"string", code}
}

type Error struct {
	Msg  string
	Code int
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code: %v, Message: %v", e.Code, e.Msg)
}
