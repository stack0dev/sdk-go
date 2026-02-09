package cdn

import "time"

// ImportJobStatus represents the status of an import job.
type ImportJobStatus string

const (
	ImportJobStatusPending    ImportJobStatus = "pending"
	ImportJobStatusValidating ImportJobStatus = "validating"
	ImportJobStatusImporting  ImportJobStatus = "importing"
	ImportJobStatusCompleted  ImportJobStatus = "completed"
	ImportJobStatusFailed     ImportJobStatus = "failed"
	ImportJobStatusCancelled  ImportJobStatus = "cancelled"
)

// ImportAuthType represents the authentication type for imports.
type ImportAuthType string

const (
	ImportAuthTypeIAMCredentials  ImportAuthType = "iam_credentials"
	ImportAuthTypeRoleAssumption ImportAuthType = "role_assumption"
)

// ImportPathMode represents how paths are handled during import.
type ImportPathMode string

const (
	ImportPathModePreserve ImportPathMode = "preserve"
	ImportPathModeFlatten  ImportPathMode = "flatten"
)

// ImportFileStatus represents the status of an import file.
type ImportFileStatus string

const (
	ImportFileStatusPending   ImportFileStatus = "pending"
	ImportFileStatusImporting ImportFileStatus = "importing"
	ImportFileStatusCompleted ImportFileStatus = "completed"
	ImportFileStatusFailed    ImportFileStatus = "failed"
	ImportFileStatusSkipped   ImportFileStatus = "skipped"
)

// ImportError represents an error during import.
type ImportError struct {
	Key       string `json:"key"`
	Error     string `json:"error"`
	Timestamp string `json:"timestamp"`
}

// CreateImportRequest is the request for creating an import job.
type CreateImportRequest struct {
	ProjectSlug     string          `json:"projectSlug"`
	Environment     *CdnEnvironment `json:"environment,omitempty"`
	SourceBucket    string          `json:"sourceBucket"`
	SourceRegion    string          `json:"sourceRegion"`
	SourcePrefix    *string         `json:"sourcePrefix,omitempty"`
	AuthType        ImportAuthType  `json:"authType"`
	AccessKeyID     *string         `json:"accessKeyId,omitempty"`
	SecretAccessKey *string         `json:"secretAccessKey,omitempty"`
	RoleARN         *string         `json:"roleArn,omitempty"`
	ExternalID      *string         `json:"externalId,omitempty"`
	PathMode        *ImportPathMode `json:"pathMode,omitempty"`
	TargetFolder    *string         `json:"targetFolder,omitempty"`
	NotifyEmail     *string         `json:"notifyEmail,omitempty"`
}

// CreateImportResponse is the response from creating an import job.
type CreateImportResponse struct {
	ImportID     string          `json:"importId"`
	Status       ImportJobStatus `json:"status"`
	SourceBucket string          `json:"sourceBucket"`
	SourceRegion string          `json:"sourceRegion"`
	SourcePrefix *string         `json:"sourcePrefix,omitempty"`
	CreatedAt    time.Time       `json:"createdAt"`
}

