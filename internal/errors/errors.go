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

type DatabaseError struct {
	Msg string
	Err error
}

func (e DatabaseError) Error() string { return e.Msg }
func (e DatabaseError) Unwrap() error { return e.Err }

func NewDatabaseError(detail error, format string, a ...interface{}) error {
	args := append(a, detail)

	return DatabaseError{
		Msg: fmt.Sprintf(format+": %v", args...),
		Err: detail,
	}
}

type InvalidSourceError struct {
	Msg string
	Err error
}

func (e InvalidSourceError) Error() string { return e.Msg }
func (e InvalidSourceError) Unwrap() error { return e.Err }

func NewInvalidSourceError(detail error, format string, a ...interface{}) error {
	args := append(a, detail)

	return InvalidSourceError{
		Msg: fmt.Sprintf(format+": %v", args...),
		Err: detail,
	}
}

type InvalidVersionError struct {
	Msg string
	Err error
}

func (e InvalidVersionError) Error() string { return e.Msg }
func (e InvalidVersionError) Unwrap() error { return e.Err }

func NewInvalidVersionError(detail error, format string, a ...interface{}) error {
	args := append(a, detail)

	return InvalidVersionError{
		Msg: fmt.Sprintf(format+": %v", args...),
		Err: detail,
	}
}
