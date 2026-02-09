package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stack0/sdk-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))
		assert.Equal(t, "/test-path", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "success"})
	}))
	defer server.Close()

	client := New("test-api-key", server.URL)

	var result map[string]string
	err := client.Get(context.Background(), "/test-path", &result)

	require.NoError(t, err)
	assert.Equal(t, "success", result["message"])
}

func TestHTTPClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "test-value", body["key"])

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "123"})
	}))
	defer server.Close()

	client := New("test-api-key", server.URL)

	var result map[string]string
	err := client.Post(context.Background(), "/test-path", map[string]string{"key": "test-value"}, &result)

	require.NoError(t, err)
	assert.Equal(t, "123", result["id"])
}

func TestHTTPClient_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.ErrorResponse{
			Message: "Invalid request",
			Code:    "INVALID_REQUEST",
		})
	}))
	defer server.Close()

	client := New("test-api-key", server.URL)

	var result map[string]string
	err := client.Get(context.Background(), "/test-path", &result)

	require.Error(t, err)
	apiErr, ok := err.(*types.APIError)
	require.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	assert.Equal(t, "INVALID_REQUEST", apiErr.Code)
	assert.Equal(t, "Invalid request", apiErr.Message)
}

func TestHTTPClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(types.SuccessResponse{Success: true})
	}))
	defer server.Close()

	client := New("test-api-key", server.URL)

	var result types.SuccessResponse
	err := client.Delete(context.Background(), "/test-path", &result)

	require.NoError(t, err)
	assert.True(t, result.Success)
}

func TestHTTPClient_DeleteWithBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "123", body["id"])

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(types.SuccessResponse{Success: true})
	}))
	defer server.Close()

	client := New("test-api-key", server.URL)

	var result types.SuccessResponse
	err := client.DeleteWithBody(context.Background(), "/test-path", map[string]string{"id": "123"}, &result)

	require.NoError(t, err)
	assert.True(t, result.Success)
}
