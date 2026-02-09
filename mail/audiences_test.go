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

func setupAudiencesTestClient(t *testing.T, handler http.HandlerFunc) (*AudiencesClient, *httptest.Server) {
	server := httptest.NewServer(handler)
	httpClient := client.New("test-api-key", server.URL)
	return NewAudiencesClient(httpClient), server
}

func TestAudiencesClient_List(t *testing.T) {
	t.Run("without filters", func(t *testing.T) {
		audiencesClient, server := setupAudiencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/mail/audiences", r.URL.Path)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListAudiencesResponse{
				Audiences: []Audience{
					{ID: "aud-1", Name: "Newsletter", TotalContacts: 1000},
					{ID: "aud-2", Name: "Beta Users", TotalContacts: 50},
				},
				Total: 2,
			})
		})
		defer server.Close()

		resp, err := audiencesClient.List(context.Background(), nil)

		require.NoError(t, err)
		assert.Len(t, resp.Audiences, 2)
		assert.Equal(t, 1000, resp.Audiences[0].TotalContacts)
	})

	t.Run("with filters", func(t *testing.T) {
		audiencesClient, server := setupAudiencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "environment=sandbox")
			assert.Contains(t, r.URL.RawQuery, "search=news")
			assert.Contains(t, r.URL.RawQuery, "limit=10")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListAudiencesResponse{Audiences: []Audience{}})
		})
		defer server.Close()

		env := types.EnvironmentSandbox
		search := "news"
		limit := 10
		resp, err := audiencesClient.List(context.Background(), &ListAudiencesRequest{
			Environment: &env,
			Search:      &search,
			Limit:       &limit,
		})

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func TestAudiencesClient_Get(t *testing.T) {
	audienceID := "aud-123"
	audiencesClient, server := setupAudiencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mail/audiences/"+audienceID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Audience{
			ID:                   audienceID,
			Name:                 "Newsletter Subscribers",
			TotalContacts:        500,
			SubscribedContacts:   480,
			UnsubscribedContacts: 20,
			CreatedAt:            time.Now(),
		})
	})
	defer server.Close()

	resp, err := audiencesClient.Get(context.Background(), audienceID)

	require.NoError(t, err)
	assert.Equal(t, audienceID, resp.ID)
	assert.Equal(t, "Newsletter Subscribers", resp.Name)
	assert.Equal(t, 500, resp.TotalContacts)
}

func TestAudiencesClient_Create(t *testing.T) {
	audiencesClient, server := setupAudiencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/audiences", r.URL.Path)

		var req CreateAudienceRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "New Audience", req.Name)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Audience{
			ID:   "aud-new",
			Name: req.Name,
		})
	})
	defer server.Close()

	resp, err := audiencesClient.Create(context.Background(), &CreateAudienceRequest{
		Name: "New Audience",
	})

	require.NoError(t, err)
	assert.Equal(t, "aud-new", resp.ID)
	assert.Equal(t, "New Audience", resp.Name)
}

func TestAudiencesClient_Update(t *testing.T) {
	audienceID := "aud-123"
	audiencesClient, server := setupAudiencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/mail/audiences/"+audienceID, r.URL.Path)

		var req UpdateAudienceRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", *req.Name)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Audience{
			ID:   audienceID,
			Name: "Updated Name",
		})
	})
	defer server.Close()

	name := "Updated Name"
	resp, err := audiencesClient.Update(context.Background(), &UpdateAudienceRequest{
		ID:   audienceID,
		Name: &name,
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated Name", resp.Name)
}

func TestAudiencesClient_Delete(t *testing.T) {
	audienceID := "aud-123"
	audiencesClient, server := setupAudiencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/mail/audiences/"+audienceID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DeleteAudienceResponse{Success: true})
	})
	defer server.Close()

	resp, err := audiencesClient.Delete(context.Background(), audienceID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestAudiencesClient_ListContacts(t *testing.T) {
	audienceID := "aud-123"
	audiencesClient, server := setupAudiencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/mail/audiences/"+audienceID+"/contacts")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListAudienceContactsResponse{
			Contacts: []AudienceContact{
				{MailContact: MailContact{ID: "contact-1", Email: "user1@example.com"}},
				{MailContact: MailContact{ID: "contact-2", Email: "user2@example.com"}},
			},
			Total: 2,
		})
	})
	defer server.Close()

	resp, err := audiencesClient.ListContacts(context.Background(), &ListAudienceContactsRequest{
		ID: audienceID,
	})

	require.NoError(t, err)
	assert.Len(t, resp.Contacts, 2)
}

func TestAudiencesClient_ListContacts_WithFilters(t *testing.T) {
	audienceID := "aud-123"
	audiencesClient, server := setupAudiencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.RawQuery, "status=subscribed")
		assert.Contains(t, r.URL.RawQuery, "limit=20")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListAudienceContactsResponse{Contacts: []AudienceContact{}})
	})
	defer server.Close()

	status := ContactStatusSubscribed
	limit := 20
	resp, err := audiencesClient.ListContacts(context.Background(), &ListAudienceContactsRequest{
		ID:     audienceID,
		Status: &status,
		Limit:  &limit,
	})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestAudiencesClient_AddContacts(t *testing.T) {
	audienceID := "aud-123"
	audiencesClient, server := setupAudiencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/audiences/"+audienceID+"/contacts", r.URL.Path)

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		contactIDs := body["contactIds"].([]interface{})
		assert.Len(t, contactIDs, 3)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AddContactsToAudienceResponse{
			Success: true,
			Added:   3,
		})
	})
	defer server.Close()

	resp, err := audiencesClient.AddContacts(context.Background(), &AddContactsToAudienceRequest{
		ID:         audienceID,
		ContactIDs: []string{"contact-1", "contact-2", "contact-3"},
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, 3, resp.Added)
}

func TestAudiencesClient_RemoveContacts(t *testing.T) {
	audienceID := "aud-123"
	audiencesClient, server := setupAudiencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/mail/audiences/"+audienceID+"/contacts", r.URL.Path)

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		contactIDs := body["contactIds"].([]interface{})
		assert.Len(t, contactIDs, 2)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(RemoveContactsFromAudienceResponse{
			Success: true,
			Removed: 2,
		})
	})
	defer server.Close()

	resp, err := audiencesClient.RemoveContacts(context.Background(), &RemoveContactsFromAudienceRequest{
		ID:         audienceID,
		ContactIDs: []string{"contact-1", "contact-2"},
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, 2, resp.Removed)
}
