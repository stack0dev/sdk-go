package cdn

import (
	"context"
	"net/url"
	"strconv"
)

// CreateImport creates an S3 import job to bulk import files.
func (c *Client) CreateImport(ctx context.Context, req *CreateImportRequest) (*CreateImportResponse, error) {
	var resp CreateImportResponse
	if err := c.http.Post(ctx, "/cdn/imports", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetImport retrieves an import job by ID.
func (c *Client) GetImport(ctx context.Context, importID string) (*ImportJob, error) {
	var resp ImportJob
	if err := c.http.Get(ctx, "/cdn/imports/"+importID, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListImports lists import jobs with pagination and filters.
func (c *Client) ListImports(ctx context.Context, req *ListImportsRequest) (*ListImportsResponse, error) {
	params := url.Values{}
	params.Set("projectSlug", req.ProjectSlug)
	if req.Environment != nil {
		params.Set("environment", string(*req.Environment))
	}
	if req.Status != nil {
		params.Set("status", string(*req.Status))
	}
	if req.SortBy != nil {
		params.Set("sortBy", *req.SortBy)
	}
	if req.SortOrder != nil {
		params.Set("sortOrder", *req.SortOrder)
	}
	if req.Limit != nil {
		params.Set("limit", strconv.Itoa(*req.Limit))
	}
	if req.Offset != nil {
		params.Set("offset", strconv.Itoa(*req.Offset))
	}

	var resp ListImportsResponse
	if err := c.http.Get(ctx, "/cdn/imports?"+params.Encode(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelImport cancels a running import job.
func (c *Client) CancelImport(ctx context.Context, importID string) (*CancelImportResponse, error) {
	var resp CancelImportResponse
	if err := c.http.Post(ctx, "/cdn/imports/"+importID+"/cancel", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RetryImport retries failed files in an import job.
func (c *Client) RetryImport(ctx context.Context, importID string) (*RetryImportResponse, error) {
	var resp RetryImportResponse
	if err := c.http.Post(ctx, "/cdn/imports/"+importID+"/retry", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListImportFiles lists files in an import job.
func (c *Client) ListImportFiles(ctx context.Context, req *ListImportFilesRequest) (*ListImportFilesResponse, error) {
	params := url.Values{}
	if req.Status != nil {
		params.Set("status", string(*req.Status))
	}
	if req.SortBy != nil {
		params.Set("sortBy", *req.SortBy)
	}
	if req.SortOrder != nil {
		params.Set("sortOrder", *req.SortOrder)
	}
	if req.Limit != nil {
		params.Set("limit", strconv.Itoa(*req.Limit))
	}
	if req.Offset != nil {
		params.Set("offset", strconv.Itoa(*req.Offset))
	}

	path := "/cdn/imports/" + req.ImportID + "/files"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListImportFilesResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
