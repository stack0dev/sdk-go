package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment_Constants(t *testing.T) {
	assert.Equal(t, Environment("sandbox"), EnvironmentSandbox)
	assert.Equal(t, Environment("production"), EnvironmentProduction)
}

func TestBatchJobStatus_Constants(t *testing.T) {
	assert.Equal(t, BatchJobStatus("pending"), BatchJobStatusPending)
	assert.Equal(t, BatchJobStatus("processing"), BatchJobStatusProcessing)
	assert.Equal(t, BatchJobStatus("completed"), BatchJobStatusCompleted)
	assert.Equal(t, BatchJobStatus("failed"), BatchJobStatusFailed)
	assert.Equal(t, BatchJobStatus("cancelled"), BatchJobStatusCancelled)
}

func TestScheduleFrequency_Constants(t *testing.T) {
	assert.Equal(t, ScheduleFrequency("hourly"), ScheduleFrequencyHourly)
	assert.Equal(t, ScheduleFrequency("daily"), ScheduleFrequencyDaily)
	assert.Equal(t, ScheduleFrequency("weekly"), ScheduleFrequencyWeekly)
	assert.Equal(t, ScheduleFrequency("monthly"), ScheduleFrequencyMonthly)
}

func TestSuccessResponse_Fields(t *testing.T) {
	resp := SuccessResponse{Success: true}
	assert.True(t, resp.Success)

	resp = SuccessResponse{Success: false}
	assert.False(t, resp.Success)
}

func TestPaginatedRequest_Fields(t *testing.T) {
	limit := 10
	offset := 20
	req := PaginatedRequest{
		Limit:  &limit,
		Offset: &offset,
	}

	assert.Equal(t, 10, *req.Limit)
	assert.Equal(t, 20, *req.Offset)
}

func TestPaginatedResponse_Fields(t *testing.T) {
	resp := PaginatedResponse{
		Total:  100,
		Limit:  10,
		Offset: 0,
	}

	assert.Equal(t, 100, resp.Total)
	assert.Equal(t, 10, resp.Limit)
	assert.Equal(t, 0, resp.Offset)
}

func TestGetBatchJobRequest_Fields(t *testing.T) {
	env := EnvironmentProduction
	projectID := "proj-123"
	req := GetBatchJobRequest{
		ID:          "batch-123",
		Environment: &env,
		ProjectID:   &projectID,
	}

	assert.Equal(t, "batch-123", req.ID)
	assert.Equal(t, EnvironmentProduction, *req.Environment)
	assert.Equal(t, "proj-123", *req.ProjectID)
}

func TestListBatchJobsRequest_Fields(t *testing.T) {
	env := EnvironmentSandbox
	status := BatchJobStatusPending
	limit := 50
	cursor := "next-cursor"
	req := ListBatchJobsRequest{
		Environment: &env,
		Status:      &status,
		Limit:       &limit,
		Cursor:      &cursor,
	}

	assert.Equal(t, EnvironmentSandbox, *req.Environment)
	assert.Equal(t, BatchJobStatusPending, *req.Status)
	assert.Equal(t, 50, *req.Limit)
	assert.Equal(t, "next-cursor", *req.Cursor)
}

func TestCreateBatchResponse_Fields(t *testing.T) {
	resp := CreateBatchResponse{
		ID:        "batch-456",
		Status:    BatchJobStatusPending,
		TotalURLs: 10,
	}

	assert.Equal(t, "batch-456", resp.ID)
	assert.Equal(t, BatchJobStatusPending, resp.Status)
	assert.Equal(t, 10, resp.TotalURLs)
}

func TestGetScheduleRequest_Fields(t *testing.T) {
	env := EnvironmentProduction
	projectID := "proj-456"
	req := GetScheduleRequest{
		ID:          "sched-123",
		Environment: &env,
		ProjectID:   &projectID,
	}

	assert.Equal(t, "sched-123", req.ID)
	assert.Equal(t, EnvironmentProduction, *req.Environment)
	assert.Equal(t, "proj-456", *req.ProjectID)
}

func TestListSchedulesRequest_Fields(t *testing.T) {
	env := EnvironmentSandbox
	isActive := true
	limit := 25
	cursor := "schedule-cursor"
	req := ListSchedulesRequest{
		Environment: &env,
		IsActive:    &isActive,
		Limit:       &limit,
		Cursor:      &cursor,
	}

	assert.Equal(t, EnvironmentSandbox, *req.Environment)
	assert.True(t, *req.IsActive)
	assert.Equal(t, 25, *req.Limit)
	assert.Equal(t, "schedule-cursor", *req.Cursor)
}

func TestCreateScheduleResponse_Fields(t *testing.T) {
	resp := CreateScheduleResponse{
		ID: "sched-789",
	}

	assert.Equal(t, "sched-789", resp.ID)
}
