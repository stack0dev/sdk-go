package cdn

import "time"

// AssetStatus represents the status of an asset.
type AssetStatus string

const (
	AssetStatusPending    AssetStatus = "pending"
	AssetStatusProcessing AssetStatus = "processing"
	AssetStatusReady      AssetStatus = "ready"
	AssetStatusFailed     AssetStatus = "failed"
	AssetStatusDeleted    AssetStatus = "deleted"
)

// AssetType represents the type of an asset.
type AssetType string

const (
	AssetTypeImage    AssetType = "image"
	AssetTypeVideo    AssetType = "video"
	AssetTypeAudio    AssetType = "audio"
	AssetTypeDocument AssetType = "document"
	AssetTypeOther    AssetType = "other"
)

// Asset represents a CDN asset.
type Asset struct {
	ID               string                 `json:"id"`
	Filename         string                 `json:"filename"`
	OriginalFilename string                 `json:"originalFilename"`
	MimeType         string                 `json:"mimeType"`
	Size             int64                  `json:"size"`
	Type             AssetType              `json:"type"`
	S3Key            string                 `json:"s3Key"`
	CDNURL           string                 `json:"cdnUrl"`
	Width            *int                   `json:"width,omitempty"`
	Height           *int                   `json:"height,omitempty"`
	Duration         *float64               `json:"duration,omitempty"`
	Status           AssetStatus            `json:"status"`
	Folder           *string                `json:"folder,omitempty"`
	Tags             []string               `json:"tags,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	Alt              *string                `json:"alt,omitempty"`
	CreatedAt        time.Time              `json:"createdAt"`
	UpdatedAt        *time.Time             `json:"updatedAt,omitempty"`
}

// UploadURLRequest is the request for getting an upload URL.
type UploadURLRequest struct {
	ProjectSlug string                  `json:"projectSlug"`
	Filename    string                  `json:"filename"`
	MimeType    string                  `json:"mimeType"`
	Size        int64                   `json:"size"`
	Folder      *string                 `json:"folder,omitempty"`
	Metadata    map[string]interface{}  `json:"metadata,omitempty"`
	Watermark   *ImageWatermarkConfig   `json:"watermark,omitempty"`
}

// UploadURLResponse is the response from getting an upload URL.
type UploadURLResponse struct {
	UploadURL string    `json:"uploadUrl"`
	AssetID   string    `json:"assetId"`
	CDNURL    string    `json:"cdnUrl"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// UpdateAssetRequest is the request for updating an asset.
type UpdateAssetRequest struct {
	ID       string                 `json:"id"`
	Filename *string                `json:"filename,omitempty"`
	Folder   *string                `json:"folder,omitempty"`
	Tags     []string               `json:"tags,omitempty"`
	Alt      *string                `json:"alt,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// DeleteAssetsResponse is the response from deleting assets.
type DeleteAssetsResponse struct {
	Success      bool `json:"success"`
	DeletedCount int  `json:"deletedCount"`
}

// ListAssetsRequest is the request for listing assets.
type ListAssetsRequest struct {
	ProjectSlug string       `json:"projectSlug"`
	Folder      *string      `json:"folder,omitempty"`
	Type        *AssetType   `json:"type,omitempty"`
	Status      *AssetStatus `json:"status,omitempty"`
	Search      *string      `json:"search,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
	SortBy      *string      `json:"sortBy,omitempty"`
	SortOrder   *string      `json:"sortOrder,omitempty"`
	Limit       *int         `json:"limit,omitempty"`
	Offset      *int         `json:"offset,omitempty"`
}

// ListAssetsResponse is the response from listing assets.
type ListAssetsResponse struct {
	Assets  []Asset `json:"assets"`
	Total   int     `json:"total"`
	HasMore bool    `json:"hasMore"`
}

// MoveAssetsRequest is the request for moving assets.
type MoveAssetsRequest struct {
	AssetIDs []string `json:"assetIds"`
	Folder   *string  `json:"folder"`
}

// MoveAssetsResponse is the response from moving assets.
type MoveAssetsResponse struct {
	Success    bool `json:"success"`
	MovedCount int  `json:"movedCount"`
}

// TransformOptions represents image transformation options.
type TransformOptions struct {
	Width      *int    `json:"width,omitempty"`
	Height     *int    `json:"height,omitempty"`
	Fit        *string `json:"fit,omitempty"`
	Format     *string `json:"format,omitempty"`
	Quality    *int    `json:"quality,omitempty"`
	Crop       *string `json:"crop,omitempty"`
	CropX      *int    `json:"cropX,omitempty"`
	CropY      *int    `json:"cropY,omitempty"`
	CropWidth  *int    `json:"cropWidth,omitempty"`
	CropHeight *int    `json:"cropHeight,omitempty"`
	Blur       *int    `json:"blur,omitempty"`
	Sharpen    *int    `json:"sharpen,omitempty"`
	Brightness *int    `json:"brightness,omitempty"`
	Saturation *int    `json:"saturation,omitempty"`
	Grayscale  bool    `json:"grayscale,omitempty"`
	Rotate     *int    `json:"rotate,omitempty"`
	Flip       bool    `json:"flip,omitempty"`
	Flop       bool    `json:"flop,omitempty"`
}

// ImageWatermarkPosition represents the position of a watermark.
type ImageWatermarkPosition string

const (
	WatermarkTopLeft      ImageWatermarkPosition = "top-left"
	WatermarkTopCenter    ImageWatermarkPosition = "top-center"
	WatermarkTopRight     ImageWatermarkPosition = "top-right"
	WatermarkCenterLeft   ImageWatermarkPosition = "center-left"
	WatermarkCenter       ImageWatermarkPosition = "center"
	WatermarkCenterRight  ImageWatermarkPosition = "center-right"
	WatermarkBottomLeft   ImageWatermarkPosition = "bottom-left"
	WatermarkBottomCenter ImageWatermarkPosition = "bottom-center"
	WatermarkBottomRight  ImageWatermarkPosition = "bottom-right"
)

// ImageWatermarkSizingMode represents the sizing mode for watermarks.
type ImageWatermarkSizingMode string

const (
	WatermarkSizingAbsolute ImageWatermarkSizingMode = "absolute"
	WatermarkSizingRelative ImageWatermarkSizingMode = "relative"
)

// ImageWatermarkConfig represents watermark configuration.
type ImageWatermarkConfig struct {
	AssetID     *string                   `json:"assetId,omitempty"`
	URL         *string                   `json:"url,omitempty"`
	Position    *ImageWatermarkPosition   `json:"position,omitempty"`
	OffsetX     *int                      `json:"offsetX,omitempty"`
	OffsetY     *int                      `json:"offsetY,omitempty"`
	SizingMode  *ImageWatermarkSizingMode `json:"sizingMode,omitempty"`
	Width       *int                      `json:"width,omitempty"`
	Height      *int                      `json:"height,omitempty"`
	Opacity     *int                      `json:"opacity,omitempty"`
	Rotation    *int                      `json:"rotation,omitempty"`
	Tile        bool                      `json:"tile,omitempty"`
	TileSpacing *int                      `json:"tileSpacing,omitempty"`
	BorderRadius *int                     `json:"borderRadius,omitempty"`
}

// FolderTreeNode represents a folder in the tree.
type FolderTreeNode struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Path       string           `json:"path"`
	AssetCount int              `json:"assetCount"`
	Children   []FolderTreeNode `json:"children"`
}

