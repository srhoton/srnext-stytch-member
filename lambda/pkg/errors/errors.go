package errors

import (
	"fmt"
	"net/http"
)

// APIError represents an API error with status code and message
type APIError struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (details: %s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAPIError creates a new API error
func NewAPIError(statusCode int, code, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

// NewAPIErrorWithDetails creates a new API error with details
func NewAPIErrorWithDetails(statusCode int, code, message, details string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Details:    details,
	}
}

// Common error types
var (
	ErrBadRequest = NewAPIError(
		http.StatusBadRequest,
		"BAD_REQUEST",
		"Invalid request parameters",
	)

	ErrUnauthorized = NewAPIError(
		http.StatusUnauthorized,
		"UNAUTHORIZED",
		"Authentication required",
	)

	ErrForbidden = NewAPIError(
		http.StatusForbidden,
		"FORBIDDEN",
		"Access denied",
	)

	ErrNotFound = NewAPIError(
		http.StatusNotFound,
		"NOT_FOUND",
		"Resource not found",
	)

	ErrConflict = NewAPIError(
		http.StatusConflict,
		"CONFLICT",
		"Resource already exists",
	)

	ErrInternalServer = NewAPIError(
		http.StatusInternalServerError,
		"INTERNAL_SERVER_ERROR",
		"An internal error occurred",
	)

	ErrServiceUnavailable = NewAPIError(
		http.StatusServiceUnavailable,
		"SERVICE_UNAVAILABLE",
		"Service temporarily unavailable",
	)
)

// ErrorResponse represents the error response structure
type ErrorResponse struct {
	Error *APIError `json:"error"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err *APIError) *ErrorResponse {
	return &ErrorResponse{Error: err}
}
