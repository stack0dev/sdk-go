package extraction

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

func setupExtractionTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	httpClient := client.New("test-api-key", server.URL)
	return NewClient(httpClient), server
}

func TestClient_Extract(t *testing.T) {
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/webdata/extractions", r.URL.Path)

		var req CreateExtractionRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "https://example.com/article", req.URL)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CreateExtractionResponse{
			ID:     "ext-123",
			Status: ExtractionStatusPending,
		})
	})
	defer server.Close()

	resp, err := extractionClient.Extract(context.Background(), &CreateExtractionRequest{
		URL: "https://example.com/article",
	})

	require.NoError(t, err)
	assert.Equal(t, "ext-123", resp.ID)
	assert.Equal(t, ExtractionStatusPending, resp.Status)
}

func TestClient_Extract_WithSchema(t *testing.T) {
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var req CreateExtractionRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, ExtractionModeSchema, *req.Mode)
		assert.NotNil(t, req.Schema)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CreateExtractionResponse{
			ID:     "ext-123",
			Status: ExtractionStatusPending,
		})
	})
	defer server.Close()

	mode := ExtractionModeSchema
	resp, err := extractionClient.Extract(context.Background(), &CreateExtractionRequest{
		URL:  "https://example.com/product",
		Mode: &mode,
		Schema: map[string]interface{}{
			"title":       "string",
			"price":       "number",
			"description": "string",
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "ext-123", resp.ID)
}

func TestClient_Get(t *testing.T) {
	extractionID := "ext-123"
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/extractions/"+extractionID)

		markdown := "# Article Title\n\nThis is the content."
		tokensUsed := 150
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ExtractionResult{
			ID:         extractionID,
			URL:        "https://example.com/article",
			Status:     ExtractionStatusCompleted,
			Markdown:   &markdown,
			TokensUsed: &tokensUsed,
		})
	})
	defer server.Close()

	resp, err := extractionClient.Get(context.Background(), &GetExtractionRequest{
		ID: extractionID,
	})

	require.NoError(t, err)
	assert.Equal(t, extractionID, resp.ID)
	assert.Equal(t, ExtractionStatusCompleted, resp.Status)
	assert.NotNil(t, resp.Markdown)
	assert.Contains(t, *resp.Markdown, "Article Title")
}

func TestClient_List(t *testing.T) {
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/extractions")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListExtractionsResponse{
			Items: []ExtractionResult{
				{ID: "ext-1", URL: "https://example1.com", Status: ExtractionStatusCompleted},
				{ID: "ext-2", URL: "https://example2.com", Status: ExtractionStatusCompleted},
			},
		})
	})
	defer server.Close()

	resp, err := extractionClient.List(context.Background(), nil)

	require.NoError(t, err)
	assert.Len(t, resp.Items, 2)
}

func TestClient_Delete(t *testing.T) {
	extractionID := "ext-123"
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/extractions/"+extractionID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SuccessResponse{Success: true})
	})
	defer server.Close()

	resp, err := extractionClient.Delete(context.Background(), &GetExtractionRequest{
		ID: extractionID,
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestClient_ExtractAndWait_Success(t *testing.T) {
	var callCount int32

	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/webdata/extractions" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CreateExtractionResponse{
				ID:     "ext-123",
				Status: ExtractionStatusPending,
			})
			return
		}

		if r.Method == http.MethodGet {
			count := atomic.AddInt32(&callCount, 1)
			markdown := "# Extracted Content"

			var status ExtractionStatus
			if count < 3 {
				status = ExtractionStatusProcessing
			} else {
				status = ExtractionStatusCompleted
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ExtractionResult{
				ID:       "ext-123",
				Status:   status,
				Markdown: &markdown,
			})
		}
	})
	defer server.Close()

	resp, err := extractionClient.ExtractAndWait(context.Background(), &CreateExtractionRequest{
		URL: "https://example.com",
	}, &ExtractAndWaitOptions{
		PollInterval: 10 * time.Millisecond,
		Timeout:      5 * time.Second,
	})

	require.NoError(t, err)
	assert.Equal(t, "ext-123", resp.ID)
	assert.Equal(t, ExtractionStatusCompleted, resp.Status)
	assert.GreaterOrEqual(t, callCount, int32(3))
}

