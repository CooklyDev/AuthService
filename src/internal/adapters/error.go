package adapters

import (
	"errors"
	"fmt"
)

var ErrInvariant = errors.New("adapter invariant violation")
var ErrServer = errors.New("adapter server error")

type AdapterError struct {
	Kind    error
	Message string
	Cause   error
}

func NewInvariantError(message string, cause error) error {
	return &AdapterError{
		Kind:    ErrInvariant,
		Message: message,
		Cause:   cause,
	}
}

func NewServerError(message string, cause error) error {
	return &AdapterError{
		Kind:    ErrServer,
		Message: message,
		Cause:   cause,
	}
}

func (e *AdapterError) Error() string {
	return fmt.Sprintf("%s: %s: %v", e.Kind.Error(), e.Message, e.Cause)
}

func (e *AdapterError) Unwrap() error {
	return e.Cause
}

func (e *AdapterError) Is(target error) bool {
	return target == e.Kind
}
