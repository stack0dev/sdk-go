package cdn

import "time"

// PrivateFileStatus represents the status of a private file.
type PrivateFileStatus string

const (
	PrivateFileStatusPending PrivateFileStatus = "pending"
	PrivateFileStatusReady   PrivateFileStatus = "ready"
	PrivateFileStatusFailed  PrivateFileStatus = "failed"
	PrivateFileStatusDeleted PrivateFileStatus = "deleted"
)

// PrivateFile represents a private file.
type PrivateFile struct {
	ID               string                 `json:"id"`
	Filename         string                 `json:"filename"`
	OriginalFilename string                 `json:"originalFilename"`
	MimeType         string                 `json:"mimeType"`
	Size             int64                  `json:"size"`
	S3Key            string                 `json:"s3Key"`
	Folder           *string                `json:"folder,omitempty"`
	Description      *string                `json:"description,omitempty"`
	Tags             []string               `json:"tags,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	Status           PrivateFileStatus      `json:"status"`
	CreatedAt        time.Time              `json:"createdAt"`
	UpdatedAt        *time.Time             `json:"updatedAt,omitempty"`
}

// PrivateUploadURLRequest is the request for getting a private upload URL.
type PrivateUploadURLRequest struct {
	ProjectSlug string                 `json:"projectSlug"`
	Filename    string                 `json:"filename"`
	MimeType    string                 `json:"mimeType"`
	Size        int64                  `json:"size"`
	Folder      *string                `json:"folder,omitempty"`
	Description *string                `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// PrivateUploadURLResponse is the response from getting a private upload URL.
type PrivateUploadURLResponse struct {
	UploadURL string    `json:"uploadUrl"`
	FileID    string    `json:"fileId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// PrivateDownloadURLRequest is the request for getting a private download URL.
type PrivateDownloadURLRequest struct {
	FileID    string `json:"fileId"`
	ExpiresIn *int   `json:"expiresIn,omitempty"`
}

// PrivateDownloadURLResponse is the response from getting a private download URL.
type PrivateDownloadURLResponse struct {
	DownloadURL string    `json:"downloadUrl"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// ListPrivateFilesRequest is the request for listing private files.
type ListPrivateFilesRequest struct {
	ProjectSlug string             `json:"projectSlug"`
	Folder      *string            `json:"folder,omitempty"`
	Status      *PrivateFileStatus `json:"status,omitempty"`
	Search      *string            `json:"search,omitempty"`
	SortBy      *string            `json:"sortBy,omitempty"`
	SortOrder   *string            `json:"sortOrder,omitempty"`
	Limit       *int               `json:"limit,omitempty"`
	Offset      *int               `json:"offset,omitempty"`
}

// ListPrivateFilesResponse is the response from listing private files.
type ListPrivateFilesResponse struct {
	Files   []PrivateFile `json:"files"`
	Total   int           `json:"total"`
	HasMore bool          `json:"hasMore"`
}

// UpdatePrivateFileRequest is the request for updating a private file.
type UpdatePrivateFileRequest struct {
	FileID      string                 `json:"fileId"`
	Description *string                `json:"description,omitempty"`
	Folder      *string                `json:"folder,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// MovePrivateFilesRequest is the request for moving private files.
type MovePrivateFilesRequest struct {
	FileIDs []string `json:"fileIds"`
	Folder  *string  `json:"folder"`
}

// MovePrivateFilesResponse is the response from moving private files.
type MovePrivateFilesResponse struct {
	Success    bool `json:"success"`
	MovedCount int  `json:"movedCount"`
}

// BundleStatus represents the status of a download bundle.
type BundleStatus string

const (
	BundleStatusPending    BundleStatus = "pending"
	BundleStatusProcessing BundleStatus = "processing"
	BundleStatusReady      BundleStatus = "ready"
	BundleStatusFailed     BundleStatus = "failed"
	BundleStatusExpired    BundleStatus = "expired"
)

// DownloadBundle represents a download bundle.
type DownloadBundle struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	Description    *string      `json:"description,omitempty"`
	AssetIDs       []string     `json:"assetIds,omitempty"`
	PrivateFileIDs []string     `json:"privateFileIds,omitempty"`
	S3Key          *string      `json:"s3Key,omitempty"`
	Size           *int64       `json:"size,omitempty"`
	FileCount      *int         `json:"fileCount,omitempty"`
	Status         BundleStatus `json:"status"`
	Error          *string      `json:"error,omitempty"`
	ExpiresAt      *time.Time   `json:"expiresAt,omitempty"`
	CreatedAt      time.Time    `json:"createdAt"`
	CompletedAt    *time.Time   `json:"completedAt,omitempty"`
}

// CreateBundleRequest is the request for creating a bundle.
type CreateBundleRequest struct {
	ProjectSlug    string   `json:"projectSlug"`
	Name           string   `json:"name"`
	Description    *string  `json:"description,omitempty"`
	AssetIDs       []string `json:"assetIds,omitempty"`
	PrivateFileIDs []string `json:"privateFileIds,omitempty"`
	ExpiresIn      *int     `json:"expiresIn,omitempty"`
}

// CreateBundleResponse is the response from creating a bundle.
type CreateBundleResponse struct {
	Bundle DownloadBundle `json:"bundle"`
}

// ListBundlesRequest is the request for listing bundles.
type ListBundlesRequest struct {
	ProjectSlug string        `json:"projectSlug"`
	Status      *BundleStatus `json:"status,omitempty"`
	Search      *string       `json:"search,omitempty"`
	Limit       *int          `json:"limit,omitempty"`
	Offset      *int          `json:"offset,omitempty"`
}

// ListBundlesResponse is the response from listing bundles.
type ListBundlesResponse struct {
	Bundles []DownloadBundle `json:"bundles"`
	Total   int              `json:"total"`
	HasMore bool             `json:"hasMore"`
}

// BundleDownloadURLRequest is the request for getting a bundle download URL.
type BundleDownloadURLRequest struct {
	BundleID  string `json:"bundleId"`
	ExpiresIn *int   `json:"expiresIn,omitempty"`
}

// BundleDownloadURLResponse is the response from getting a bundle download URL.
type BundleDownloadURLResponse struct {
	DownloadURL string    `json:"downloadUrl"`
	ExpiresAt   time.Time `json:"expiresAt"`
}
