// Package stack0 provides the official Go SDK for Stack0 services.
//
// Example usage:
//
//	client := stack0.New("stack0_...")
//
//	// Send an email
//	resp, err := client.Mail.Send(ctx, &mail.SendEmailRequest{
//	    From: &mail.EmailAddress{Email: "noreply@example.com"},
//	    To:   []mail.EmailAddress{{Email: "user@example.com"}},
//	    Subject: "Hello",
//	    HTML: ptr("&lt;p&gt;World&lt;/p&gt;"),
//	})
//
//	// Capture a screenshot
//	screenshot, err := client.Screenshots.CaptureAndWait(ctx, &screenshots.CreateScreenshotRequest{
//	    URL: "https://example.com",
//	})
package stack0

import (
	"github.com/stack0/sdk-go/cdn"
	"github.com/stack0/sdk-go/client"
	"github.com/stack0/sdk-go/extraction"
	"github.com/stack0/sdk-go/mail"
	"github.com/stack0/sdk-go/screenshots"
)

const (
	// DefaultBaseURL is the default Stack0 API endpoint.
	DefaultBaseURL = "https://api.stack0.dev"
)

// Client is the main Stack0 SDK client.
type Client struct {
	// Mail provides access to email sending and management.
	Mail *mail.Client

	// CDN provides access to asset upload and management.
	CDN *cdn.Client

	// Screenshots provides access to webpage screenshot capture.
	Screenshots *screenshots.Client

	// Extraction provides access to AI content extraction.
	Extraction *extraction.Client
}

// Option is a functional option for configuring the Client.
type Option func(*options)

type options struct {
	baseURL string
	cdnURL  string
}

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(url string) Option {
	return func(o *options) {
		o.baseURL = url
	}
}

// WithCDNURL sets the CDN URL for image transformations.
// When provided, transform URLs are generated client-side without API calls.
func WithCDNURL(url string) Option {
	return func(o *options) {
		o.cdnURL = url
	}
}

// New creates a new Stack0 client with the given API key.
func New(apiKey string, opts ...Option) *Client {
	o := &options{
		baseURL: DefaultBaseURL,
	}
	for _, opt := range opts {
		opt(o)
	}

	httpClient := client.NewHTTPClient(apiKey, o.baseURL)

	return &Client{
		Mail:        mail.NewClient(httpClient),
		CDN:         cdn.NewClient(httpClient, o.cdnURL),
		Screenshots: screenshots.NewClient(httpClient),
		Extraction:  extraction.NewClient(httpClient),
	}
}
