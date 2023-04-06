package httperror

import "fmt"

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
