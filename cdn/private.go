package cdn

import (
	"context"
	"net/url"
	"strconv"
)

// GetPrivateUploadURL generates a presigned URL for uploading a private file.
func (c *Client) GetPrivateUploadURL(ctx context.Context, req *PrivateUploadURLRequest) (*PrivateUploadURLResponse, error) {
	var resp PrivateUploadURLResponse
	if err := c.http.Post(ctx, "/cdn/private/upload", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ConfirmPrivateUpload confirms that a private file upload has completed.
func (c *Client) ConfirmPrivateUpload(ctx context.Context, fileID string) (*PrivateFile, error) {
	var resp PrivateFile
	if err := c.http.Post(ctx, "/cdn/private/upload/"+fileID+"/confirm", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPrivateDownloadURL generates a presigned download URL for a private file.
func (c *Client) GetPrivateDownloadURL(ctx context.Context, req *PrivateDownloadURLRequest) (*PrivateDownloadURLResponse, error) {
	body := map[string]interface{}{}
	if req.ExpiresIn != nil {
		body["expiresIn"] = *req.ExpiresIn
	}
	var resp PrivateDownloadURLResponse
	if err := c.http.Post(ctx, "/cdn/private/"+req.FileID+"/download", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPrivateFile retrieves a private file by ID.
func (c *Client) GetPrivateFile(ctx context.Context, fileID string) (*PrivateFile, error) {
	var resp PrivateFile
	if err := c.http.Get(ctx, "/cdn/private/"+fileID, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdatePrivateFile updates a private file's metadata.
func (c *Client) UpdatePrivateFile(ctx context.Context, req *UpdatePrivateFileRequest) (*PrivateFile, error) {
	var resp PrivateFile
	if err := c.http.Patch(ctx, "/cdn/private/"+req.FileID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeletePrivateFile deletes a private file.
func (c *Client) DeletePrivateFile(ctx context.Context, fileID string) (*SuccessResponse, error) {
	var resp SuccessResponse
	if err := c.http.DeleteWithBody(ctx, "/cdn/private/"+fileID, map[string]string{"fileId": fileID}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeletePrivateFiles deletes multiple private files.
func (c *Client) DeletePrivateFiles(ctx context.Context, fileIDs []string) (*DeleteAssetsResponse, error) {
	var resp DeleteAssetsResponse
	if err := c.http.Post(ctx, "/cdn/private/delete", map[string][]string{"fileIds": fileIDs}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListPrivateFiles lists private files with filters and pagination.
func (c *Client) ListPrivateFiles(ctx context.Context, req *ListPrivateFilesRequest) (*ListPrivateFilesResponse, error) {
	params := url.Values{}
	params.Set("projectSlug", req.ProjectSlug)
	if req.Folder != nil {
		params.Set("folder", *req.Folder)
	}
	if req.Status != nil {
		params.Set("status", string(*req.Status))
	}
	if req.Search != nil {
		params.Set("search", *req.Search)
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

	var resp ListPrivateFilesResponse
	if err := c.http.Get(ctx, "/cdn/private?"+params.Encode(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MovePrivateFiles moves private files to a different folder.
func (c *Client) MovePrivateFiles(ctx context.Context, req *MovePrivateFilesRequest) (*MovePrivateFilesResponse, error) {
	var resp MovePrivateFilesResponse
	if err := c.http.Post(ctx, "/cdn/private/move", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateBundle creates a download bundle from assets and/or private files.
func (c *Client) CreateBundle(ctx context.Context, req *CreateBundleRequest) (*CreateBundleResponse, error) {
	var resp CreateBundleResponse
	if err := c.http.Post(ctx, "/cdn/bundles", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBundle retrieves a download bundle by ID.
func (c *Client) GetBundle(ctx context.Context, bundleID string) (*DownloadBundle, error) {
	var resp DownloadBundle
	if err := c.http.Get(ctx, "/cdn/bundles/"+bundleID, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListBundles lists download bundles with filters and pagination.
func (c *Client) ListBundles(ctx context.Context, req *ListBundlesRequest) (*ListBundlesResponse, error) {
	params := url.Values{}
	params.Set("projectSlug", req.ProjectSlug)
	if req.Status != nil {
		params.Set("status", string(*req.Status))
	}
	if req.Search != nil {
		params.Set("search", *req.Search)
	}
	if req.Limit != nil {
		params.Set("limit", strconv.Itoa(*req.Limit))
	}
	if req.Offset != nil {
		params.Set("offset", strconv.Itoa(*req.Offset))
	}

	var resp ListBundlesResponse
	if err := c.http.Get(ctx, "/cdn/bundles?"+params.Encode(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBundleDownloadURL generates a presigned download URL for a bundle.
func (c *Client) GetBundleDownloadURL(ctx context.Context, req *BundleDownloadURLRequest) (*BundleDownloadURLResponse, error) {
	body := map[string]interface{}{}
	if req.ExpiresIn != nil {
		body["expiresIn"] = *req.ExpiresIn
	}
	var resp BundleDownloadURLResponse
	if err := c.http.Post(ctx, "/cdn/bundles/"+req.BundleID+"/download", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteBundle deletes a download bundle.
func (c *Client) DeleteBundle(ctx context.Context, bundleID string) (*SuccessResponse, error) {
	var resp SuccessResponse
	if err := c.http.DeleteWithBody(ctx, "/cdn/bundles/"+bundleID, map[string]string{"bundleId": bundleID}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
