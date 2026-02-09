package mail

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stack0/sdk-go/client"
	"github.com/stack0/sdk-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	httpClient := client.New("test-api-key", server.URL)
	return New(httpClient), server
}

func TestClient_Send(t *testing.T) {
	t.Run("successful send", func(t *testing.T) {
		expectedID := "email-123"
		mailClient, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/mail/send", r.URL.Path)
			assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var req SendEmailRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)
			assert.Equal(t, "Hello", req.Subject)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(SendEmailResponse{
				ID:        expectedID,
				From:      "sender@example.com",
				To:        "recipient@example.com",
				Subject:   "Hello",
				Status:    "pending",
				CreatedAt: time.Now(),
			})
		})
		defer server.Close()

		resp, err := mailClient.Send(context.Background(), &SendEmailRequest{
			From:    "sender@example.com",
			To:      "recipient@example.com",
			Subject: "Hello",
			HTML:    ptr("<p>World</p>"),
		})

		require.NoError(t, err)
		assert.Equal(t, expectedID, resp.ID)
		assert.Equal(t, "pending", resp.Status)
	})

	t.Run("error response", func(t *testing.T) {
		mailClient, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(types.ErrorResponse{
				Message: "Invalid email address",
				Code:    "INVALID_EMAIL",
			})
		})
		defer server.Close()

		resp, err := mailClient.Send(context.Background(), &SendEmailRequest{
			From:    "invalid",
			To:      "recipient@example.com",
			Subject: "Hello",
		})

		require.Error(t, err)
		assert.Nil(t, resp)
		apiErr, ok := err.(*types.APIError)
		require.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
		assert.Equal(t, "INVALID_EMAIL", apiErr.Code)
	})
}

func TestClient_SendBatch(t *testing.T) {
	mailClient, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/send/batch", r.URL.Path)

		var req SendBatchEmailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Len(t, req.Emails, 2)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SendBatchEmailResponse{
			Success: true,
			Data: []BatchEmailResult{
				{ID: "email-1", Success: true},
				{ID: "email-2", Success: true},
			},
		})
	})
	defer server.Close()

	resp, err := mailClient.SendBatch(context.Background(), &SendBatchEmailRequest{
		Emails: []SendEmailRequest{
			{From: "sender@example.com", To: "user1@example.com", Subject: "Hi 1"},
			{From: "sender@example.com", To: "user2@example.com", Subject: "Hi 2"},
		},
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Len(t, resp.Data, 2)
}

func TestClient_SendBroadcast(t *testing.T) {
	mailClient, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/send/broadcast", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SendBroadcastEmailResponse{
			Success: true,
			Count:   100,
		})
	})
	defer server.Close()

	resp, err := mailClient.SendBroadcast(context.Background(), &SendBroadcastEmailRequest{
		From:    "sender@example.com",
		To:      []interface{}{"user1@example.com", "user2@example.com"},
		Subject: "Broadcast",
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, 100, resp.Count)
}

func TestClient_Get(t *testing.T) {
	emailID := "email-123"
	mailClient, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mail/"+emailID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetEmailResponse{
			ID:      emailID,
			From:    "sender@example.com",
			To:      "recipient@example.com",
			Subject: "Hello",
			Status:  "delivered",
		})
	})
	defer server.Close()

	resp, err := mailClient.Get(context.Background(), emailID)

	require.NoError(t, err)
	assert.Equal(t, emailID, resp.ID)
	assert.Equal(t, "delivered", resp.Status)
}

func TestClient_List(t *testing.T) {
	t.Run("without filters", func(t *testing.T) {
		mailClient, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/mail", r.URL.Path)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListEmailsResponse{
				Emails: []Email{{ID: "email-1"}, {ID: "email-2"}},
				Total:  2,
				Limit:  10,
				Offset: 0,
			})
		})
		defer server.Close()

		resp, err := mailClient.List(context.Background(), nil)

		require.NoError(t, err)
		assert.Len(t, resp.Emails, 2)
		assert.Equal(t, 2, resp.Total)
	})

	t.Run("with filters", func(t *testing.T) {
		mailClient, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Contains(t, r.URL.String(), "limit=5")
			assert.Contains(t, r.URL.String(), "status=sent")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListEmailsResponse{
				Emails: []Email{{ID: "email-1"}},
				Total:  1,
			})
		})
		defer server.Close()

		limit := 5
		status := EmailStatusSent
		resp, err := mailClient.List(context.Background(), &ListEmailsRequest{
			Limit:  &limit,
			Status: &status,
		})

		require.NoError(t, err)
		assert.Len(t, resp.Emails, 1)
	})
}

