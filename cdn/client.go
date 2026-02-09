package cdn

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/stack0/sdk-go/client"
)

// AllowedWidths are the widths that match CloudFront url-rewriter configuration.
var AllowedWidths = []int{256, 384, 640, 750, 828, 1080, 1200, 1920, 2048, 3840}

// Client handles CDN operations.
type Client struct {
	http   *client.HTTPClient
	cdnURL string
}

// NewClient creates a new CDN client.
func NewClient(http *client.HTTPClient, cdnURL string) *Client {
	return &Client{http: http, cdnURL: cdnURL}
}

// GetUploadURL generates a presigned URL for uploading a file.
func (c *Client) GetUploadURL(ctx context.Context, req *UploadURLRequest) (*UploadURLResponse, error) {
	var resp UploadURLResponse
	if err := c.http.Post(ctx, "/cdn/upload", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ConfirmUpload confirms that an upload has completed.
func (c *Client) ConfirmUpload(ctx context.Context, assetID string) (*Asset, error) {
	var resp Asset
	if err := c.http.Post(ctx, "/cdn/upload/"+assetID+"/confirm", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves an asset by ID.
func (c *Client) Get(ctx context.Context, id string) (*Asset, error) {
	var resp Asset
	if err := c.http.Get(ctx, "/cdn/assets/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates asset metadata.
func (c *Client) Update(ctx context.Context, req *UpdateAssetRequest) (*Asset, error) {
	var resp Asset
	if err := c.http.Patch(ctx, "/cdn/assets/"+req.ID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes an asset.
func (c *Client) Delete(ctx context.Context, id string) (*SuccessResponse, error) {
	var resp SuccessResponse
	if err := c.http.DeleteWithBody(ctx, "/cdn/assets/"+id, map[string]string{"id": id}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteMany deletes multiple assets.
func (c *Client) DeleteMany(ctx context.Context, ids []string) (*DeleteAssetsResponse, error) {
	var resp DeleteAssetsResponse
	if err := c.http.Post(ctx, "/cdn/assets/delete", map[string][]string{"ids": ids}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// List lists assets with filters and pagination.
func (c *Client) List(ctx context.Context, req *ListAssetsRequest) (*ListAssetsResponse, error) {
	params := url.Values{}
	params.Set("projectSlug", req.ProjectSlug)
	if req.Folder != nil {
		params.Set("folder", *req.Folder)
	}
	if req.Type != nil {
		params.Set("type", string(*req.Type))
	}
	if req.Status != nil {
		params.Set("status", string(*req.Status))
	}
	if req.Search != nil {
		params.Set("search", *req.Search)
	}
	if len(req.Tags) > 0 {
		params.Set("tags", strings.Join(req.Tags, ","))
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

	path := "/cdn/assets?" + params.Encode()
	var resp ListAssetsResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Move moves assets to a different folder.
func (c *Client) Move(ctx context.Context, req *MoveAssetsRequest) (*MoveAssetsResponse, error) {
	var resp MoveAssetsResponse
	if err := c.http.Post(ctx, "/cdn/assets/move", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetTransformURL generates a transformed image URL client-side.
func (c *Client) GetTransformURL(assetURLOrS3Key string, options *TransformOptions) (string, error) {
	var baseURL string
	if strings.HasPrefix(assetURLOrS3Key, "http://") || strings.HasPrefix(assetURLOrS3Key, "https://") {
		parsed, err := url.Parse(assetURLOrS3Key)
		if err != nil {
			return "", err
		}
		baseURL = fmt.Sprintf("%s://%s%s", parsed.Scheme, parsed.Host, parsed.Path)
	} else if c.cdnURL != "" {
		cdnBase := strings.TrimSuffix(c.cdnURL, "/")
		baseURL = cdnBase + "/" + assetURLOrS3Key
	} else {
		return "", fmt.Errorf("getTransformURL requires either a full URL or cdnURL to be configured")
	}

	params := c.buildTransformQuery(options)
	if params == "" {
		return baseURL, nil
	}
	return baseURL + "?" + params, nil
}

func (c *Client) buildTransformQuery(options *TransformOptions) string {
	if options == nil {
		return ""
	}
	params := url.Values{}
	if options.Format != nil {
		params.Set("f", *options.Format)
	}
	if options.Quality != nil {
		params.Set("q", strconv.Itoa(*options.Quality))
	}
	if options.Width != nil {
		width := c.getNearestWidth(*options.Width)
		params.Set("w", strconv.Itoa(width))
	}
	if options.Height != nil {
		params.Set("h", strconv.Itoa(*options.Height))
	}
	if options.Fit != nil {
		params.Set("fit", *options.Fit)
	}
	if options.Crop != nil {
		params.Set("crop", *options.Crop)
	}
	if options.CropX != nil {
		params.Set("crop-x", strconv.Itoa(*options.CropX))
	}
	if options.CropY != nil {
		params.Set("crop-y", strconv.Itoa(*options.CropY))
	}
	if options.CropWidth != nil {
		params.Set("crop-w", strconv.Itoa(*options.CropWidth))
	}
	if options.CropHeight != nil {
		params.Set("crop-h", strconv.Itoa(*options.CropHeight))
	}
	if options.Blur != nil {
		params.Set("blur", strconv.Itoa(*options.Blur))
	}
	if options.Sharpen != nil {
		params.Set("sharpen", strconv.Itoa(*options.Sharpen))
	}
	if options.Brightness != nil {
		params.Set("brightness", strconv.Itoa(*options.Brightness))
	}
	if options.Saturation != nil {
		params.Set("saturation", strconv.Itoa(*options.Saturation))
	}
	if options.Grayscale {
		params.Set("grayscale", "true")
	}
	if options.Rotate != nil {
		params.Set("rotate", strconv.Itoa(*options.Rotate))
	}
	if options.Flip {
		params.Set("flip", "y")
	}
	if options.Flop {
		params.Set("flop", "x")
	}
	return params.Encode()
}

func (c *Client) getNearestWidth(width int) int {
	nearest := AllowedWidths[0]
	for _, w := range AllowedWidths {
		if abs(w-width) < abs(nearest-width) {
			nearest = w
		}
	}
	return nearest
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GetFolderTree gets the folder tree for navigation.
func (c *Client) GetFolderTree(ctx context.Context, req *GetFolderTreeRequest) ([]FolderTreeNode, error) {
	params := url.Values{}
	params.Set("projectSlug", req.ProjectSlug)
	if req.MaxDepth != nil {
		params.Set("maxDepth", strconv.Itoa(*req.MaxDepth))
	}

	var resp struct {
		Tree []FolderTreeNode `json:"tree"`
	}
	if err := c.http.Get(ctx, "/cdn/folders/tree?"+params.Encode(), &resp); err != nil {
		return nil, err
	}
	return resp.Tree, nil
}

// CreateFolder creates a new folder.
func (c *Client) CreateFolder(ctx context.Context, req *CreateFolderRequest) (*Folder, error) {
	var resp Folder
	if err := c.http.Post(ctx, "/cdn/folders", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetFolder retrieves a folder by ID.
func (c *Client) GetFolder(ctx context.Context, id string) (*Folder, error) {
	var resp Folder
	if err := c.http.Get(ctx, "/cdn/folders/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetFolderByPath retrieves a folder by its path.
func (c *Client) GetFolderByPath(ctx context.Context, path string) (*Folder, error) {
	var resp Folder
	if err := c.http.Get(ctx, "/cdn/folders/path/"+url.PathEscape(path), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateFolder updates a folder's name.
func (c *Client) UpdateFolder(ctx context.Context, req *UpdateFolderRequest) (*Folder, error) {
	var resp Folder
	if err := c.http.Patch(ctx, "/cdn/folders/"+req.ID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListFolders lists folders with optional filters.
func (c *Client) ListFolders(ctx context.Context, req *ListFoldersRequest) (*ListFoldersResponse, error) {
	params := url.Values{}
	if req != nil {
		if req.ParentID != nil {
			params.Set("parentId", *req.ParentID)
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

	path := "/cdn/folders"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListFoldersResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MoveFolder moves a folder to a new parent.
func (c *Client) MoveFolder(ctx context.Context, req *MoveFolderRequest) (*MoveFolderResponse, error) {
	var resp MoveFolderResponse
	if err := c.http.Post(ctx, "/cdn/folders/move", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteFolder deletes a folder.
func (c *Client) DeleteFolder(ctx context.Context, id string, deleteContents bool) (*SuccessResponse, error) {
	params := url.Values{}
	if deleteContents {
		params.Set("deleteContents", "true")
	}
	path := "/cdn/folders/" + id
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	var resp SuccessResponse
	if err := c.http.DeleteWithBody(ctx, path, map[string]interface{}{"id": id, "deleteContents": deleteContents}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
