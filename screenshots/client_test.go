package screenshots

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stack0/sdk-go/client"
	"github.com/stack0/sdk-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupScreenshotsTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	httpClient := client.New("test-api-key", server.URL)
	return NewClient(httpClient), server
}

func TestClient_Capture(t *testing.T) {
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/webdata/screenshots", r.URL.Path)

		var req CreateScreenshotRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "https://example.com", req.URL)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CreateScreenshotResponse{
			ID:     "ss-123",
			Status: ScreenshotStatusPending,
		})
	})
	defer server.Close()

	resp, err := screenshotsClient.Capture(context.Background(), &CreateScreenshotRequest{
		URL: "https://example.com",
	})

	require.NoError(t, err)
	assert.Equal(t, "ss-123", resp.ID)
	assert.Equal(t, ScreenshotStatusPending, resp.Status)
}

func TestClient_Capture_WithOptions(t *testing.T) {
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var req CreateScreenshotRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, ScreenshotFormatPNG, *req.Format)
		assert.Equal(t, DeviceTypeMobile, *req.DeviceType)
		assert.True(t, *req.FullPage)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CreateScreenshotResponse{
			ID:     "ss-123",
			Status: ScreenshotStatusPending,
		})
	})
	defer server.Close()

	format := ScreenshotFormatPNG
	deviceType := DeviceTypeMobile
	fullPage := true
	resp, err := screenshotsClient.Capture(context.Background(), &CreateScreenshotRequest{
		URL:        "https://example.com",
		Format:     &format,
		DeviceType: &deviceType,
		FullPage:   &fullPage,
	})

	require.NoError(t, err)
	assert.Equal(t, "ss-123", resp.ID)
}

func TestClient_Get(t *testing.T) {
	screenshotID := "ss-123"
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/screenshots/"+screenshotID)

		imageURL := "https://cdn.example.com/screenshot.png"
		imageSize := int64(50000)
		imageWidth := 1920
		imageHeight := 1080
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Screenshot{
			ID:          screenshotID,
			URL:         "https://example.com",
			Status:      ScreenshotStatusCompleted,
			ImageURL:    &imageURL,
			ImageSize:   &imageSize,
			ImageWidth:  &imageWidth,
			ImageHeight: &imageHeight,
		})
	})
	defer server.Close()

	resp, err := screenshotsClient.Get(context.Background(), &GetScreenshotRequest{
		ID: screenshotID,
	})

	require.NoError(t, err)
	assert.Equal(t, screenshotID, resp.ID)
	assert.Equal(t, ScreenshotStatusCompleted, resp.Status)
	assert.NotNil(t, resp.ImageURL)
}

func TestClient_List(t *testing.T) {
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/screenshots")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListScreenshotsResponse{
			Items: []Screenshot{
				{ID: "ss-1", URL: "https://example1.com", Status: ScreenshotStatusCompleted},
				{ID: "ss-2", URL: "https://example2.com", Status: ScreenshotStatusCompleted},
			},
		})
	})
	defer server.Close()

	resp, err := screenshotsClient.List(context.Background(), nil)

	require.NoError(t, err)
	assert.Len(t, resp.Items, 2)
}

func TestClient_Delete(t *testing.T) {
	screenshotID := "ss-123"
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/screenshots/"+screenshotID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SuccessResponse{Success: true})
	})
	defer server.Close()

	resp, err := screenshotsClient.Delete(context.Background(), &GetScreenshotRequest{
		ID: screenshotID,
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestClient_CaptureAndWait_Success(t *testing.T) {
	var callCount int32

	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/webdata/screenshots" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CreateScreenshotResponse{
				ID:     "ss-123",
				Status: ScreenshotStatusPending,
			})
			return
		}

		if r.Method == http.MethodGet {
			count := atomic.AddInt32(&callCount, 1)
			imageURL := "https://cdn.example.com/screenshot.png"

			var status ScreenshotStatus
			if count < 3 {
				status = ScreenshotStatusProcessing
			} else {
				status = ScreenshotStatusCompleted
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Screenshot{
				ID:       "ss-123",
				Status:   status,
				ImageURL: &imageURL,
			})
		}
	})
	defer server.Close()

	resp, err := screenshotsClient.CaptureAndWait(context.Background(), &CreateScreenshotRequest{
		URL: "https://example.com",
	}, &CaptureAndWaitOptions{
		PollInterval: 10 * time.Millisecond,
		Timeout:      5 * time.Second,
	})

	require.NoError(t, err)
	assert.Equal(t, "ss-123", resp.ID)
	assert.Equal(t, ScreenshotStatusCompleted, resp.Status)
	assert.GreaterOrEqual(t, callCount, int32(3))
}