// ImportJob represents an import job.
type ImportJob struct {
	ID             string          `json:"id"`
	OrganizationID string          `json:"organizationId"`
	ProjectID      string          `json:"projectId"`
	Environment    CdnEnvironment  `json:"environment"`
	SourceBucket   string          `json:"sourceBucket"`
	SourceRegion   string          `json:"sourceRegion"`
	SourcePrefix   *string         `json:"sourcePrefix,omitempty"`
	AuthType       ImportAuthType  `json:"authType"`
	PathMode       ImportPathMode  `json:"pathMode"`
	TargetFolder   *string         `json:"targetFolder,omitempty"`
	Status         ImportJobStatus `json:"status"`
	TotalFiles     int             `json:"totalFiles"`
	ProcessedFiles int             `json:"processedFiles"`
	SkippedFiles   int             `json:"skippedFiles"`
	FailedFiles    int             `json:"failedFiles"`
	TotalBytes     int64           `json:"totalBytes"`
	ProcessedBytes int64           `json:"processedBytes"`
	Errors         []ImportError   `json:"errors,omitempty"`
	NotifyEmail    *string         `json:"notifyEmail,omitempty"`
	StartedAt      *time.Time      `json:"startedAt,omitempty"`
	CompletedAt    *time.Time      `json:"completedAt,omitempty"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      *time.Time      `json:"updatedAt,omitempty"`
}

// ListImportsRequest is the request for listing import jobs.
type ListImportsRequest struct {
	ProjectSlug string           `json:"projectSlug"`
	Environment *CdnEnvironment  `json:"environment,omitempty"`
	Status      *ImportJobStatus `json:"status,omitempty"`
	SortBy      *string          `json:"sortBy,omitempty"`
	SortOrder   *string          `json:"sortOrder,omitempty"`
	Limit       *int             `json:"limit,omitempty"`
	Offset      *int             `json:"offset,omitempty"`
}

// ImportJobSummary represents an import job summary.
type ImportJobSummary struct {
	ID             string          `json:"id"`
	SourceBucket   string          `json:"sourceBucket"`
	SourceRegion   string          `json:"sourceRegion"`
	SourcePrefix   *string         `json:"sourcePrefix,omitempty"`
	Status         ImportJobStatus `json:"status"`
	TotalFiles     int             `json:"totalFiles"`
	ProcessedFiles int             `json:"processedFiles"`
	SkippedFiles   int             `json:"skippedFiles"`
	FailedFiles    int             `json:"failedFiles"`
	TotalBytes     int64           `json:"totalBytes"`
	ProcessedBytes int64           `json:"processedBytes"`
	StartedAt      *time.Time      `json:"startedAt,omitempty"`
	CompletedAt    *time.Time      `json:"completedAt,omitempty"`
	CreatedAt      time.Time       `json:"createdAt"`
}

// ListImportsResponse is the response from listing import jobs.
type ListImportsResponse struct {
	Imports []ImportJobSummary `json:"imports"`
	Total   int                `json:"total"`
	HasMore bool               `json:"hasMore"`
}

// CancelImportResponse is the response from cancelling an import.
type CancelImportResponse struct {
	Success bool            `json:"success"`
	Status  ImportJobStatus `json:"status"`
}

// RetryImportResponse is the response from retrying an import.
type RetryImportResponse struct {
	Success      bool            `json:"success"`
	RetriedCount int             `json:"retriedCount"`
	Status       ImportJobStatus `json:"status"`
}

// ListImportFilesRequest is the request for listing import files.
type ListImportFilesRequest struct {
	ImportID  string            `json:"importId"`
	Status    *ImportFileStatus `json:"status,omitempty"`
	SortBy    *string           `json:"sortBy,omitempty"`
	SortOrder *string           `json:"sortOrder,omitempty"`
	Limit     *int              `json:"limit,omitempty"`
	Offset    *int              `json:"offset,omitempty"`
}

// ImportFile represents a file in an import job.
type ImportFile struct {
	ID             string           `json:"id"`
	ImportJobID    string           `json:"importJobId"`
	SourceKey      string           `json:"sourceKey"`
	SourceSize     int64            `json:"sourceSize"`
	SourceMimeType *string          `json:"sourceMimeType,omitempty"`
	SourceEtag     *string          `json:"sourceEtag,omitempty"`
	AssetID        *string          `json:"assetId,omitempty"`
	Status         ImportFileStatus `json:"status"`
	ErrorMessage   *string          `json:"errorMessage,omitempty"`
	RetryCount     int              `json:"retryCount"`
	LastAttemptAt  *time.Time       `json:"lastAttemptAt,omitempty"`
	CreatedAt      time.Time        `json:"createdAt"`
	CompletedAt    *time.Time       `json:"completedAt,omitempty"`
}

// ListImportFilesResponse is the response from listing import files.
type ListImportFilesResponse struct {
	Files   []ImportFile `json:"files"`
	Total   int          `json:"total"`
	HasMore bool         `json:"hasMore"`
}
