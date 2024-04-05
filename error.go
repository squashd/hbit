package hbit

import (
	"errors"
	"fmt"
	"strings"
)

type AppError string

const (
	ECONFLICT       AppError = "conflict"
	EINTERNAL       AppError = "internal"
	EINVALID        AppError = "invalid"
	ENOTFOUND       AppError = "not_found"
	ENOTIMPLEMENTED AppError = "not_implemented"
	EUNAUTHORIZED   AppError = "unauthorized"
	EFORBIDDEN      AppError = "forbidden"
	EASYNC          AppError = "async"
)

type Error struct {
	Code    AppError
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("hbit error: code=%s message=%s", e.Code, e.Message)
}

func ErrorCode(err error) AppError {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}
	return EINTERNAL
}

func ErrorMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}
	return "Internal error."
}

func Errorf(code AppError, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Join joins multiple errors into a single error.
func NewMultiError(errors ...*Error) error {
	return &MultiError{Errors: errors}
}

// MultiError is an error that contains multiple errors.
// Intended to be used for validation errors.
type MultiError struct {
	Errors []*Error
}

// MultiError implements the error interface
// Not intended to be used directly.
func (m *MultiError) Error() string {
	var messages []string
	for _, err := range m.Errors {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("multiple errors: %s", strings.Join(messages, ", "))
}

func (m *MultiError) GetErrorMessages() []string {
	var messages []string
	for _, err := range m.Errors {
		messages = append(messages, err.Message)
	}
	return messages
}