func TestClient_ExtractAndWait_Failed(t *testing.T) {
	var callCount int32
	errorMessage := "Failed to extract content"

	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CreateExtractionResponse{
				ID:     "ext-123",
				Status: ExtractionStatusPending,
			})
			return
		}

		count := atomic.AddInt32(&callCount, 1)
		var status ExtractionStatus
		var errPtr *string
		if count < 2 {
			status = ExtractionStatusProcessing
		} else {
			status = ExtractionStatusFailed
			errPtr = &errorMessage
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ExtractionResult{
			ID:     "ext-123",
			Status: status,
			Error:  errPtr,
		})
	})
	defer server.Close()

	resp, err := extractionClient.ExtractAndWait(context.Background(), &CreateExtractionRequest{
		URL: "https://example.com",
	}, &ExtractAndWaitOptions{
		PollInterval: 10 * time.Millisecond,
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Failed to extract content")
}

func TestClient_ExtractAndWait_Timeout(t *testing.T) {
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CreateExtractionResponse{
				ID:     "ext-123",
				Status: ExtractionStatusPending,
			})
			return
		}

		// Always return processing
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ExtractionResult{
			ID:     "ext-123",
			Status: ExtractionStatusProcessing,
		})
	})
	defer server.Close()

	resp, err := extractionClient.ExtractAndWait(context.Background(), &CreateExtractionRequest{
		URL: "https://example.com",
	}, &ExtractAndWaitOptions{
		PollInterval: 10 * time.Millisecond,
		Timeout:      50 * time.Millisecond,
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "timeout")
}

func TestClient_ExtractAndWait_ContextCancelled(t *testing.T) {
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CreateExtractionResponse{
				ID:     "ext-123",
				Status: ExtractionStatusPending,
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ExtractionResult{
			ID:     "ext-123",
			Status: ExtractionStatusProcessing,
		})
	})
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()

	resp, err := extractionClient.ExtractAndWait(ctx, &CreateExtractionRequest{
		URL: "https://example.com",
	}, &ExtractAndWaitOptions{
		PollInterval: 10 * time.Millisecond,
		Timeout:      5 * time.Second,
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, context.Canceled, err)
}

func TestClient_Batch(t *testing.T) {
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/webdata/batch/extractions", r.URL.Path)

		var req CreateBatchExtractionsRequest
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

	resp, err := extractionClient.Batch(context.Background(), &CreateBatchExtractionsRequest{
		URLs: []string{"https://example1.com", "https://example2.com", "https://example3.com"},
	})

	require.NoError(t, err)
	assert.Equal(t, "batch-123", resp.ID)
	assert.Equal(t, 3, resp.TotalURLs)
}

func TestClient_GetBatchJob(t *testing.T) {
	batchID := "batch-123"
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/batch/"+batchID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(BatchExtractionJob{
			ID:             batchID,
			Status:         types.BatchJobStatusCompleted,
			TotalURLs:      3,
			ProcessedURLs:  3,
			SuccessfulURLs: 2,
			FailedURLs:     1,
		})
	})
	defer server.Close()

	resp, err := extractionClient.GetBatchJob(context.Background(), &GetBatchJobRequest{
		ID: batchID,
	})

	require.NoError(t, err)
	assert.Equal(t, batchID, resp.ID)
	assert.Equal(t, types.BatchJobStatusCompleted, resp.Status)
	assert.Equal(t, 2, resp.SuccessfulURLs)
}

func TestClient_ListBatchJobs(t *testing.T) {
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/batch")
		assert.Contains(t, r.URL.RawQuery, "type=extraction")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(BatchJobsResponse{
			Items: []BatchExtractionJob{
				{ID: "batch-1", Status: types.BatchJobStatusCompleted},
				{ID: "batch-2", Status: types.BatchJobStatusProcessing},
			},
		})
	})
	defer server.Close()

	resp, err := extractionClient.ListBatchJobs(context.Background(), nil)

	require.NoError(t, err)
	assert.Len(t, resp.Items, 2)
}

func TestClient_CancelBatchJob(t *testing.T) {
	batchID := "batch-123"
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/batch/"+batchID+"/cancel")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SuccessResponse{Success: true})
	})
	defer server.Close()

	resp, err := extractionClient.CancelBatchJob(context.Background(), &GetBatchJobRequest{
		ID: batchID,
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestClient_BatchAndWait_Success(t *testing.T) {
	var callCount int32

	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/webdata/batch/extractions" {
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
			json.NewEncoder(w).Encode(BatchExtractionJob{
				ID:             "batch-123",
				Status:         status,
				TotalURLs:      3,
				ProcessedURLs:  3,
				SuccessfulURLs: 3,
			})
		}
	})
	defer server.Close()

	resp, err := extractionClient.BatchAndWait(context.Background(), &CreateBatchExtractionsRequest{
		URLs: []string{"https://example1.com", "https://example2.com", "https://example3.com"},
	}, &ExtractAndWaitOptions{
		PollInterval: 10 * time.Millisecond,
		Timeout:      5 * time.Second,
	})

	require.NoError(t, err)
	assert.Equal(t, "batch-123", resp.ID)
	assert.Equal(t, types.BatchJobStatusCompleted, resp.Status)
}

