package interrors

import "fmt"

const (
	ErrInvalidMetricType   = "invalid metric type"
	ErrInvalidMetricValue  = "invalid value"
	ErrMetricAlreadyExists = "metric with same name and different type already exists"
	ErrBadUrl              = "bad request"
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
