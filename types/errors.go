// Package types provides common type definitions for the Stack0 SDK.
package types

import "fmt"

// ErrorResponse represents an error response from the API.
type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// APIError represents an error returned by the Stack0 API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Response   ErrorResponse
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("stack0: %s (code: %s, status: %d)", e.Message, e.Code, e.StatusCode)
	}
	return fmt.Sprintf("stack0: %s (status: %d)", e.Message, e.StatusCode)
}

// TimeoutError represents a timeout error during polling operations.
type TimeoutError struct {
	Message string
}

// Error implements the error interface.
func (e *TimeoutError) Error() string {
	return fmt.Sprintf("stack0: timeout: %s", e.Message)
}

// NewTimeoutError creates a new TimeoutError with the given message.
func NewTimeoutError(message string) *TimeoutError {
	return &TimeoutError{Message: message}
}