func TestClient_CreateSchedule(t *testing.T) {
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/webdata/schedules", r.URL.Path)

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "extraction", body["type"])
		assert.Equal(t, "My Extraction Schedule", body["name"])

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CreateScheduleResponse{
			ID: "sched-123",
		})
	})
	defer server.Close()

	resp, err := extractionClient.CreateSchedule(context.Background(), &CreateExtractionScheduleRequest{
		Name: "My Extraction Schedule",
		URL:  "https://example.com",
	})

	require.NoError(t, err)
	assert.Equal(t, "sched-123", resp.ID)
}

func TestClient_GetSchedule(t *testing.T) {
	scheduleID := "sched-123"
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/schedules/"+scheduleID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ExtractionSchedule{
			ID:        scheduleID,
			Name:      "My Schedule",
			URL:       "https://example.com",
			Frequency: types.ScheduleFrequencyDaily,
			IsActive:  true,
		})
	})
	defer server.Close()

	resp, err := extractionClient.GetSchedule(context.Background(), &GetScheduleRequest{
		ID: scheduleID,
	})

	require.NoError(t, err)
	assert.Equal(t, scheduleID, resp.ID)
	assert.True(t, resp.IsActive)
}

func TestClient_ListSchedules(t *testing.T) {
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/schedules")
		assert.Contains(t, r.URL.RawQuery, "type=extraction")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SchedulesResponse{
			Items: []ExtractionSchedule{
				{ID: "sched-1", Name: "Schedule 1"},
				{ID: "sched-2", Name: "Schedule 2"},
			},
		})
	})
	defer server.Close()

	resp, err := extractionClient.ListSchedules(context.Background(), nil)

	require.NoError(t, err)
	assert.Len(t, resp.Items, 2)
}

func TestClient_DeleteSchedule(t *testing.T) {
	scheduleID := "sched-123"
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/schedules/"+scheduleID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SuccessResponse{Success: true})
	})
	defer server.Close()

	resp, err := extractionClient.DeleteSchedule(context.Background(), &GetScheduleRequest{
		ID: scheduleID,
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestClient_ToggleSchedule(t *testing.T) {
	scheduleID := "sched-123"
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/schedules/"+scheduleID+"/toggle")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ToggleResponse{IsActive: false})
	})
	defer server.Close()

	resp, err := extractionClient.ToggleSchedule(context.Background(), &GetScheduleRequest{
		ID: scheduleID,
	})

	require.NoError(t, err)
	assert.False(t, resp.IsActive)
}

func TestClient_GetUsage(t *testing.T) {
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/usage")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ExtractionUsage{
			PeriodStart:           time.Now().AddDate(0, -1, 0),
			PeriodEnd:             time.Now(),
			ExtractionsTotal:      1000,
			ExtractionsSuccessful: 950,
			ExtractionsFailed:     50,
			ExtractionCreditsUsed: 1000,
			ExtractionTokensUsed:  500000,
		})
	})
	defer server.Close()

	resp, err := extractionClient.GetUsage(context.Background(), nil)

	require.NoError(t, err)
	assert.Equal(t, 1000, resp.ExtractionsTotal)
	assert.Equal(t, 950, resp.ExtractionsSuccessful)
}

func TestClient_GetUsageDaily(t *testing.T) {
	extractionClient, server := setupExtractionTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/webdata/usage/daily")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetDailyUsageResponse{
			Days: []DailyUsageItem{
				{Date: "2024-01-01", Screenshots: 50, Extractions: 100, CreditsUsed: 150},
				{Date: "2024-01-02", Screenshots: 60, Extractions: 120, CreditsUsed: 180},
			},
		})
	})
	defer server.Close()

	resp, err := extractionClient.GetUsageDaily(context.Background(), nil)

	require.NoError(t, err)
	assert.Len(t, resp.Days, 2)
}

func TestExtractionStatus_Constants(t *testing.T) {
	assert.Equal(t, ExtractionStatus("pending"), ExtractionStatusPending)
	assert.Equal(t, ExtractionStatus("processing"), ExtractionStatusProcessing)
	assert.Equal(t, ExtractionStatus("completed"), ExtractionStatusCompleted)
	assert.Equal(t, ExtractionStatus("failed"), ExtractionStatusFailed)
}

func TestExtractionMode_Constants(t *testing.T) {
	assert.Equal(t, ExtractionMode("auto"), ExtractionModeAuto)
	assert.Equal(t, ExtractionMode("schema"), ExtractionModeSchema)
	assert.Equal(t, ExtractionMode("markdown"), ExtractionModeMarkdown)
	assert.Equal(t, ExtractionMode("raw"), ExtractionModeRaw)
}
