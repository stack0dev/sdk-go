package mail

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/stack0/sdk-go/client"
)

// Client is the mail client for the Stack0 SDK.
type Client struct {
	http      *client.HTTPClient
	Domains   *DomainsClient
	Templates *TemplatesClient
	Audiences *AudiencesClient
	Contacts  *ContactsClient
	Campaigns *CampaignsClient
	Sequences *SequencesClient
	Events    *EventsClient
}

// New creates a new mail client.
func New(http *client.HTTPClient) *Client {
	return &Client{
		http:      http,
		Domains:   NewDomainsClient(http),
		Templates: NewTemplatesClient(http),
		Audiences: NewAudiencesClient(http),
		Contacts:  NewContactsClient(http),
		Campaigns: NewCampaignsClient(http),
		Sequences: NewSequencesClient(http),
		Events:    NewEventsClient(http),
	}
}

// Send sends a single email.
func (c *Client) Send(ctx context.Context, req *SendEmailRequest) (*SendEmailResponse, error) {
	var resp SendEmailResponse
	if err := c.http.Post(ctx, "/mail/send", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendBatch sends multiple emails in a batch.
func (c *Client) SendBatch(ctx context.Context, req *SendBatchEmailRequest) (*SendBatchEmailResponse, error) {
	var resp SendBatchEmailResponse
	if err := c.http.Post(ctx, "/mail/send/batch", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendBroadcast sends a broadcast email to multiple recipients.
func (c *Client) SendBroadcast(ctx context.Context, req *SendBroadcastEmailRequest) (*SendBroadcastEmailResponse, error) {
	var resp SendBroadcastEmailResponse
	if err := c.http.Post(ctx, "/mail/send/broadcast", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves an email by ID.
func (c *Client) Get(ctx context.Context, id string) (*GetEmailResponse, error) {
	var resp GetEmailResponse
	if err := c.http.Get(ctx, "/mail/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// List lists emails with optional filters.
func (c *Client) List(ctx context.Context, req *ListEmailsRequest) (*ListEmailsResponse, error) {
	params := url.Values{}
	if req != nil {
		if req.ProjectSlug != nil {
			params.Set("projectSlug", *req.ProjectSlug)
		}
		if req.Environment != nil {
			params.Set("environment", string(*req.Environment))
		}
		if req.Limit != nil {
			params.Set("limit", strconv.Itoa(*req.Limit))
		}
		if req.Offset != nil {
			params.Set("offset", strconv.Itoa(*req.Offset))
		}
		if req.Status != nil {
			params.Set("status", string(*req.Status))
		}
		if req.From != nil {
			params.Set("from", *req.From)
		}
		if req.To != nil {
			params.Set("to", *req.To)
		}
		if req.Subject != nil {
			params.Set("subject", *req.Subject)
		}
		if req.Tag != nil {
			params.Set("tag", *req.Tag)
		}
		if req.StartDate != nil {
			params.Set("startDate", req.StartDate.Format("2006-01-02T15:04:05Z07:00"))
		}
		if req.EndDate != nil {
			params.Set("endDate", req.EndDate.Format("2006-01-02T15:04:05Z07:00"))
		}
		if req.SortBy != nil {
			params.Set("sortBy", *req.SortBy)
		}
		if req.SortOrder != nil {
			params.Set("sortOrder", *req.SortOrder)
		}
	}

	path := "/mail"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListEmailsResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Resend resends an email by ID.
func (c *Client) Resend(ctx context.Context, id string) (*ResendEmailResponse, error) {
	var resp ResendEmailResponse
	if err := c.http.Post(ctx, "/mail/"+id+"/resend", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Cancel cancels a scheduled email by ID.
func (c *Client) Cancel(ctx context.Context, id string) (*CancelEmailResponse, error) {
	var resp CancelEmailResponse
	if err := c.http.Post(ctx, "/mail/"+id+"/cancel", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAnalytics retrieves overall email analytics.
func (c *Client) GetAnalytics(ctx context.Context) (*EmailAnalyticsResponse, error) {
	var resp EmailAnalyticsResponse
	if err := c.http.Get(ctx, "/mail/analytics", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetTimeSeriesAnalytics retrieves time series analytics.
func (c *Client) GetTimeSeriesAnalytics(ctx context.Context, days *int) (*TimeSeriesAnalyticsResponse, error) {
	path := "/mail/analytics/timeseries"
	if days != nil {
		path += fmt.Sprintf("?days=%d", *days)
	}

	var resp TimeSeriesAnalyticsResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetHourlyAnalytics retrieves hourly analytics.
func (c *Client) GetHourlyAnalytics(ctx context.Context) (*HourlyAnalyticsResponse, error) {
	var resp HourlyAnalyticsResponse
	if err := c.http.Get(ctx, "/mail/analytics/hourly", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListSenders lists unique senders with statistics.
func (c *Client) ListSenders(ctx context.Context, req *ListSendersRequest) (*ListSendersResponse, error) {
	params := url.Values{}
	if req != nil {
		if req.ProjectSlug != nil {
			params.Set("projectSlug", *req.ProjectSlug)
		}
		if req.Environment != nil {
			params.Set("environment", string(*req.Environment))
		}
		if req.Search != nil {
			params.Set("search", *req.Search)
		}
	}

	path := "/mail/senders"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListSendersResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