func TestClient_CaptureAndWait_Failed(t *testing.T) {
	var callCount int32
	errorMessage := "Page load failed"

	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CreateScreenshotResponse{
				ID:     "ss-123",
				Status: ScreenshotStatusPending,
			})
			return
		}

		count := atomic.AddInt32(&callCount, 1)
		var status ScreenshotStatus
		var errPtr *string
		if count < 2 {
			status = ScreenshotStatusProcessing
		} else {
			status = ScreenshotStatusFailed
			errPtr = &errorMessage
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Screenshot{
			ID:     "ss-123",
			Status: status,
			Error:  errPtr,
		})
	})
	defer server.Close()

	resp, err := screenshotsClient.CaptureAndWait(context.Background(), &CreateScreenshotRequest{
		URL: "https://example.com",
	}, &CaptureAndWaitOptions{
		PollInterval: 10 * time.Millisecond,
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Page load failed")
}

func TestClient_CaptureAndWait_Timeout(t *testing.T) {
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CreateScreenshotResponse{
				ID:     "ss-123",
				Status: ScreenshotStatusPending,
			})
			return
		}

		// Always return processing
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Screenshot{
			ID:     "ss-123",
			Status: ScreenshotStatusProcessing,
		})
	})
	defer server.Close()

	resp, err := screenshotsClient.CaptureAndWait(context.Background(), &CreateScreenshotRequest{
		URL: "https://example.com",
	}, &CaptureAndWaitOptions{
		PollInterval: 10 * time.Millisecond,
		Timeout:      50 * time.Millisecond,
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "timeout")
}

func TestClient_CaptureAndWait_ContextCancelled(t *testing.T) {
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CreateScreenshotResponse{
				ID:     "ss-123",
				Status: ScreenshotStatusPending,
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Screenshot{
			ID:     "ss-123",
			Status: ScreenshotStatusProcessing,
		})
	})
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()

	resp, err := screenshotsClient.CaptureAndWait(ctx, &CreateScreenshotRequest{
		URL: "https://example.com",
	}, &CaptureAndWaitOptions{
		PollInterval: 10 * time.Millisecond,
		Timeout:      5 * time.Second,
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, context.Canceled, err)
}

func TestClient_Batch(t *testing.T) {
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/webdata/batch/screenshots", r.URL.Path)

		var req CreateBatchScreenshotsRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Len(t, req.URLs, 3)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CreateBatchResponse{
			ID:        "batch-123",
			TotalURLs: 3,
		})
	})
	defer server.Close()

	resp, err := screenshotsClient.Batch(context.Background(), &CreateBatchScreenshotsRequest{
		URLs: []string{"https://example1.com", "https://example2.com", "https://example3.com"},
	})

	require.NoError(t, err)
	assert.Equal(t, "batch-123", resp.ID)
	assert.Equal(t, 3, resp.TotalURLs)
}

func TestClient_GetBatchJob(t *testing.T) {
	batchID := "batch-123"
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/batch/"+batchID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(BatchScreenshotJob{
			ID:             batchID,
			Status:         types.BatchJobStatusCompleted,
			TotalURLs:      3,
			ProcessedURLs:  3,
			SuccessfulURLs: 2,
			FailedURLs:     1,
		})
	})
	defer server.Close()

	resp, err := screenshotsClient.GetBatchJob(context.Background(), &GetBatchJobRequest{
		ID: batchID,
	})

	require.NoError(t, err)
	assert.Equal(t, batchID, resp.ID)
	assert.Equal(t, types.BatchJobStatusCompleted, resp.Status)
	assert.Equal(t, 2, resp.SuccessfulURLs)
}

func TestClient_ListBatchJobs(t *testing.T) {
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/batch")
		assert.Contains(t, r.URL.RawQuery, "type=screenshot")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(BatchJobsResponse{
			Items: []BatchScreenshotJob{
				{ID: "batch-1", Status: types.BatchJobStatusCompleted},
				{ID: "batch-2", Status: types.BatchJobStatusProcessing},
			},
		})
	})
	defer server.Close()

	resp, err := screenshotsClient.ListBatchJobs(context.Background(), nil)

	require.NoError(t, err)
	assert.Len(t, resp.Items, 2)
}

