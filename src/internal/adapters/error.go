package adapters

import (
	"errors"
	"fmt"
)

var ErrAdapter = errors.New("adapter error")

type AdapterError struct {
	Message string
	Cause   error
}

func NewAdapterError(message string, cause error) error {
	return &AdapterError{Message: message, Cause: cause}
}

func (e *AdapterError) Error() string {
	return fmt.Sprintf("%s: %s: %v", ErrAdapter.Error(), e.Message, e.Cause)
}

func (e *AdapterError) Unwrap() error {
	return ErrAdapter
}
