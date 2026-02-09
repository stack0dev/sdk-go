// Package client provides the HTTP client for the Stack0 SDK.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/stack0/sdk-go/types"
)

// HTTPClient handles HTTP communication with the Stack0 API.
type HTTPClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// New creates a new HTTP client.
func New(apiKey, baseURL string) *HTTPClient {
	return &HTTPClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs an HTTP request.
func (c *HTTPClient) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp types.ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			errResp.Message = string(respBody)
		}
		return nil, &types.APIError{
			StatusCode: resp.StatusCode,
			Code:       errResp.Code,
			Message:    errResp.Message,
			Response:   errResp,
		}
	}

	return respBody, nil
}

// Get performs a GET request.
func (c *HTTPClient) Get(ctx context.Context, path string, result interface{}) error {
	respBody, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	if result != nil && len(respBody) > 0 {
		return json.Unmarshal(respBody, result)
	}
	return nil
}

// Post performs a POST request.
func (c *HTTPClient) Post(ctx context.Context, path string, body, result interface{}) error {
	respBody, err := c.doRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	if result != nil && len(respBody) > 0 {
		return json.Unmarshal(respBody, result)
	}
	return nil
}

// Put performs a PUT request.
func (c *HTTPClient) Put(ctx context.Context, path string, body, result interface{}) error {
	respBody, err := c.doRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return err
	}
	if result != nil && len(respBody) > 0 {
		return json.Unmarshal(respBody, result)
	}
	return nil
}

// Patch performs a PATCH request.
func (c *HTTPClient) Patch(ctx context.Context, path string, body, result interface{}) error {
	respBody, err := c.doRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return err
	}
	if result != nil && len(respBody) > 0 {
		return json.Unmarshal(respBody, result)
	}
	return nil
}

// Delete performs a DELETE request.
func (c *HTTPClient) Delete(ctx context.Context, path string, result interface{}) error {
	respBody, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	if result != nil && len(respBody) > 0 {
		return json.Unmarshal(respBody, result)
	}
	return nil
}

// DeleteWithBody performs a DELETE request with a body.
func (c *HTTPClient) DeleteWithBody(ctx context.Context, path string, body, result interface{}) error {
	respBody, err := c.doRequest(ctx, http.MethodDelete, path, body)
	if err != nil {
		return err
	}
	if result != nil && len(respBody) > 0 {
		return json.Unmarshal(respBody, result)
	}
	return nil
}

// BaseURL returns the base URL of the client.
func (c *HTTPClient) BaseURL() string {
	return c.baseURL
}
