package screenshots

import (
	"time"

	"github.com/stack0/sdk-go/types"
)

// ScreenshotStatus represents the status of a screenshot.
type ScreenshotStatus string

const (
	ScreenshotStatusPending    ScreenshotStatus = "pending"
	ScreenshotStatusProcessing ScreenshotStatus = "processing"
	ScreenshotStatusCompleted  ScreenshotStatus = "completed"
	ScreenshotStatusFailed     ScreenshotStatus = "failed"
)

// ScreenshotFormat represents the output format.
type ScreenshotFormat string

const (
	ScreenshotFormatPNG  ScreenshotFormat = "png"
	ScreenshotFormatJPEG ScreenshotFormat = "jpeg"
	ScreenshotFormatWebP ScreenshotFormat = "webp"
	ScreenshotFormatPDF  ScreenshotFormat = "pdf"
)

// DeviceType represents the device type for rendering.
type DeviceType string

const (
	DeviceTypeDesktop DeviceType = "desktop"
	DeviceTypeTablet  DeviceType = "tablet"
	DeviceTypeMobile  DeviceType = "mobile"
)

// ResourceType represents blocked resource types.
type ResourceType string

const (
	ResourceTypeImage      ResourceType = "image"
	ResourceTypeStylesheet ResourceType = "stylesheet"
	ResourceTypeScript     ResourceType = "script"
	ResourceTypeFont       ResourceType = "font"
	ResourceTypeMedia      ResourceType = "media"
	ResourceTypeXHR        ResourceType = "xhr"
	ResourceTypeFetch      ResourceType = "fetch"
	ResourceTypeWebSocket  ResourceType = "websocket"
)

