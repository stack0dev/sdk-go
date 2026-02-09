package stack0

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_DefaultConfiguration(t *testing.T) {
	client := New("test-api-key")

	require.NotNil(t, client)
	assert.NotNil(t, client.Mail)
	assert.NotNil(t, client.CDN)
	assert.NotNil(t, client.Screenshots)
	assert.NotNil(t, client.Extraction)
}

func TestNew_WithBaseURL(t *testing.T) {
	customURL := "https://custom.api.example.com"
	client := New("test-api-key", WithBaseURL(customURL))

	require.NotNil(t, client)
	assert.NotNil(t, client.Mail)
	assert.NotNil(t, client.CDN)
	assert.NotNil(t, client.Screenshots)
	assert.NotNil(t, client.Extraction)
}

func TestNew_WithCDNURL(t *testing.T) {
	cdnURL := "https://cdn.example.com"
	client := New("test-api-key", WithCDNURL(cdnURL))

	require.NotNil(t, client)
	assert.NotNil(t, client.CDN)
}

func TestNew_WithMultipleOptions(t *testing.T) {
	customURL := "https://custom.api.example.com"
	cdnURL := "https://cdn.example.com"
	client := New("test-api-key", WithBaseURL(customURL), WithCDNURL(cdnURL))

	require.NotNil(t, client)
	assert.NotNil(t, client.Mail)
	assert.NotNil(t, client.CDN)
	assert.NotNil(t, client.Screenshots)
	assert.NotNil(t, client.Extraction)
}

func TestWithBaseURL(t *testing.T) {
	o := &options{}
	customURL := "https://custom.api.example.com"

	opt := WithBaseURL(customURL)
	opt(o)

	assert.Equal(t, customURL, o.baseURL)
}

func TestWithCDNURL(t *testing.T) {
	o := &options{}
	cdnURL := "https://cdn.example.com"

	opt := WithCDNURL(cdnURL)
	opt(o)

	assert.Equal(t, cdnURL, o.cdnURL)
}

func TestDefaultBaseURL(t *testing.T) {
	assert.Equal(t, "https://api.stack0.dev", DefaultBaseURL)
}
