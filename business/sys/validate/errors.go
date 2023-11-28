package validate

import (
	"encoding/json"
	"errors"
)

// ErrInvalidID occurs when ID is not in a valid form
var ErrInvalidID = errors.New("ID is not in its proper form")

// ErrorResponce is the form used for API responses from failures in the API
type ErrorResponse struct {
	Error  string `json:"error"`
	Fields string `json:"fields,omitempty"`
}

// RequestError is used to pass an error during the request though the
// application with web specific context
type RequestError struct {
	Err    error
	Status int
	Fields error
}

// NewRequestError wrpas a provided error with an HTTP status code.
// This function should be used when handlerrs encounter expected errors.
func NewRequestError(err error, status int) error {
	return &RequestError{err, status, nil}
}

// Error implements the error inteface. It uses the defaule message of the wrapped
// error. This is what will be shown in the services' logs.
func (err *RequestError) Error() string {
	return err.Err.Error()
}

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// FieldErrors represents a collection of field errors.
type FieldErrors []FieldError

// Error implments the error interface.
func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

// Cause iterates through all the wrapped errors until the root
// error value is reached.
func Cause(err error) error {
	root := err
	for {
		if err = errors.Unwrap(root); err == nil {
			return root
		}
		root = err
	}
}
