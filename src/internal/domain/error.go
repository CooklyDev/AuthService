package domain

import "errors"

var ErrBusinessRule = errors.New("business rule violation")

type BusinessRuleError struct {
	Message string
}

func NewBusinessRuleError(message string) error {
	return &BusinessRuleError{Message: message}
}

func (e *BusinessRuleError) Error() string {
	return ErrBusinessRule.Error() + ": " + e.Message
}

func (e *BusinessRuleError) Unwrap() error {
	return ErrBusinessRule
}
