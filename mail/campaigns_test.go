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

func setupCampaignsTestClient(t *testing.T, handler http.HandlerFunc) (*CampaignsClient, *httptest.Server) {
	server := httptest.NewServer(handler)
	httpClient := client.New("test-api-key", server.URL)
	return NewCampaignsClient(httpClient), server
}

func TestCampaignsClient_List(t *testing.T) {
	t.Run("without filters", func(t *testing.T) {
		campaignsClient, server := setupCampaignsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/mail/campaigns", r.URL.Path)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListCampaignsResponse{
				Campaigns: []Campaign{
					{ID: "camp-1", Name: "Summer Sale", Status: "sent"},
					{ID: "camp-2", Name: "Newsletter", Status: "draft"},
				},
				Total: 2,
			})
		})
		defer server.Close()

		resp, err := campaignsClient.List(context.Background(), nil)

		require.NoError(t, err)
		assert.Len(t, resp.Campaigns, 2)
	})

	t.Run("with filters", func(t *testing.T) {
		campaignsClient, server := setupCampaignsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "status=draft")
			assert.Contains(t, r.URL.RawQuery, "environment=sandbox")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListCampaignsResponse{Campaigns: []Campaign{}})
		})
		defer server.Close()

		status := CampaignStatusDraft
		env := types.EnvironmentSandbox
		resp, err := campaignsClient.List(context.Background(), &ListCampaignsRequest{
			Status:      &status,
			Environment: &env,
		})

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func TestCampaignsClient_Get(t *testing.T) {
	campaignID := "camp-123"
	campaignsClient, server := setupCampaignsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mail/campaigns/"+campaignID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Campaign{
			ID:              campaignID,
			Name:            "Holiday Campaign",
			Subject:         "Happy Holidays!",
			FromEmail:       "marketing@example.com",
			Status:          "draft",
			TotalRecipients: 1000,
			CreatedAt:       time.Now(),
		})
	})
	defer server.Close()

	resp, err := campaignsClient.Get(context.Background(), campaignID)

	require.NoError(t, err)
	assert.Equal(t, campaignID, resp.ID)
	assert.Equal(t, "Holiday Campaign", resp.Name)
	assert.Equal(t, 1000, resp.TotalRecipients)
}

func TestCampaignsClient_Create(t *testing.T) {
	campaignsClient, server := setupCampaignsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/campaigns", r.URL.Path)

		var req CreateCampaignRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "New Campaign", req.Name)
		assert.Equal(t, "Check this out!", req.Subject)
		assert.Equal(t, "sender@example.com", req.FromEmail)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Campaign{
			ID:        "camp-new",
			Name:      req.Name,
			Subject:   req.Subject,
			FromEmail: req.FromEmail,
			Status:    "draft",
		})
	})
	defer server.Close()

	resp, err := campaignsClient.Create(context.Background(), &CreateCampaignRequest{
		Name:      "New Campaign",
		Subject:   "Check this out!",
		FromEmail: "sender@example.com",
	})

	require.NoError(t, err)
	assert.Equal(t, "camp-new", resp.ID)
	assert.Equal(t, "draft", resp.Status)
}

