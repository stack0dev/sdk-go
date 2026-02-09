package mail

import (
	"context"
	"net/url"
	"strconv"

	"github.com/stack0/sdk-go/client"
)

// ContactsClient handles contact operations.
type ContactsClient struct {
	http *client.HTTPClient
}

// NewContactsClient creates a new contacts client.
func NewContactsClient(http *client.HTTPClient) *ContactsClient {
	return &ContactsClient{http: http}
}

// List lists all contacts.
func (c *ContactsClient) List(ctx context.Context, req *ListContactsRequest) (*ListContactsResponse, error) {
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

	path := "/mail/contacts"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListContactsResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a contact by ID.
func (c *ContactsClient) Get(ctx context.Context, id string) (*MailContact, error) {
	var resp MailContact
	if err := c.http.Get(ctx, "/mail/contacts/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new contact.
func (c *ContactsClient) Create(ctx context.Context, req *CreateContactRequest) (*MailContact, error) {
	var resp MailContact
	if err := c.http.Post(ctx, "/mail/contacts", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates a contact.
func (c *ContactsClient) Update(ctx context.Context, req *UpdateContactRequest) (*MailContact, error) {
	var resp MailContact
	if err := c.http.Put(ctx, "/mail/contacts/"+req.ID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes a contact.
func (c *ContactsClient) Delete(ctx context.Context, id string) (*DeleteContactResponse, error) {
	var resp DeleteContactResponse
	if err := c.http.Delete(ctx, "/mail/contacts/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Import imports contacts in bulk.
func (c *ContactsClient) Import(ctx context.Context, req *ImportContactsRequest) (*ImportContactsResponse, error) {
	var resp ImportContactsResponse
	if err := c.http.Post(ctx, "/mail/contacts/import", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