func TestClient_CancelBatchJob(t *testing.T) {
	batchID := "batch-123"
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/batch/"+batchID+"/cancel")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SuccessResponse{Success: true})
	})
	defer server.Close()

	resp, err := screenshotsClient.CancelBatchJob(context.Background(), &GetBatchJobRequest{
		ID: batchID,
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestClient_BatchAndWait_Success(t *testing.T) {
	var callCount int32

	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/webdata/batch/screenshots" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CreateBatchResponse{
				ID:        "batch-123",
				TotalURLs: 3,
			})
			return
		}

		if r.Method == http.MethodGet {
			count := atomic.AddInt32(&callCount, 1)

			var status types.BatchJobStatus
			if count < 2 {
				status = types.BatchJobStatusProcessing
			} else {
				status = types.BatchJobStatusCompleted
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(BatchScreenshotJob{
				ID:             "batch-123",
				Status:         status,
				TotalURLs:      3,
				ProcessedURLs:  3,
				SuccessfulURLs: 3,
			})
		}
	})
	defer server.Close()

	resp, err := screenshotsClient.BatchAndWait(context.Background(), &CreateBatchScreenshotsRequest{
		URLs: []string{"https://example1.com", "https://example2.com", "https://example3.com"},
	}, &CaptureAndWaitOptions{
		PollInterval: 10 * time.Millisecond,
		Timeout:      5 * time.Second,
	})

	require.NoError(t, err)
	assert.Equal(t, "batch-123", resp.ID)
	assert.Equal(t, types.BatchJobStatusCompleted, resp.Status)
}

func TestClient_CreateSchedule(t *testing.T) {
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/webdata/schedules", r.URL.Path)

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "screenshot", body["type"])
		assert.Equal(t, "My Schedule", body["name"])

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CreateScheduleResponse{
			ID: "sched-123",
		})
	})
	defer server.Close()

	resp, err := screenshotsClient.CreateSchedule(context.Background(), &CreateScreenshotScheduleRequest{
		Name: "My Schedule",
		URL:  "https://example.com",
	})

	require.NoError(t, err)
	assert.Equal(t, "sched-123", resp.ID)
}

func TestClient_GetSchedule(t *testing.T) {
	scheduleID := "sched-123"
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/schedules/"+scheduleID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ScreenshotSchedule{
			ID:        scheduleID,
			Name:      "My Schedule",
			URL:       "https://example.com",
			Frequency: types.ScheduleFrequencyDaily,
			IsActive:  true,
		})
	})
	defer server.Close()

	resp, err := screenshotsClient.GetSchedule(context.Background(), &GetScheduleRequest{
		ID: scheduleID,
	})

	require.NoError(t, err)
	assert.Equal(t, scheduleID, resp.ID)
	assert.True(t, resp.IsActive)
}

func TestClient_ListSchedules(t *testing.T) {
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/schedules")
		assert.Contains(t, r.URL.RawQuery, "type=screenshot")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SchedulesResponse{
			Items: []ScreenshotSchedule{
				{ID: "sched-1", Name: "Schedule 1"},
				{ID: "sched-2", Name: "Schedule 2"},
			},
		})
	})
	defer server.Close()

	resp, err := screenshotsClient.ListSchedules(context.Background(), nil)

	require.NoError(t, err)
	assert.Len(t, resp.Items, 2)
}

func TestClient_DeleteSchedule(t *testing.T) {
	scheduleID := "sched-123"
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/schedules/"+scheduleID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SuccessResponse{Success: true})
	})
	defer server.Close()

	resp, err := screenshotsClient.DeleteSchedule(context.Background(), &GetScheduleRequest{
		ID: scheduleID,
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestClient_ToggleSchedule(t *testing.T) {
	scheduleID := "sched-123"
	screenshotsClient, server := setupScreenshotsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/schedules/"+scheduleID+"/toggle")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ToggleResponse{IsActive: false})
	})
	defer server.Close()

	resp, err := screenshotsClient.ToggleSchedule(context.Background(), &GetScheduleRequest{
		ID: scheduleID,
	})

	require.NoError(t, err)
	assert.False(t, resp.IsActive)
}

func TestScreenshotStatus_Constants(t *testing.T) {
	assert.Equal(t, ScreenshotStatus("pending"), ScreenshotStatusPending)
	assert.Equal(t, ScreenshotStatus("processing"), ScreenshotStatusProcessing)
	assert.Equal(t, ScreenshotStatus("completed"), ScreenshotStatusCompleted)
	assert.Equal(t, ScreenshotStatus("failed"), ScreenshotStatusFailed)
}

func TestScreenshotFormat_Constants(t *testing.T) {
	assert.Equal(t, ScreenshotFormat("png"), ScreenshotFormatPNG)
	assert.Equal(t, ScreenshotFormat("jpeg"), ScreenshotFormatJPEG)
	assert.Equal(t, ScreenshotFormat("webp"), ScreenshotFormatWebP)
	assert.Equal(t, ScreenshotFormat("pdf"), ScreenshotFormatPDF)
}

func TestDeviceType_Constants(t *testing.T) {
	assert.Equal(t, DeviceType("desktop"), DeviceTypeDesktop)
	assert.Equal(t, DeviceType("tablet"), DeviceTypeTablet)
	assert.Equal(t, DeviceType("mobile"), DeviceTypeMobile)
}
