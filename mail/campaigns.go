package mail

import (
	"context"
	"net/url"
	"strconv"

	"github.com/stack0/sdk-go/client"
)

// CampaignsClient handles campaign operations.
type CampaignsClient struct {
	http *client.HTTPClient
}

// NewCampaignsClient creates a new campaigns client.
func NewCampaignsClient(http *client.HTTPClient) *CampaignsClient {
	return &CampaignsClient{http: http}
}

// List lists all campaigns.
func (c *CampaignsClient) List(ctx context.Context, req *ListCampaignsRequest) (*ListCampaignsResponse, error) {
	params := url.Values{}
	if req != nil {
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
		if req.Status != nil {
			params.Set("status", string(*req.Status))
		}
	}

	path := "/mail/campaigns"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListCampaignsResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a campaign by ID.
func (c *CampaignsClient) Get(ctx context.Context, id string) (*Campaign, error) {
	var resp Campaign
	if err := c.http.Get(ctx, "/mail/campaigns/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new campaign.
func (c *CampaignsClient) Create(ctx context.Context, req *CreateCampaignRequest) (*Campaign, error) {
	var resp Campaign
	if err := c.http.Post(ctx, "/mail/campaigns", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates a campaign.
func (c *CampaignsClient) Update(ctx context.Context, req *UpdateCampaignRequest) (*Campaign, error) {
	var resp Campaign
	if err := c.http.Put(ctx, "/mail/campaigns/"+req.ID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes a campaign.
func (c *CampaignsClient) Delete(ctx context.Context, id string) (*DeleteCampaignResponse, error) {
	var resp DeleteCampaignResponse
	if err := c.http.Delete(ctx, "/mail/campaigns/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Send sends a campaign.
func (c *CampaignsClient) Send(ctx context.Context, req *SendCampaignRequest) (*SendCampaignResponse, error) {
	var resp SendCampaignResponse
	body := map[string]interface{}{}
	if req.SendNow != nil {
		body["sendNow"] = *req.SendNow
	}
	if req.ScheduledAt != nil {
		body["scheduledAt"] = req.ScheduledAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if err := c.http.Post(ctx, "/mail/campaigns/"+req.ID+"/send", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Pause pauses a sending campaign.
func (c *CampaignsClient) Pause(ctx context.Context, id string) (*PauseCampaignResponse, error) {
	var resp PauseCampaignResponse
	if err := c.http.Post(ctx, "/mail/campaigns/"+id+"/pause", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Cancel cancels a campaign.
func (c *CampaignsClient) Cancel(ctx context.Context, id string) (*CancelCampaignResponse, error) {
	var resp CancelCampaignResponse
	if err := c.http.Post(ctx, "/mail/campaigns/"+id+"/cancel", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Duplicate duplicates a campaign.
func (c *CampaignsClient) Duplicate(ctx context.Context, id string) (*Campaign, error) {
	var resp Campaign
	if err := c.http.Post(ctx, "/mail/campaigns/"+id+"/duplicate", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetStats retrieves campaign statistics.
func (c *CampaignsClient) GetStats(ctx context.Context, id string) (*CampaignStatsResponse, error) {
	var resp CampaignStatsResponse
	if err := c.http.Get(ctx, "/mail/campaigns/"+id+"/stats", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