func TestClient_Resend(t *testing.T) {
	emailID := "email-123"
	mailClient, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/"+emailID+"/resend", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ResendEmailResponse{
			Success: true,
			Data: struct {
				ID      string `json:"id"`
				Success bool   `json:"success"`
				Error   string `json:"error,omitempty"`
			}{
				ID:      "email-456",
				Success: true,
			},
		})
	})
	defer server.Close()

	resp, err := mailClient.Resend(context.Background(), emailID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "email-456", resp.Data.ID)
}

func TestClient_Cancel(t *testing.T) {
	emailID := "email-123"
	mailClient, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/"+emailID+"/cancel", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CancelEmailResponse{Success: true})
	})
	defer server.Close()

	resp, err := mailClient.Cancel(context.Background(), emailID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestClient_GetAnalytics(t *testing.T) {
	mailClient, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mail/analytics", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(EmailAnalyticsResponse{
			Total:        1000,
			Sent:         950,
			Delivered:    900,
			Bounced:      50,
			DeliveryRate: 0.95,
			OpenRate:     0.25,
			ClickRate:    0.10,
		})
	})
	defer server.Close()

	resp, err := mailClient.GetAnalytics(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 1000, resp.Total)
	assert.Equal(t, 0.95, resp.DeliveryRate)
}

func TestClient_GetTimeSeriesAnalytics(t *testing.T) {
	t.Run("without days parameter", func(t *testing.T) {
		mailClient, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/mail/analytics/timeseries", r.URL.Path)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(TimeSeriesAnalyticsResponse{
				Data: []TimeSeriesDataPoint{{Date: "2024-01-01", Sent: 100}},
			})
		})
		defer server.Close()

		resp, err := mailClient.GetTimeSeriesAnalytics(context.Background(), nil)

		require.NoError(t, err)
		assert.Len(t, resp.Data, 1)
	})

	t.Run("with days parameter", func(t *testing.T) {
		mailClient, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.String(), "days=7")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(TimeSeriesAnalyticsResponse{
				Data: []TimeSeriesDataPoint{},
			})
		})
		defer server.Close()

		days := 7
		resp, err := mailClient.GetTimeSeriesAnalytics(context.Background(), &days)

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func TestClient_GetHourlyAnalytics(t *testing.T) {
	mailClient, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mail/analytics/hourly", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HourlyAnalyticsResponse{
			Data: []HourlyAnalyticsDataPoint{{Hour: 9, Sent: 50}},
		})
	})
	defer server.Close()

	resp, err := mailClient.GetHourlyAnalytics(context.Background())

	require.NoError(t, err)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, 9, resp.Data[0].Hour)
}

func TestClient_ListSenders(t *testing.T) {
	mailClient, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/mail/senders")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListSendersResponse{
			Senders: []Sender{{From: "sender@example.com", Total: 100}},
		})
	})
	defer server.Close()

	resp, err := mailClient.ListSenders(context.Background(), nil)

	require.NoError(t, err)
	assert.Len(t, resp.Senders, 1)
}

func TestNewClient(t *testing.T) {
	httpClient := client.New("test-key", "https://api.example.com")
	mailClient := New(httpClient)

	assert.NotNil(t, mailClient)
	assert.NotNil(t, mailClient.Domains)
	assert.NotNil(t, mailClient.Templates)
	assert.NotNil(t, mailClient.Audiences)
	assert.NotNil(t, mailClient.Contacts)
	assert.NotNil(t, mailClient.Campaigns)
	assert.NotNil(t, mailClient.Sequences)
	assert.NotNil(t, mailClient.Events)
}

// Helper function
func ptr[T any](v T) *T {
	return &v
}
