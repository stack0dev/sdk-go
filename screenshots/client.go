package screenshots

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/stack0/sdk-go/client"
	"github.com/stack0/sdk-go/types"
)

// Client handles screenshot operations.
type Client struct {
	http *client.HTTPClient
}

// NewClient creates a new screenshots client.
func NewClient(http *client.HTTPClient) *Client {
	return &Client{http: http}
}

// Capture captures a screenshot of a URL.
func (c *Client) Capture(ctx context.Context, req *CreateScreenshotRequest) (*CreateScreenshotResponse, error) {
	var resp CreateScreenshotResponse
	if err := c.http.Post(ctx, "/webdata/screenshots", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a screenshot by ID.
func (c *Client) Get(ctx context.Context, req *GetScreenshotRequest) (*Screenshot, error) {
	params := url.Values{}
	if req.Environment != nil {
		params.Set("environment", string(*req.Environment))
	}
	if req.ProjectID != nil {
		params.Set("projectId", *req.ProjectID)
	}

	path := "/webdata/screenshots/" + req.ID
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp Screenshot
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// List lists screenshots with pagination and filters.
func (c *Client) List(ctx context.Context, req *ListScreenshotsRequest) (*ListScreenshotsResponse, error) {
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

	path := "/webdata/screenshots"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListScreenshotsResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes a screenshot.
func (c *Client) Delete(ctx context.Context, req *GetScreenshotRequest) (*SuccessResponse, error) {
	params := url.Values{}
	if req.Environment != nil {
		params.Set("environment", string(*req.Environment))
	}
	if req.ProjectID != nil {
		params.Set("projectId", *req.ProjectID)
	}

	path := "/webdata/screenshots/" + req.ID
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

// CaptureAndWaitOptions are options for CaptureAndWait.
type CaptureAndWaitOptions struct {
	PollInterval time.Duration
	Timeout      time.Duration
}

// CaptureAndWait captures a screenshot and waits for completion.
func (c *Client) CaptureAndWait(ctx context.Context, req *CreateScreenshotRequest, opts *CaptureAndWaitOptions) (*Screenshot, error) {
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

	resp, err := c.Capture(ctx, req)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	for time.Since(startTime) < timeout {
		screenshot, err := c.Get(ctx, &GetScreenshotRequest{
			ID:          resp.ID,
			Environment: req.Environment,
			ProjectID:   req.ProjectID,
		})
		if err != nil {
			return nil, err
		}

		if screenshot.Status == ScreenshotStatusCompleted || screenshot.Status == ScreenshotStatusFailed {
			if screenshot.Status == ScreenshotStatusFailed {
				errMsg := "Screenshot failed"
				if screenshot.Error != nil {
					errMsg = *screenshot.Error
				}
				return nil, errors.New(errMsg)
			}
			return screenshot, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(pollInterval):
		}
	}

	return nil, types.NewTimeoutError("Screenshot timed out")
}

// Batch creates a batch screenshot job for multiple URLs.
func (c *Client) Batch(ctx context.Context, req *CreateBatchScreenshotsRequest) (*CreateBatchResponse, error) {
	var resp CreateBatchResponse
	if err := c.http.Post(ctx, "/webdata/batch/screenshots", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBatchJob retrieves a batch job by ID.
func (c *Client) GetBatchJob(ctx context.Context, req *GetBatchJobRequest) (*BatchScreenshotJob, error) {
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

	var resp BatchScreenshotJob
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
		params.Set("type", "screenshot")
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
func (c *Client) BatchAndWait(ctx context.Context, req *CreateBatchScreenshotsRequest, opts *CaptureAndWaitOptions) (*BatchScreenshotJob, error) {
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

// CreateSchedule creates a scheduled screenshot job.
func (c *Client) CreateSchedule(ctx context.Context, req *CreateScreenshotScheduleRequest) (*CreateScheduleResponse, error) {
	body := map[string]interface{}{
		"type": "screenshot",
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
func (c *Client) UpdateSchedule(ctx context.Context, req *UpdateScreenshotScheduleRequest) (*SuccessResponse, error) {
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
func (c *Client) GetSchedule(ctx context.Context, req *GetScheduleRequest) (*ScreenshotSchedule, error) {
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

	var resp ScreenshotSchedule
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
		params.Set("type", "screenshot")
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