func TestCampaignsClient_Update(t *testing.T) {
	campaignID := "camp-123"
	campaignsClient, server := setupCampaignsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/mail/campaigns/"+campaignID, r.URL.Path)

		var req UpdateCampaignRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "Updated Subject", *req.Subject)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Campaign{
			ID:      campaignID,
			Subject: "Updated Subject",
		})
	})
	defer server.Close()

	subject := "Updated Subject"
	resp, err := campaignsClient.Update(context.Background(), &UpdateCampaignRequest{
		ID:      campaignID,
		Subject: &subject,
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated Subject", resp.Subject)
}

func TestCampaignsClient_Delete(t *testing.T) {
	campaignID := "camp-123"
	campaignsClient, server := setupCampaignsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/mail/campaigns/"+campaignID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DeleteCampaignResponse{Success: true})
	})
	defer server.Close()

	resp, err := campaignsClient.Delete(context.Background(), campaignID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestCampaignsClient_Send(t *testing.T) {
	t.Run("send now", func(t *testing.T) {
		campaignID := "camp-123"
		campaignsClient, server := setupCampaignsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/mail/campaigns/"+campaignID+"/send", r.URL.Path)

			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			require.NoError(t, err)
			assert.True(t, body["sendNow"].(bool))

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(SendCampaignResponse{
				Success:         true,
				SentCount:       500,
				TotalRecipients: 500,
			})
		})
		defer server.Close()

		sendNow := true
		resp, err := campaignsClient.Send(context.Background(), &SendCampaignRequest{
			ID:      campaignID,
			SendNow: &sendNow,
		})

		require.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, 500, resp.SentCount)
	})

	t.Run("scheduled send", func(t *testing.T) {
		campaignID := "camp-123"
		scheduledTime := time.Now().Add(24 * time.Hour)
		campaignsClient, server := setupCampaignsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			require.NoError(t, err)
			assert.NotEmpty(t, body["scheduledAt"])

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(SendCampaignResponse{Success: true})
		})
		defer server.Close()

		resp, err := campaignsClient.Send(context.Background(), &SendCampaignRequest{
			ID:          campaignID,
			ScheduledAt: &scheduledTime,
		})

		require.NoError(t, err)
		assert.True(t, resp.Success)
	})
}

func TestCampaignsClient_Pause(t *testing.T) {
	campaignID := "camp-123"
	campaignsClient, server := setupCampaignsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/campaigns/"+campaignID+"/pause", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PauseCampaignResponse{Success: true})
	})
	defer server.Close()

	resp, err := campaignsClient.Pause(context.Background(), campaignID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestCampaignsClient_Cancel(t *testing.T) {
	campaignID := "camp-123"
	campaignsClient, server := setupCampaignsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/campaigns/"+campaignID+"/cancel", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CancelCampaignResponse{Success: true})
	})
	defer server.Close()

	resp, err := campaignsClient.Cancel(context.Background(), campaignID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestCampaignsClient_Duplicate(t *testing.T) {
	campaignID := "camp-123"
	campaignsClient, server := setupCampaignsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/campaigns/"+campaignID+"/duplicate", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Campaign{
			ID:     "camp-456",
			Name:   "Original Campaign (Copy)",
			Status: "draft",
		})
	})
	defer server.Close()

	resp, err := campaignsClient.Duplicate(context.Background(), campaignID)

	require.NoError(t, err)
	assert.Equal(t, "camp-456", resp.ID)
	assert.Equal(t, "draft", resp.Status)
}

func TestCampaignsClient_GetStats(t *testing.T) {
	campaignID := "camp-123"
	campaignsClient, server := setupCampaignsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mail/campaigns/"+campaignID+"/stats", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CampaignStatsResponse{
			Total:        1000,
			Sent:         950,
			Delivered:    900,
			Opened:       400,
			Clicked:      100,
			Bounced:      50,
			DeliveryRate: 0.95,
			OpenRate:     0.44,
			ClickRate:    0.11,
		})
	})
	defer server.Close()

	resp, err := campaignsClient.GetStats(context.Background(), campaignID)

	require.NoError(t, err)
	assert.Equal(t, 1000, resp.Total)
	assert.Equal(t, 0.95, resp.DeliveryRate)
	assert.Equal(t, 0.44, resp.OpenRate)
}

func TestCampaignStatus_Constants(t *testing.T) {
	assert.Equal(t, CampaignStatus("draft"), CampaignStatusDraft)
	assert.Equal(t, CampaignStatus("scheduled"), CampaignStatusScheduled)
	assert.Equal(t, CampaignStatus("sending"), CampaignStatusSending)
	assert.Equal(t, CampaignStatus("sent"), CampaignStatusSent)
	assert.Equal(t, CampaignStatus("paused"), CampaignStatusPaused)
	assert.Equal(t, CampaignStatus("cancelled"), CampaignStatusCancelled)
	assert.Equal(t, CampaignStatus("failed"), CampaignStatusFailed)
}
