package mail

import (
	"context"
	"net/url"
	"strconv"

	"github.com/stack0/sdk-go/client"
)

// TemplatesClient handles template operations.
type TemplatesClient struct {
	http *client.HTTPClient
}

// NewTemplatesClient creates a new templates client.
func NewTemplatesClient(http *client.HTTPClient) *TemplatesClient {
	return &TemplatesClient{http: http}
}

// List lists all templates.
func (c *TemplatesClient) List(ctx context.Context, req *ListTemplatesRequest) (*ListTemplatesResponse, error) {
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
		if req.IsActive != nil {
			params.Set("isActive", strconv.FormatBool(*req.IsActive))
		}
		if req.Search != nil {
			params.Set("search", *req.Search)
		}
	}

	path := "/mail/templates"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListTemplatesResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a template by ID.
func (c *TemplatesClient) Get(ctx context.Context, id string) (*Template, error) {
	var resp Template
	if err := c.http.Get(ctx, "/mail/templates/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBySlug retrieves a template by slug.
func (c *TemplatesClient) GetBySlug(ctx context.Context, slug string) (*Template, error) {
	var resp Template
	if err := c.http.Get(ctx, "/mail/templates/slug/"+slug, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new template.
func (c *TemplatesClient) Create(ctx context.Context, req *CreateTemplateRequest) (*Template, error) {
	var resp Template
	if err := c.http.Post(ctx, "/mail/templates", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates a template.
func (c *TemplatesClient) Update(ctx context.Context, req *UpdateTemplateRequest) (*Template, error) {
	var resp Template
	if err := c.http.Put(ctx, "/mail/templates/"+req.ID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes a template.
func (c *TemplatesClient) Delete(ctx context.Context, id string) (*DeleteTemplateResponse, error) {
	var resp DeleteTemplateResponse
	if err := c.http.Delete(ctx, "/mail/templates/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Preview previews a template with variables.
func (c *TemplatesClient) Preview(ctx context.Context, req *PreviewTemplateRequest) (*PreviewTemplateResponse, error) {
	var resp PreviewTemplateResponse
	body := map[string]interface{}{"variables": req.Variables}
	if err := c.http.Post(ctx, "/mail/templates/"+req.ID+"/preview", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
