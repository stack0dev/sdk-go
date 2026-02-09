package extraction

import (
	"time"

	"github.com/stack0/sdk-go/types"
)

// ExtractionStatus represents the status of an extraction.
type ExtractionStatus string

const (
	ExtractionStatusPending    ExtractionStatus = "pending"
	ExtractionStatusProcessing ExtractionStatus = "processing"
	ExtractionStatusCompleted  ExtractionStatus = "completed"
	ExtractionStatusFailed     ExtractionStatus = "failed"
)

// ExtractionMode represents the extraction mode.
type ExtractionMode string

const (
	ExtractionModeAuto     ExtractionMode = "auto"
	ExtractionModeSchema   ExtractionMode = "schema"
	ExtractionModeMarkdown ExtractionMode = "markdown"
	ExtractionModeRaw      ExtractionMode = "raw"
)

// PageMetadata represents extracted page metadata.
type PageMetadata struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	OgImage     *string  `json:"ogImage,omitempty"`
	Favicon     *string  `json:"favicon,omitempty"`
	Links       []string `json:"links,omitempty"`
	Images      []string `json:"images,omitempty"`
}

// ExtractionResult represents an extraction result.
type ExtractionResult struct {
	ID               string                 `json:"id"`
	OrganizationID   string                 `json:"organizationId"`
	ProjectID        *string                `json:"projectId,omitempty"`
	Environment      types.Environment      `json:"environment"`
	URL              string                 `json:"url"`
	Mode             string                 `json:"mode"`
	Status           ExtractionStatus       `json:"status"`
	ExtractedData    map[string]interface{} `json:"extractedData,omitempty"`
	Markdown         *string                `json:"markdown,omitempty"`
	RawHTML          *string                `json:"rawHtml,omitempty"`
	PageMetadata     *PageMetadata          `json:"pageMetadata,omitempty"`
	Error            *string                `json:"error,omitempty"`
	ProcessingTimeMs *int64                 `json:"processingTimeMs,omitempty"`
	TokensUsed       *int                   `json:"tokensUsed,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt        time.Time              `json:"createdAt"`
	CompletedAt      *time.Time             `json:"completedAt,omitempty"`
}

// Cookie represents a browser cookie.
type Cookie struct {
	Name   string  `json:"name"`
	Value  string  `json:"value"`
	Domain *string `json:"domain,omitempty"`
}

// CreateExtractionRequest is the request for extracting content.
type CreateExtractionRequest struct {
	URL             string                 `json:"url"`
	Environment     *types.Environment     `json:"environment,omitempty"`
	ProjectID       *string                `json:"projectId,omitempty"`
	Mode            *ExtractionMode        `json:"mode,omitempty"`
	Schema          map[string]interface{} `json:"schema,omitempty"`
	Prompt          *string                `json:"prompt,omitempty"`
	IncludeLinks    *bool                  `json:"includeLinks,omitempty"`
	IncludeImages   *bool                  `json:"includeImages,omitempty"`
	IncludeMetadata *bool                  `json:"includeMetadata,omitempty"`
	WaitForSelector *string                `json:"waitForSelector,omitempty"`
	WaitForTimeout  *int                   `json:"waitForTimeout,omitempty"`
	Headers         map[string]string      `json:"headers,omitempty"`
	Cookies         []Cookie               `json:"cookies,omitempty"`
	WebhookURL      *string                `json:"webhookUrl,omitempty"`
	WebhookSecret   *string                `json:"webhookSecret,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// CreateExtractionResponse is the response from extracting content.
type CreateExtractionResponse struct {
	ID     string           `json:"id"`
	Status ExtractionStatus `json:"status"`
}

// GetExtractionRequest is the request for getting an extraction.
type GetExtractionRequest struct {
	ID          string             `json:"id"`
	Environment *types.Environment `json:"environment,omitempty"`
	ProjectID   *string            `json:"projectId,omitempty"`
}

// ListExtractionsRequest is the request for listing extractions.
type ListExtractionsRequest struct {
	Environment *types.Environment `json:"environment,omitempty"`
	ProjectID   *string            `json:"projectId,omitempty"`
	Status      *ExtractionStatus  `json:"status,omitempty"`
	URL         *string            `json:"url,omitempty"`
	Limit       *int               `json:"limit,omitempty"`
	Cursor      *string            `json:"cursor,omitempty"`
}

// ListExtractionsResponse is the response from listing extractions.
type ListExtractionsResponse struct {
	Items      []ExtractionResult `json:"items"`
	NextCursor *string            `json:"nextCursor,omitempty"`
}

// BatchExtractionJob represents a batch extraction job.
type BatchExtractionJob struct {
	ID             string                 `json:"id"`
	OrganizationID string                 `json:"organizationId"`
	ProjectID      *string                `json:"projectId,omitempty"`
	Environment    types.Environment      `json:"environment"`
	Type           string                 `json:"type"`
	Name           *string                `json:"name,omitempty"`
	Status         types.BatchJobStatus   `json:"status"`
	URLs           []string               `json:"urls"`
	Config         map[string]interface{} `json:"config"`
	TotalURLs      int                    `json:"totalUrls"`
	ProcessedURLs  int                    `json:"processedUrls"`
	SuccessfulURLs int                    `json:"successfulUrls"`
	FailedURLs     int                    `json:"failedUrls"`
	WebhookURL     *string                `json:"webhookUrl,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt      time.Time              `json:"createdAt"`
	StartedAt      *time.Time             `json:"startedAt,omitempty"`
	CompletedAt    *time.Time             `json:"completedAt,omitempty"`
}

// BatchExtractionConfig represents batch extraction configuration.
type BatchExtractionConfig struct {
	Mode            *ExtractionMode        `json:"mode,omitempty"`
	Schema          map[string]interface{} `json:"schema,omitempty"`
	Prompt          *string                `json:"prompt,omitempty"`
	IncludeLinks    *bool                  `json:"includeLinks,omitempty"`
	IncludeImages   *bool                  `json:"includeImages,omitempty"`
	IncludeMetadata *bool                  `json:"includeMetadata,omitempty"`
	WaitForSelector *string                `json:"waitForSelector,omitempty"`
	WaitForTimeout  *int                   `json:"waitForTimeout,omitempty"`
}

// CreateBatchExtractionsRequest is the request for creating a batch job.
type CreateBatchExtractionsRequest struct {
	URLs          []string                `json:"urls"`
	Environment   *types.Environment      `json:"environment,omitempty"`
	ProjectID     *string                 `json:"projectId,omitempty"`
	Name          *string                 `json:"name,omitempty"`
	Config        *BatchExtractionConfig  `json:"config,omitempty"`
	WebhookURL    *string                 `json:"webhookUrl,omitempty"`
	WebhookSecret *string                 `json:"webhookSecret,omitempty"`
	Metadata      map[string]interface{}  `json:"metadata,omitempty"`
}

// CreateBatchResponse is the response from creating a batch job.
type CreateBatchResponse struct {
	ID        string `json:"id"`
	TotalURLs int    `json:"totalUrls"`
}

// GetBatchJobRequest is the request for getting a batch job.
type GetBatchJobRequest struct {
	ID          string             `json:"id"`
	Environment *types.Environment `json:"environment,omitempty"`
	ProjectID   *string            `json:"projectId,omitempty"`
}

// ListBatchJobsRequest is the request for listing batch jobs.
type ListBatchJobsRequest struct {
	Environment *types.Environment    `json:"environment,omitempty"`
	ProjectID   *string               `json:"projectId,omitempty"`
	Status      *types.BatchJobStatus `json:"status,omitempty"`
	Limit       *int                  `json:"limit,omitempty"`
	Cursor      *string               `json:"cursor,omitempty"`
}

// BatchJobsResponse is the response from listing batch jobs.
type BatchJobsResponse struct {
	Items      []BatchExtractionJob `json:"items"`
	NextCursor *string              `json:"nextCursor,omitempty"`
}

// ExtractionSchedule represents an extraction schedule.
type ExtractionSchedule struct {
	ID              string                  `json:"id"`
	OrganizationID  string                  `json:"organizationId"`
	ProjectID       *string                 `json:"projectId,omitempty"`
	Environment     types.Environment       `json:"environment"`
	Name            string                  `json:"name"`
	URL             string                  `json:"url"`
	Type            string                  `json:"type"`
	Frequency       types.ScheduleFrequency `json:"frequency"`
	Config          map[string]interface{}  `json:"config"`
	IsActive        bool                    `json:"isActive"`
	DetectChanges   bool                    `json:"detectChanges"`
	ChangeThreshold *int                    `json:"changeThreshold,omitempty"`
	WebhookURL      *string                 `json:"webhookUrl,omitempty"`
	TotalRuns       int                     `json:"totalRuns"`
	SuccessfulRuns  int                     `json:"successfulRuns"`
	FailedRuns      int                     `json:"failedRuns"`
	LastRunAt       *time.Time              `json:"lastRunAt,omitempty"`
	NextRunAt       *time.Time              `json:"nextRunAt,omitempty"`
	Metadata        map[string]interface{}  `json:"metadata,omitempty"`
	CreatedAt       time.Time               `json:"createdAt"`
	UpdatedAt       time.Time               `json:"updatedAt"`
}

// CreateExtractionScheduleRequest is the request for creating a schedule.
type CreateExtractionScheduleRequest struct {
	Name            string                   `json:"name"`
	URL             string                   `json:"url"`
	Environment     *types.Environment       `json:"environment,omitempty"`
	ProjectID       *string                  `json:"projectId,omitempty"`
	Frequency       *types.ScheduleFrequency `json:"frequency,omitempty"`
	Config          *BatchExtractionConfig   `json:"config,omitempty"`
	DetectChanges   *bool                    `json:"detectChanges,omitempty"`
	ChangeThreshold *int                     `json:"changeThreshold,omitempty"`
	WebhookURL      *string                  `json:"webhookUrl,omitempty"`
	WebhookSecret   *string                  `json:"webhookSecret,omitempty"`
	Metadata        map[string]interface{}   `json:"metadata,omitempty"`
}

// UpdateExtractionScheduleRequest is the request for updating a schedule.
type UpdateExtractionScheduleRequest struct {
	ID              string                   `json:"id"`
	Environment     *types.Environment       `json:"environment,omitempty"`
	ProjectID       *string                  `json:"projectId,omitempty"`
	Name            *string                  `json:"name,omitempty"`
	Frequency       *types.ScheduleFrequency `json:"frequency,omitempty"`
	Config          map[string]interface{}   `json:"config,omitempty"`
	IsActive        *bool                    `json:"isActive,omitempty"`
	DetectChanges   *bool                    `json:"detectChanges,omitempty"`
	ChangeThreshold *int                     `json:"changeThreshold,omitempty"`
	WebhookURL      *string                  `json:"webhookUrl,omitempty"`
	WebhookSecret   *string                  `json:"webhookSecret,omitempty"`
	Metadata        map[string]interface{}   `json:"metadata,omitempty"`
}

// CreateScheduleResponse is the response from creating a schedule.
type CreateScheduleResponse struct {
	ID string `json:"id"`
}

// GetScheduleRequest is the request for getting a schedule.
type GetScheduleRequest struct {
	ID          string             `json:"id"`
	Environment *types.Environment `json:"environment,omitempty"`
	ProjectID   *string            `json:"projectId,omitempty"`
}

// ListSchedulesRequest is the request for listing schedules.
type ListSchedulesRequest struct {
	Environment *types.Environment `json:"environment,omitempty"`
	ProjectID   *string            `json:"projectId,omitempty"`
	IsActive    *bool              `json:"isActive,omitempty"`
	Limit       *int               `json:"limit,omitempty"`
	Cursor      *string            `json:"cursor,omitempty"`
}

// SchedulesResponse is the response from listing schedules.
type SchedulesResponse struct {
	Items      []ExtractionSchedule `json:"items"`
	NextCursor *string              `json:"nextCursor,omitempty"`
}

// SuccessResponse is a generic success response.
type SuccessResponse struct {
	Success bool `json:"success"`
}

// ToggleResponse is the response from toggling a schedule.
type ToggleResponse struct {
	IsActive bool `json:"isActive"`
}

// ExtractionUsage represents extraction usage stats.
type ExtractionUsage struct {
	PeriodStart              time.Time `json:"periodStart"`
	PeriodEnd                time.Time `json:"periodEnd"`
	ExtractionsTotal         int       `json:"extractionsTotal"`
	ExtractionsSuccessful    int       `json:"extractionsSuccessful"`
	ExtractionsFailed        int       `json:"extractionsFailed"`
	ExtractionCreditsUsed    int       `json:"extractionCreditsUsed"`
	ExtractionTokensUsed     int       `json:"extractionTokensUsed"`
}

// GetUsageRequest is the request for getting usage stats.
type GetUsageRequest struct {
	Environment *types.Environment `json:"environment,omitempty"`
	PeriodStart *string            `json:"periodStart,omitempty"`
	PeriodEnd   *string            `json:"periodEnd,omitempty"`
}

// DailyUsageItem represents daily usage.
type DailyUsageItem struct {
	Date        string `json:"date"`
	Screenshots int    `json:"screenshots"`
	Extractions int    `json:"extractions"`
	CreditsUsed int    `json:"creditsUsed"`
}

// GetDailyUsageResponse is the response from getting daily usage.
type GetDailyUsageResponse struct {
	Days []DailyUsageItem `json:"days"`
}
