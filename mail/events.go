package mail

import (
	"context"
	"net/url"
	"strconv"

	"github.com/stack0/sdk-go/client"
)

// EventsClient handles event operations.
type EventsClient struct {
	http *client.HTTPClient
}

// NewEventsClient creates a new events client.
func NewEventsClient(http *client.HTTPClient) *EventsClient {
	return &EventsClient{http: http}
}

// List lists all event definitions.
func (c *EventsClient) List(ctx context.Context, req *ListEventsRequest) (*ListEventsResponse, error) {
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
		if req.Search != nil {
			params.Set("search", *req.Search)
		}
	}

	path := "/mail/events"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListEventsResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves an event definition by ID.
func (c *EventsClient) Get(ctx context.Context, id string) (*MailEvent, error) {
	var resp MailEvent
	if err := c.http.Get(ctx, "/mail/events/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new event definition.
func (c *EventsClient) Create(ctx context.Context, req *CreateEventRequest) (*MailEvent, error) {
	var resp MailEvent
	if err := c.http.Post(ctx, "/mail/events", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates an event definition.
func (c *EventsClient) Update(ctx context.Context, req *UpdateEventRequest) (*MailEvent, error) {
	var resp MailEvent
	if err := c.http.Put(ctx, "/mail/events/"+req.ID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes an event definition.
func (c *EventsClient) Delete(ctx context.Context, id string) (*DeleteEventResponse, error) {
	var resp DeleteEventResponse
	if err := c.http.Delete(ctx, "/mail/events/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Track tracks a single event.
func (c *EventsClient) Track(ctx context.Context, req *TrackEventRequest) (*TrackEventResponse, error) {
	var resp TrackEventResponse
	if err := c.http.Post(ctx, "/mail/events/track", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TrackBatch tracks multiple events in a batch.
func (c *EventsClient) TrackBatch(ctx context.Context, req *BatchTrackEventsRequest) (*BatchTrackEventsResponse, error) {
	var resp BatchTrackEventsResponse
	if err := c.http.Post(ctx, "/mail/events/track/batch", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListOccurrences lists event occurrences.
func (c *EventsClient) ListOccurrences(ctx context.Context, req *ListEventOccurrencesRequest) (*ListEventOccurrencesResponse, error) {
	params := url.Values{}
	if req != nil {
		if req.EventID != nil {
			params.Set("eventId", *req.EventID)
		}
		if req.ContactID != nil {
			params.Set("contactId", *req.ContactID)
		}
		if req.Limit != nil {
			params.Set("limit", strconv.Itoa(*req.Limit))
		}
		if req.Offset != nil {
			params.Set("offset", strconv.Itoa(*req.Offset))
		}
		if req.StartDate != nil {
			params.Set("startDate", req.StartDate.Format("2006-01-02T15:04:05Z07:00"))
		}
		if req.EndDate != nil {
			params.Set("endDate", req.EndDate.Format("2006-01-02T15:04:05Z07:00"))
		}
	}

	path := "/mail/events/occurrences"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListEventOccurrencesResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAnalytics retrieves analytics for an event.
func (c *EventsClient) GetAnalytics(ctx context.Context, id string) (*EventAnalyticsResponse, error) {
	var resp EventAnalyticsResponse
	if err := c.http.Get(ctx, "/mail/events/analytics/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
