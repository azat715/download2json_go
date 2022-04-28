package core

import (
	"errors"
)

type ErrorWrapper struct {
	Message string
	Context string
	Err     error
}

func (err ErrorWrapper) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return err.Message
}

func (err ErrorWrapper) Unwrap() error {
	return err.Err // Returns inner error
}

// Returns the inner most CustomErrorWrapper
func (err ErrorWrapper) Dig() ErrorWrapper {
	var ew ErrorWrapper
	if errors.As(err.Err, &ew) {
		// Recursively digs until wrapper error is not in which case it will stop
		return ew.Dig()
	}
	return err
}

func NewErrorWrapper(context string, err error, message string) error {
	return ErrorWrapper{
		Message: message,
		Context: context,
		Err:     err,
	}
}
