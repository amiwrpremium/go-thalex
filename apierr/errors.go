// Package apierr defines error types returned by the Thalex SDK.
//
// All error types implement the standard error interface.
// [ConnectionError], [AuthError], and [TimeoutError] also implement
// the Unwrap interface for use with [errors.As] and [errors.Is].
package apierr

import (
	"errors"
	"fmt"
)

// APIError represents an error returned by the Thalex API.
type APIError struct {
	// Code is the numeric error code from the API.
	Code int `json:"code"`
	// Message is the human-readable error description.
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("thalex: API error %d: %s", e.Code, e.Message)
}

// ConnectionError represents a connection-level error.
type ConnectionError struct {
	// Message describes the connection error.
	Message string
	// Err is the underlying error, if any.
	Err error
}

// Error implements the error interface.
func (e *ConnectionError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("thalex: connection error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("thalex: connection error: %s", e.Message)
}

// Unwrap returns the underlying error.
func (e *ConnectionError) Unwrap() error {
	return e.Err
}

// AuthError represents an authentication error.
type AuthError struct {
	// Message describes the authentication error.
	Message string
	// Err is the underlying error, if any.
	Err error
}

// Error implements the error interface.
func (e *AuthError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("thalex: auth error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("thalex: auth error: %s", e.Message)
}

// Unwrap returns the underlying error.
func (e *AuthError) Unwrap() error {
	return e.Err
}

// TimeoutError represents a request timeout error.
type TimeoutError struct {
	// Message describes what timed out.
	Message string
	// Err is the underlying error, if any.
	Err error
}

// Error implements the error interface.
func (e *TimeoutError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("thalex: timeout: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("thalex: timeout: %s", e.Message)
}

// Unwrap returns the underlying error.
func (e *TimeoutError) Unwrap() error {
	return e.Err
}

// IsAPIError checks if an error is an APIError and returns it.
func IsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}
