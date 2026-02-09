package mail

import (
	"context"
	"net/url"

	"github.com/stack0/sdk-go/client"
)

// DomainsClient handles domain operations.
type DomainsClient struct {
	http *client.HTTPClient
}

// NewDomainsClient creates a new domains client.
func NewDomainsClient(http *client.HTTPClient) *DomainsClient {
	return &DomainsClient{http: http}
}

// List lists all domains.
func (c *DomainsClient) List(ctx context.Context, req *ListDomainsRequest) ([]Domain, error) {
	params := url.Values{}
	params.Set("projectSlug", req.ProjectSlug)
	if req.Environment != nil {
		params.Set("environment", string(*req.Environment))
	}

	var resp []Domain
	if err := c.http.Get(ctx, "/mail/domains?"+params.Encode(), &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Add adds a new domain.
func (c *DomainsClient) Add(ctx context.Context, req *AddDomainRequest) (*AddDomainResponse, error) {
	var resp AddDomainResponse
	if err := c.http.Post(ctx, "/mail/domains", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDNSRecords retrieves DNS records for a domain.
func (c *DomainsClient) GetDNSRecords(ctx context.Context, domainID string) (*GetDNSRecordsResponse, error) {
	var resp GetDNSRecordsResponse
	if err := c.http.Get(ctx, "/mail/domains/"+domainID+"/dns", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Verify verifies a domain.
func (c *DomainsClient) Verify(ctx context.Context, domainID string) (*VerifyDomainResponse, error) {
	var resp VerifyDomainResponse
	if err := c.http.Post(ctx, "/mail/domains/"+domainID+"/verify", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes a domain.
func (c *DomainsClient) Delete(ctx context.Context, domainID string) (*DeleteDomainResponse, error) {
	var resp DeleteDomainResponse
	if err := c.http.Delete(ctx, "/mail/domains/"+domainID, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetDefault sets a domain as the default.
func (c *DomainsClient) SetDefault(ctx context.Context, domainID string) (*Domain, error) {
	var resp Domain
	if err := c.http.Post(ctx, "/mail/domains/"+domainID+"/default", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
