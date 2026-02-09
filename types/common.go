package types

// Environment represents the deployment environment.
type Environment string

const (
	EnvironmentSandbox    Environment = "sandbox"
	EnvironmentProduction Environment = "production"
)

// BatchJobStatus represents the status of a batch job.
type BatchJobStatus string

const (
	BatchJobStatusPending    BatchJobStatus = "pending"
	BatchJobStatusProcessing BatchJobStatus = "processing"
	BatchJobStatusCompleted  BatchJobStatus = "completed"
	BatchJobStatusFailed     BatchJobStatus = "failed"
	BatchJobStatusCancelled  BatchJobStatus = "cancelled"
)

// ScheduleFrequency represents how often a scheduled job runs.
type ScheduleFrequency string

const (
	ScheduleFrequencyHourly  ScheduleFrequency = "hourly"
	ScheduleFrequencyDaily   ScheduleFrequency = "daily"
	ScheduleFrequencyWeekly  ScheduleFrequency = "weekly"
	ScheduleFrequencyMonthly ScheduleFrequency = "monthly"
)

// SuccessResponse represents a simple success response.
type SuccessResponse struct {
	Success bool `json:"success"`
}

// PaginatedRequest contains common pagination parameters.
type PaginatedRequest struct {
	Limit  *int `url:"limit,omitempty"`
	Offset *int `url:"offset,omitempty"`
}

// PaginatedResponse contains common pagination response fields.
type PaginatedResponse struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// GetBatchJobRequest is the request to get a batch job.
type GetBatchJobRequest struct {
	ID          string
	Environment *Environment
	ProjectID   *string
}

// ListBatchJobsRequest is the request to list batch jobs.
type ListBatchJobsRequest struct {
	Environment *Environment
	ProjectID   *string
	Status      *BatchJobStatus
	Limit       *int
	Cursor      *string
}

// CreateBatchResponse is the response when creating a batch job.
type CreateBatchResponse struct {
	ID        string         `json:"id"`
	Status    BatchJobStatus `json:"status"`
	TotalURLs int            `json:"totalUrls"`
}

// GetScheduleRequest is the request to get a schedule.
type GetScheduleRequest struct {
	ID          string
	Environment *Environment
	ProjectID   *string
}

// ListSchedulesRequest is the request to list schedules.
type ListSchedulesRequest struct {
	Environment *Environment
	ProjectID   *string
	IsActive    *bool
	Limit       *int
	Cursor      *string
}

// CreateScheduleResponse is the response when creating a schedule.
type CreateScheduleResponse struct {
	ID string `json:"id"`
}