// Screenshot represents a screenshot result.
type Screenshot struct {
	ID               string                 `json:"id"`
	OrganizationID   string                 `json:"organizationId"`
	ProjectID        *string                `json:"projectId,omitempty"`
	Environment      types.Environment      `json:"environment"`
	URL              string                 `json:"url"`
	Format           ScreenshotFormat       `json:"format"`
	Quality          *int                   `json:"quality,omitempty"`
	FullPage         bool                   `json:"fullPage"`
	DeviceType       DeviceType             `json:"deviceType"`
	ViewportWidth    *int                   `json:"viewportWidth,omitempty"`
	ViewportHeight   *int                   `json:"viewportHeight,omitempty"`
	Status           ScreenshotStatus       `json:"status"`
	ImageURL         *string                `json:"imageUrl,omitempty"`
	ImageSize        *int64                 `json:"imageSize,omitempty"`
	ImageWidth       *int                   `json:"imageWidth,omitempty"`
	ImageHeight      *int                   `json:"imageHeight,omitempty"`
	Error            *string                `json:"error,omitempty"`
	ProcessingTimeMs *int64                 `json:"processingTimeMs,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt        time.Time              `json:"createdAt"`
	CompletedAt      *time.Time             `json:"completedAt,omitempty"`
}

// Clip represents a clip region for screenshots.
type Clip struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Cookie represents a browser cookie.
type Cookie struct {
	Name   string  `json:"name"`
	Value  string  `json:"value"`
	Domain *string `json:"domain,omitempty"`
}

// CreateScreenshotRequest is the request for capturing a screenshot.
type CreateScreenshotRequest struct {
	URL                string                 `json:"url"`
	Environment        *types.Environment     `json:"environment,omitempty"`
	ProjectID          *string                `json:"projectId,omitempty"`
	Format             *ScreenshotFormat      `json:"format,omitempty"`
	Quality            *int                   `json:"quality,omitempty"`
	FullPage           *bool                  `json:"fullPage,omitempty"`
	DeviceType         *DeviceType            `json:"deviceType,omitempty"`
	ViewportWidth      *int                   `json:"viewportWidth,omitempty"`
	ViewportHeight     *int                   `json:"viewportHeight,omitempty"`
	DeviceScaleFactor  *int                   `json:"deviceScaleFactor,omitempty"`
	WaitForSelector    *string                `json:"waitForSelector,omitempty"`
	WaitForTimeout     *int                   `json:"waitForTimeout,omitempty"`
	BlockAds           *bool                  `json:"blockAds,omitempty"`
	BlockCookieBanners *bool                  `json:"blockCookieBanners,omitempty"`
	BlockChatWidgets   *bool                  `json:"blockChatWidgets,omitempty"`
	BlockTrackers      *bool                  `json:"blockTrackers,omitempty"`
	BlockURLs          []string               `json:"blockUrls,omitempty"`
	BlockResources     []ResourceType         `json:"blockResources,omitempty"`
	DarkMode           *bool                  `json:"darkMode,omitempty"`
	CustomCSS          *string                `json:"customCss,omitempty"`
	CustomJS           *string                `json:"customJs,omitempty"`
	Headers            map[string]string      `json:"headers,omitempty"`
	Cookies            []Cookie               `json:"cookies,omitempty"`
	Selector           *string                `json:"selector,omitempty"`
	HideSelectors      []string               `json:"hideSelectors,omitempty"`
	ClickSelector      *string                `json:"clickSelector,omitempty"`
	OmitBackground     *bool                  `json:"omitBackground,omitempty"`
	UserAgent          *string                `json:"userAgent,omitempty"`
	Clip               *Clip                  `json:"clip,omitempty"`
	ThumbnailWidth     *int                   `json:"thumbnailWidth,omitempty"`
	ThumbnailHeight    *int                   `json:"thumbnailHeight,omitempty"`
	CacheKey           *string                `json:"cacheKey,omitempty"`
	CacheTTL           *int                   `json:"cacheTtl,omitempty"`
	WebhookURL         *string                `json:"webhookUrl,omitempty"`
	WebhookSecret      *string                `json:"webhookSecret,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// CreateScreenshotResponse is the response from capturing a screenshot.
type CreateScreenshotResponse struct {
	ID     string           `json:"id"`
	Status ScreenshotStatus `json:"status"`
}

// GetScreenshotRequest is the request for getting a screenshot.
type GetScreenshotRequest struct {
	ID          string             `json:"id"`
	Environment *types.Environment `json:"environment,omitempty"`
	ProjectID   *string            `json:"projectId,omitempty"`
}

// ListScreenshotsRequest is the request for listing screenshots.
type ListScreenshotsRequest struct {
	Environment *types.Environment `json:"environment,omitempty"`
	ProjectID   *string            `json:"projectId,omitempty"`
	Status      *ScreenshotStatus  `json:"status,omitempty"`
	URL         *string            `json:"url,omitempty"`
	Limit       *int               `json:"limit,omitempty"`
	Cursor      *string            `json:"cursor,omitempty"`
}

// ListScreenshotsResponse is the response from listing screenshots.
type ListScreenshotsResponse struct {
	Items      []Screenshot `json:"items"`
	NextCursor *string      `json:"nextCursor,omitempty"`
}

// BatchScreenshotJob represents a batch screenshot job.
type BatchScreenshotJob struct {
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

// BatchScreenshotConfig represents batch screenshot configuration.
type BatchScreenshotConfig struct {
	Format             *ScreenshotFormat `json:"format,omitempty"`
	Quality            *int              `json:"quality,omitempty"`
	FullPage           *bool             `json:"fullPage,omitempty"`
	DeviceType         *DeviceType       `json:"deviceType,omitempty"`
	ViewportWidth      *int              `json:"viewportWidth,omitempty"`
	ViewportHeight     *int              `json:"viewportHeight,omitempty"`
	BlockAds           *bool             `json:"blockAds,omitempty"`
	BlockCookieBanners *bool             `json:"blockCookieBanners,omitempty"`
	WaitForSelector    *string           `json:"waitForSelector,omitempty"`
	WaitForTimeout     *int              `json:"waitForTimeout,omitempty"`
}

// CreateBatchScreenshotsRequest is the request for creating a batch job.
type CreateBatchScreenshotsRequest struct {
	URLs          []string                `json:"urls"`
	Environment   *types.Environment      `json:"environment,omitempty"`
	ProjectID     *string                 `json:"projectId,omitempty"`
	Name          *string                 `json:"name,omitempty"`
	Config        *BatchScreenshotConfig  `json:"config,omitempty"`
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
	Items      []BatchScreenshotJob `json:"items"`
	NextCursor *string              `json:"nextCursor,omitempty"`
}

// ScreenshotSchedule represents a screenshot schedule.
type ScreenshotSchedule struct {
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

// CreateScreenshotScheduleRequest is the request for creating a schedule.
type CreateScreenshotScheduleRequest struct {
	Name            string                   `json:"name"`
	URL             string                   `json:"url"`
	Environment     *types.Environment       `json:"environment,omitempty"`
	ProjectID       *string                  `json:"projectId,omitempty"`
	Frequency       *types.ScheduleFrequency `json:"frequency,omitempty"`
	Config          *BatchScreenshotConfig   `json:"config,omitempty"`
	DetectChanges   *bool                    `json:"detectChanges,omitempty"`
	ChangeThreshold *int                     `json:"changeThreshold,omitempty"`
	WebhookURL      *string                  `json:"webhookUrl,omitempty"`
	WebhookSecret   *string                  `json:"webhookSecret,omitempty"`
	Metadata        map[string]interface{}   `json:"metadata,omitempty"`
}

// UpdateScreenshotScheduleRequest is the request for updating a schedule.
type UpdateScreenshotScheduleRequest struct {
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
	Items      []ScreenshotSchedule `json:"items"`
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
