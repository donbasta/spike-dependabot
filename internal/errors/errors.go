package errors

import "fmt"

type OperationError struct {
	Msg string
	Err error
}

func (e OperationError) Error() string { return e.Msg }
func (e OperationError) Unwrap() error { return e.Err }

func NewOperationError(detail error, format string, a ...interface{}) error {
	args := append(a, detail)

	return OperationError{
		Msg: fmt.Sprintf(format+": %v", args...),
		Err: detail,
	}
}
