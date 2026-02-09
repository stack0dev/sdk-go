package mail

import (
	"context"
	"net/url"
	"strconv"

	"github.com/stack0/sdk-go/client"
)

// AudiencesClient handles audience operations.
type AudiencesClient struct {
	http *client.HTTPClient
}

// NewAudiencesClient creates a new audiences client.
func NewAudiencesClient(http *client.HTTPClient) *AudiencesClient {
	return &AudiencesClient{http: http}
}

// List lists all audiences.
func (c *AudiencesClient) List(ctx context.Context, req *ListAudiencesRequest) (*ListAudiencesResponse, error) {
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
	}

	path := "/mail/audiences"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListAudiencesResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves an audience by ID.
func (c *AudiencesClient) Get(ctx context.Context, id string) (*Audience, error) {
	var resp Audience
	if err := c.http.Get(ctx, "/mail/audiences/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new audience.
func (c *AudiencesClient) Create(ctx context.Context, req *CreateAudienceRequest) (*Audience, error) {
	var resp Audience
	if err := c.http.Post(ctx, "/mail/audiences", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates an audience.
func (c *AudiencesClient) Update(ctx context.Context, req *UpdateAudienceRequest) (*Audience, error) {
	var resp Audience
	if err := c.http.Put(ctx, "/mail/audiences/"+req.ID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes an audience.
func (c *AudiencesClient) Delete(ctx context.Context, id string) (*DeleteAudienceResponse, error) {
	var resp DeleteAudienceResponse
	if err := c.http.Delete(ctx, "/mail/audiences/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListContacts lists contacts in an audience.
func (c *AudiencesClient) ListContacts(ctx context.Context, req *ListAudienceContactsRequest) (*ListAudienceContactsResponse, error) {
	params := url.Values{}
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

	path := "/mail/audiences/" + req.ID + "/contacts"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListAudienceContactsResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AddContacts adds contacts to an audience.
func (c *AudiencesClient) AddContacts(ctx context.Context, req *AddContactsToAudienceRequest) (*AddContactsToAudienceResponse, error) {
	var resp AddContactsToAudienceResponse
	body := map[string]interface{}{"contactIds": req.ContactIDs}
	if err := c.http.Post(ctx, "/mail/audiences/"+req.ID+"/contacts", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RemoveContacts removes contacts from an audience.
func (c *AudiencesClient) RemoveContacts(ctx context.Context, req *RemoveContactsFromAudienceRequest) (*RemoveContactsFromAudienceResponse, error) {
	var resp RemoveContactsFromAudienceResponse
	body := map[string]interface{}{"contactIds": req.ContactIDs}
	if err := c.http.DeleteWithBody(ctx, "/mail/audiences/"+req.ID+"/contacts", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
