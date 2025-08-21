package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		code       string
		message    string
		expected   *APIError
	}{
		{
			name:       "Creates error with all fields",
			statusCode: http.StatusBadRequest,
			code:       "BAD_REQUEST",
			message:    "Invalid input",
			expected: &APIError{
				StatusCode: http.StatusBadRequest,
				Code:       "BAD_REQUEST",
				Message:    "Invalid input",
			},
		},
		{
			name:       "Creates error with empty message",
			statusCode: http.StatusInternalServerError,
			code:       "INTERNAL_ERROR",
			message:    "",
			expected: &APIError{
				StatusCode: http.StatusInternalServerError,
				Code:       "INTERNAL_ERROR",
				Message:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIError(tt.statusCode, tt.code, tt.message)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestNewAPIErrorWithDetails(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		code       string
		message    string
		details    string
		expected   *APIError
	}{
		{
			name:       "Creates error with details",
			statusCode: http.StatusBadRequest,
			code:       "VALIDATION_ERROR",
			message:    "Validation failed",
			details:    "Email is invalid",
			expected: &APIError{
				StatusCode: http.StatusBadRequest,
				Code:       "VALIDATION_ERROR",
				Message:    "Validation failed",
				Details:    "Email is invalid",
			},
		},
		{
			name:       "Creates error without details",
			statusCode: http.StatusNotFound,
			code:       "NOT_FOUND",
			message:    "Resource not found",
			details:    "",
			expected: &APIError{
				StatusCode: http.StatusNotFound,
				Code:       "NOT_FOUND",
				Message:    "Resource not found",
				Details:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIErrorWithDetails(tt.statusCode, tt.code, tt.message, tt.details)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		apiError *APIError
		expected string
	}{
		{
			name: "Error message without details",
			apiError: &APIError{
				StatusCode: http.StatusBadRequest,
				Code:       "BAD_REQUEST",
				Message:    "Invalid request",
			},
			expected: "BAD_REQUEST: Invalid request",
		},
		{
			name: "Error message with details",
			apiError: &APIError{
				StatusCode: http.StatusBadRequest,
				Code:       "VALIDATION_ERROR",
				Message:    "Validation failed",
				Details:    "Email format is invalid",
			},
			expected: "VALIDATION_ERROR: Validation failed (details: Email format is invalid)",
		},
		{
			name: "Error message with empty details",
			apiError: &APIError{
				StatusCode: http.StatusInternalServerError,
				Code:       "INTERNAL_ERROR",
				Message:    "Something went wrong",
				Details:    "",
			},
			expected: "INTERNAL_ERROR: Something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.apiError.Error()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewErrorResponse(t *testing.T) {
	apiErr := &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Message:    "Invalid input",
		Details:    "Missing required field",
	}

	response := NewErrorResponse(apiErr)

	assert.NotNil(t, response)
	assert.Equal(t, apiErr, response.Error)
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        *APIError
		statusCode int
		code       string
		message    string
	}{
		{
			name:       "Bad Request Error",
			err:        ErrBadRequest,
			statusCode: http.StatusBadRequest,
			code:       "BAD_REQUEST",
			message:    "Invalid request parameters",
		},
		{
			name:       "Unauthorized Error",
			err:        ErrUnauthorized,
			statusCode: http.StatusUnauthorized,
			code:       "UNAUTHORIZED",
			message:    "Authentication required",
		},
		{
			name:       "Forbidden Error",
			err:        ErrForbidden,
			statusCode: http.StatusForbidden,
			code:       "FORBIDDEN",
			message:    "Access denied",
		},
		{
			name:       "Not Found Error",
			err:        ErrNotFound,
			statusCode: http.StatusNotFound,
			code:       "NOT_FOUND",
			message:    "Resource not found",
		},
		{
			name:       "Conflict Error",
			err:        ErrConflict,
			statusCode: http.StatusConflict,
			code:       "CONFLICT",
			message:    "Resource already exists",
		},
		{
			name:       "Internal Server Error",
			err:        ErrInternalServer,
			statusCode: http.StatusInternalServerError,
			code:       "INTERNAL_SERVER_ERROR",
			message:    "An internal error occurred",
		},
		{
			name:       "Service Unavailable Error",
			err:        ErrServiceUnavailable,
			statusCode: http.StatusServiceUnavailable,
			code:       "SERVICE_UNAVAILABLE",
			message:    "Service temporarily unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.statusCode, tt.err.StatusCode)
			assert.Equal(t, tt.code, tt.err.Code)
			assert.Equal(t, tt.message, tt.err.Message)
			assert.Empty(t, tt.err.Details)
		})
	}
}

func TestAPIError_ImplementsError(t *testing.T) {
	var err error = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "TEST",
		Message:    "Test error",
	}

	// Should be able to use as error interface
	assert.NotNil(t, err)
	assert.Equal(t, "TEST: Test error", err.Error())
}
