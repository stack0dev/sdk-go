package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError_Error_WithCode(t *testing.T) {
	err := &APIError{
		StatusCode: 400,
		Code:       "INVALID_REQUEST",
		Message:    "Invalid request body",
		Response: ErrorResponse{
			Message: "Invalid request body",
			Code:    "INVALID_REQUEST",
		},
	}

	expected := "stack0: Invalid request body (code: INVALID_REQUEST, status: 400)"
	assert.Equal(t, expected, err.Error())
}

func TestAPIError_Error_WithoutCode(t *testing.T) {
	err := &APIError{
		StatusCode: 500,
		Code:       "",
		Message:    "Internal server error",
		Response: ErrorResponse{
			Message: "Internal server error",
		},
	}

	expected := "stack0: Internal server error (status: 500)"
	assert.Equal(t, expected, err.Error())
}

func TestAPIError_ImplementsError(t *testing.T) {
	var err error = &APIError{
		StatusCode: 404,
		Message:    "Not found",
	}

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Not found")
}

func TestTimeoutError_Error(t *testing.T) {
	err := &TimeoutError{
		Message: "Screenshot timed out",
	}

	expected := "stack0: timeout: Screenshot timed out"
	assert.Equal(t, expected, err.Error())
}

func TestTimeoutError_ImplementsError(t *testing.T) {
	var err error = &TimeoutError{
		Message: "Operation timed out",
	}

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "timeout")
}

func TestNewTimeoutError(t *testing.T) {
	err := NewTimeoutError("Screenshot timed out")

	assert.NotNil(t, err)
	assert.Equal(t, "Screenshot timed out", err.Message)
	assert.Equal(t, "stack0: timeout: Screenshot timed out", err.Error())
}

func TestErrorResponse_Fields(t *testing.T) {
	resp := ErrorResponse{
		Message: "Bad request",
		Code:    "BAD_REQUEST",
	}

	assert.Equal(t, "Bad request", resp.Message)
	assert.Equal(t, "BAD_REQUEST", resp.Code)
}

func TestAPIError_AllFields(t *testing.T) {
	err := &APIError{
		StatusCode: 422,
		Code:       "VALIDATION_ERROR",
		Message:    "Email is required",
		Response: ErrorResponse{
			Message: "Email is required",
			Code:    "VALIDATION_ERROR",
		},
	}

	assert.Equal(t, 422, err.StatusCode)
	assert.Equal(t, "VALIDATION_ERROR", err.Code)
	assert.Equal(t, "Email is required", err.Message)
	assert.Equal(t, "Email is required", err.Response.Message)
	assert.Equal(t, "VALIDATION_ERROR", err.Response.Code)
}
