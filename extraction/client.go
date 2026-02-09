package extraction

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/stack0/sdk-go/client"
	"github.com/stack0/sdk-go/types"
)

// Client handles extraction operations.
type Client struct {
	http *client.HTTPClient
}

// NewClient creates a new extraction client.
func NewClient(http *client.HTTPClient) *Client {
	return &Client{http: http}
}

// Extract extracts content from a URL.
func (c *Client) Extract(ctx context.Context, req *CreateExtractionRequest) (*CreateExtractionResponse, error) {
	var resp CreateExtractionResponse
	if err := c.http.Post(ctx, "/webdata/extractions", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves an extraction by ID.
func (c *Client) Get(ctx context.Context, req *GetExtractionRequest) (*ExtractionResult, error) {
	params := url.Values{}
	if req.Environment != nil {
		params.Set("environment", string(*req.Environment))
	}
	if req.ProjectID != nil {
		params.Set("projectId", *req.ProjectID)
	}

	path := "/webdata/extractions/" + req.ID
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ExtractionResult
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// List lists extractions with pagination and filters.
func (c *Client) List(ctx context.Context, req *ListExtractionsRequest) (*ListExtractionsResponse, error) {
	params := url.Values{}
	if req != nil {
		if req.Environment != nil {
			params.Set("environment", string(*req.Environment))
		}
		if req.ProjectID != nil {
			params.Set("projectId", *req.ProjectID)
		}
		if req.Status != nil {
			params.Set("status", string(*req.Status))
		}
		if req.URL != nil {
			params.Set("url", *req.URL)
		}
		if req.Limit != nil {
			params.Set("limit", strconv.Itoa(*req.Limit))
		}
		if req.Cursor != nil {
			params.Set("cursor", *req.Cursor)
		}
	}

	path := "/webdata/extractions"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListExtractionsResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes an extraction.
func (c *Client) Delete(ctx context.Context, req *GetExtractionRequest) (*SuccessResponse, error) {
	params := url.Values{}
	if req.Environment != nil {
		params.Set("environment", string(*req.Environment))
	}
	if req.ProjectID != nil {
		params.Set("projectId", *req.ProjectID)
	}

	path := "/webdata/extractions/" + req.ID
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	body := map[string]interface{}{"id": req.ID}
	if req.Environment != nil {
		body["environment"] = *req.Environment
	}
	if req.ProjectID != nil {
		body["projectId"] = *req.ProjectID
	}

	var resp SuccessResponse
	if err := c.http.DeleteWithBody(ctx, path, body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ExtractAndWaitOptions are options for ExtractAndWait.
type ExtractAndWaitOptions struct {
	PollInterval time.Duration
	Timeout      time.Duration
}

// ExtractAndWait extracts content and waits for completion.
func (c *Client) ExtractAndWait(ctx context.Context, req *CreateExtractionRequest, opts *ExtractAndWaitOptions) (*ExtractionResult, error) {
	pollInterval := 1 * time.Second
	timeout := 60 * time.Second
	if opts != nil {
		if opts.PollInterval > 0 {
			pollInterval = opts.PollInterval
		}
		if opts.Timeout > 0 {
			timeout = opts.Timeout
		}
	}

	resp, err := c.Extract(ctx, req)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	for time.Since(startTime) < timeout {
		extraction, err := c.Get(ctx, &GetExtractionRequest{
			ID:          resp.ID,
			Environment: req.Environment,
			ProjectID:   req.ProjectID,
		})
		if err != nil {
			return nil, err
		}

		if extraction.Status == ExtractionStatusCompleted || extraction.Status == ExtractionStatusFailed {
			if extraction.Status == ExtractionStatusFailed {
				errMsg := "Extraction failed"
				if extraction.Error != nil {
					errMsg = *extraction.Error
				}
				return nil, errors.New(errMsg)
			}
			return extraction, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(pollInterval):
		}
	}

	return nil, types.NewTimeoutError("Extraction timed out")
}

// Batch creates a batch extraction job for multiple URLs.
func (c *Client) Batch(ctx context.Context, req *CreateBatchExtractionsRequest) (*CreateBatchResponse, error) {
	var resp CreateBatchResponse
	if err := c.http.Post(ctx, "/webdata/batch/extractions", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBatchJob retrieves a batch job by ID.
func (c *Client) GetBatchJob(ctx context.Context, req *GetBatchJobRequest) (*BatchExtractionJob, error) {
	params := url.Values{}
	if req.Environment != nil {
		params.Set("environment", string(*req.Environment))
	}
	if req.ProjectID != nil {
		params.Set("projectId", *req.ProjectID)
	}

	path := "/webdata/batch/" + req.ID
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp BatchExtractionJob
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListBatchJobs lists batch jobs with pagination and filters.
func (c *Client) ListBatchJobs(ctx context.Context, req *ListBatchJobsRequest) (*BatchJobsResponse, error) {
	params := url.Values{}
	if req != nil {
		if req.Environment != nil {
			params.Set("environment", string(*req.Environment))
		}
		if req.ProjectID != nil {
			params.Set("projectId", *req.ProjectID)
		}
		if req.Status != nil {
			params.Set("status", string(*req.Status))
		}
		params.Set("type", "extraction")
		if req.Limit != nil {
			params.Set("limit", strconv.Itoa(*req.Limit))
		}
		if req.Cursor != nil {
			params.Set("cursor", *req.Cursor)
		}
	}

	path := "/webdata/batch"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp BatchJobsResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelBatchJob cancels a batch job.
func (c *Client) CancelBatchJob(ctx context.Context, req *GetBatchJobRequest) (*SuccessResponse, error) {
	params := url.Values{}
	if req.Environment != nil {
		params.Set("environment", string(*req.Environment))
	}
	if req.ProjectID != nil {
		params.Set("projectId", *req.ProjectID)
	}

	path := "/webdata/batch/" + req.ID + "/cancel"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp SuccessResponse
	if err := c.http.Post(ctx, path, map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BatchAndWait creates a batch job and waits for completion.
func (c *Client) BatchAndWait(ctx context.Context, req *CreateBatchExtractionsRequest, opts *ExtractAndWaitOptions) (*BatchExtractionJob, error) {
	pollInterval := 2 * time.Second
	timeout := 300 * time.Second
	if opts != nil {
		if opts.PollInterval > 0 {
			pollInterval = opts.PollInterval
		}
		if opts.Timeout > 0 {
			timeout = opts.Timeout
		}
	}

	resp, err := c.Batch(ctx, req)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	for time.Since(startTime) < timeout {
		job, err := c.GetBatchJob(ctx, &GetBatchJobRequest{
			ID:          resp.ID,
			Environment: req.Environment,
			ProjectID:   req.ProjectID,
		})
		if err != nil {
			return nil, err
		}

		if job.Status == types.BatchJobStatusCompleted || job.Status == types.BatchJobStatusFailed || job.Status == types.BatchJobStatusCancelled {
			return job, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(pollInterval):
		}
	}

	return nil, types.NewTimeoutError("Batch job timed out")
}

// CreateSchedule creates a scheduled extraction job.
func (c *Client) CreateSchedule(ctx context.Context, req *CreateExtractionScheduleRequest) (*CreateScheduleResponse, error) {
	body := map[string]interface{}{
		"type": "extraction",
		"name": req.Name,
		"url":  req.URL,
	}
	if req.Environment != nil {
		body["environment"] = *req.Environment
	}
	if req.ProjectID != nil {
		body["projectId"] = *req.ProjectID
	}
	if req.Frequency != nil {
		body["frequency"] = *req.Frequency
	}
	if req.Config != nil {
		body["config"] = req.Config
	}
	if req.DetectChanges != nil {
		body["detectChanges"] = *req.DetectChanges
	}
	if req.ChangeThreshold != nil {
		body["changeThreshold"] = *req.ChangeThreshold
	}
	if req.WebhookURL != nil {
		body["webhookUrl"] = *req.WebhookURL
	}
	if req.WebhookSecret != nil {
		body["webhookSecret"] = *req.WebhookSecret
	}
	if req.Metadata != nil {
		body["metadata"] = req.Metadata
	}

	var resp CreateScheduleResponse
	if err := c.http.Post(ctx, "/webdata/schedules", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateSchedule updates a schedule.
func (c *Client) UpdateSchedule(ctx context.Context, req *UpdateExtractionScheduleRequest) (*SuccessResponse, error) {
	params := url.Values{}
	if req.Environment != nil {
		params.Set("environment", string(*req.Environment))
	}
	if req.ProjectID != nil {
		params.Set("projectId", *req.ProjectID)
	}

	path := "/webdata/schedules/" + req.ID
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	body := map[string]interface{}{}
	if req.Name != nil {
		body["name"] = *req.Name
	}
	if req.Frequency != nil {
		body["frequency"] = *req.Frequency
	}
	if req.Config != nil {
		body["config"] = req.Config
	}
	if req.IsActive != nil {
		body["isActive"] = *req.IsActive
	}
	if req.DetectChanges != nil {
		body["detectChanges"] = *req.DetectChanges
	}
	if req.ChangeThreshold != nil {
		body["changeThreshold"] = *req.ChangeThreshold
	}
	if req.WebhookURL != nil {
		body["webhookUrl"] = *req.WebhookURL
	}
	if req.WebhookSecret != nil {
		body["webhookSecret"] = *req.WebhookSecret
	}
	if req.Metadata != nil {
		body["metadata"] = req.Metadata
	}

	var resp SuccessResponse
	if err := c.http.Post(ctx, path, body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSchedule retrieves a schedule by ID.
func (c *Client) GetSchedule(ctx context.Context, req *GetScheduleRequest) (*ExtractionSchedule, error) {
	params := url.Values{}
	if req.Environment != nil {
		params.Set("environment", string(*req.Environment))
	}
	if req.ProjectID != nil {
		params.Set("projectId", *req.ProjectID)
	}

	path := "/webdata/schedules/" + req.ID
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ExtractionSchedule
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListSchedules lists schedules with pagination and filters.
func (c *Client) ListSchedules(ctx context.Context, req *ListSchedulesRequest) (*SchedulesResponse, error) {
	params := url.Values{}
	if req != nil {
		if req.Environment != nil {
			params.Set("environment", string(*req.Environment))
		}
		if req.ProjectID != nil {
			params.Set("projectId", *req.ProjectID)
		}
		params.Set("type", "extraction")
		if req.IsActive != nil {
			params.Set("isActive", strconv.FormatBool(*req.IsActive))
		}
		if req.Limit != nil {
			params.Set("limit", strconv.Itoa(*req.Limit))
		}
		if req.Cursor != nil {
			params.Set("cursor", *req.Cursor)
		}
	}

	path := "/webdata/schedules"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp SchedulesResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteSchedule deletes a schedule.
func (c *Client) DeleteSchedule(ctx context.Context, req *GetScheduleRequest) (*SuccessResponse, error) {
	params := url.Values{}
	if req.Environment != nil {
		params.Set("environment", string(*req.Environment))
	}
	if req.ProjectID != nil {
		params.Set("projectId", *req.ProjectID)
	}

	path := "/webdata/schedules/" + req.ID
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	body := map[string]interface{}{"id": req.ID}
	if req.Environment != nil {
		body["environment"] = *req.Environment
	}
	if req.ProjectID != nil {
		body["projectId"] = *req.ProjectID
	}

	var resp SuccessResponse
	if err := c.http.DeleteWithBody(ctx, path, body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ToggleSchedule toggles a schedule on or off.
func (c *Client) ToggleSchedule(ctx context.Context, req *GetScheduleRequest) (*ToggleResponse, error) {
	params := url.Values{}
	if req.Environment != nil {
		params.Set("environment", string(*req.Environment))
	}
	if req.ProjectID != nil {
		params.Set("projectId", *req.ProjectID)
	}

	path := "/webdata/schedules/" + req.ID + "/toggle"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ToggleResponse
	if err := c.http.Post(ctx, path, map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetUsage gets usage statistics.
func (c *Client) GetUsage(ctx context.Context, req *GetUsageRequest) (*ExtractionUsage, error) {
	params := url.Values{}
	if req != nil {
		if req.Environment != nil {
			params.Set("environment", string(*req.Environment))
		}
		if req.PeriodStart != nil {
			params.Set("periodStart", *req.PeriodStart)
		}
		if req.PeriodEnd != nil {
			params.Set("periodEnd", *req.PeriodEnd)
		}
	}

	path := "/webdata/usage"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ExtractionUsage
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetUsageDaily gets daily usage breakdown.
func (c *Client) GetUsageDaily(ctx context.Context, req *GetUsageRequest) (*GetDailyUsageResponse, error) {
	params := url.Values{}
	if req != nil {
		if req.Environment != nil {
			params.Set("environment", string(*req.Environment))
		}
		if req.PeriodStart != nil {
			params.Set("periodStart", *req.PeriodStart)
		}
		if req.PeriodEnd != nil {
			params.Set("periodEnd", *req.PeriodEnd)
		}
	}

	path := "/webdata/usage/daily"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp GetDailyUsageResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