// GetFolderTreeRequest is the request for getting the folder tree.
type GetFolderTreeRequest struct {
	ProjectSlug string `json:"projectSlug"`
	MaxDepth    *int   `json:"maxDepth,omitempty"`
}

// CreateFolderRequest is the request for creating a folder.
type CreateFolderRequest struct {
	ProjectSlug string  `json:"projectSlug"`
	Name        string  `json:"name"`
	ParentID    *string `json:"parentId,omitempty"`
}

// Folder represents a CDN folder.
type Folder struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Path       string     `json:"path"`
	ParentID   *string    `json:"parentId,omitempty"`
	AssetCount int        `json:"assetCount"`
	TotalSize  int64      `json:"totalSize"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

// UpdateFolderRequest is the request for updating a folder.
type UpdateFolderRequest struct {
	ID   string  `json:"id"`
	Name *string `json:"name,omitempty"`
}

// ListFoldersRequest is the request for listing folders.
type ListFoldersRequest struct {
	ParentID *string `json:"parentId,omitempty"`
	Limit    *int    `json:"limit,omitempty"`
	Offset   *int    `json:"offset,omitempty"`
	Search   *string `json:"search,omitempty"`
}

// FolderListItem represents a folder in list responses.
type FolderListItem struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	ParentID   *string   `json:"parentId,omitempty"`
	AssetCount int       `json:"assetCount"`
	TotalSize  int64     `json:"totalSize"`
	CreatedAt  time.Time `json:"createdAt"`
}

// ListFoldersResponse is the response from listing folders.
type ListFoldersResponse struct {
	Folders []FolderListItem `json:"folders"`
	Total   int              `json:"total"`
	HasMore bool             `json:"hasMore"`
}

// MoveFolderRequest is the request for moving a folder.
type MoveFolderRequest struct {
	ID          string  `json:"id"`
	NewParentID *string `json:"newParentId"`
}

// MoveFolderResponse is the response from moving a folder.
type MoveFolderResponse struct {
	Success bool `json:"success"`
}

// SuccessResponse is a generic success response.
type SuccessResponse struct {
	Success bool `json:"success"`
}
