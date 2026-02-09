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

func setupEventsTestClient(t *testing.T, handler http.HandlerFunc) (*EventsClient, *httptest.Server) {
	server := httptest.NewServer(handler)
	httpClient := client.New("test-api-key", server.URL)
	return NewEventsClient(httpClient), server
}

func TestEventsClient_List(t *testing.T) {
	t.Run("without filters", func(t *testing.T) {
		eventsClient, server := setupEventsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/mail/events", r.URL.Path)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListEventsResponse{
				Events: []MailEvent{
					{ID: "evt-1", Name: "user_signup", TotalReceived: 100},
					{ID: "evt-2", Name: "purchase_completed", TotalReceived: 50},
				},
				Total: 2,
			})
		})
		defer server.Close()

		resp, err := eventsClient.List(context.Background(), nil)

		require.NoError(t, err)
		assert.Len(t, resp.Events, 2)
	})

	t.Run("with filters", func(t *testing.T) {
		eventsClient, server := setupEventsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "environment=production")
			assert.Contains(t, r.URL.RawQuery, "search=signup")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListEventsResponse{Events: []MailEvent{}})
		})
		defer server.Close()

		env := types.EnvironmentProduction
		search := "signup"
		resp, err := eventsClient.List(context.Background(), &ListEventsRequest{
			Environment: &env,
			Search:      &search,
		})

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func TestEventsClient_Get(t *testing.T) {
	eventID := "evt-123"
	eventsClient, server := setupEventsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mail/events/"+eventID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MailEvent{
			ID:            eventID,
			Name:          "user_signup",
			TotalReceived: 500,
			CreatedAt:     time.Now(),
		})
	})
	defer server.Close()

	resp, err := eventsClient.Get(context.Background(), eventID)

	require.NoError(t, err)
	assert.Equal(t, eventID, resp.ID)
	assert.Equal(t, "user_signup", resp.Name)
	assert.Equal(t, 500, resp.TotalReceived)
}

func TestEventsClient_Create(t *testing.T) {
	eventsClient, server := setupEventsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/events", r.URL.Path)

		var req CreateEventRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "new_event", req.Name)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MailEvent{
			ID:   "evt-new",
			Name: req.Name,
		})
	})
	defer server.Close()

	resp, err := eventsClient.Create(context.Background(), &CreateEventRequest{
		Name: "new_event",
	})

	require.NoError(t, err)
	assert.Equal(t, "evt-new", resp.ID)
}

func TestEventsClient_Create_WithSchema(t *testing.T) {
	eventsClient, server := setupEventsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var req CreateEventRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.NotNil(t, req.PropertiesSchema)
		assert.Len(t, req.PropertiesSchema.Properties, 2)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MailEvent{
			ID:               "evt-new",
			Name:             req.Name,
			PropertiesSchema: req.PropertiesSchema,
		})
	})
	defer server.Close()

	resp, err := eventsClient.Create(context.Background(), &CreateEventRequest{
		Name: "purchase_completed",
		PropertiesSchema: &EventPropertiesSchema{
			Properties: []EventProperty{
				{Name: "amount", Type: "number"},
				{Name: "currency", Type: "string"},
			},
		},
	})

	require.NoError(t, err)
	assert.NotNil(t, resp.PropertiesSchema)
}

func TestEventsClient_Update(t *testing.T) {
	eventID := "evt-123"
	eventsClient, server := setupEventsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/mail/events/"+eventID, r.URL.Path)

		var req UpdateEventRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "updated_event", *req.Name)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MailEvent{
			ID:   eventID,
			Name: "updated_event",
		})
	})
	defer server.Close()

	name := "updated_event"
	resp, err := eventsClient.Update(context.Background(), &UpdateEventRequest{
		ID:   eventID,
		Name: &name,
	})

	require.NoError(t, err)
	assert.Equal(t, "updated_event", resp.Name)
}

func TestEventsClient_Delete(t *testing.T) {
	eventID := "evt-123"
	eventsClient, server := setupEventsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/mail/events/"+eventID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DeleteEventResponse{Success: true})
	})
	defer server.Close()

	resp, err := eventsClient.Delete(context.Background(), eventID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestEventsClient_Track(t *testing.T) {
	eventsClient, server := setupEventsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/events/track", r.URL.Path)

		var req TrackEventRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "user_signup", req.EventName)
		assert.Equal(t, "user@example.com", *req.ContactEmail)

		occurrenceID := "occ-123"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(TrackEventResponse{
			Success:           true,
			EventOccurrenceID: &occurrenceID,
		})
	})
	defer server.Close()

	contactEmail := "user@example.com"
	resp, err := eventsClient.Track(context.Background(), &TrackEventRequest{
		EventName:    "user_signup",
		ContactEmail: &contactEmail,
		Properties: map[string]interface{}{
			"source": "website",
		},
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.EventOccurrenceID)
}

func TestEventsClient_Track_WithContactID(t *testing.T) {
	eventsClient, server := setupEventsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var req TrackEventRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "contact-123", *req.ContactID)
		assert.Nil(t, req.ContactEmail)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(TrackEventResponse{Success: true})
	})
	defer server.Close()

	contactID := "contact-123"
	resp, err := eventsClient.Track(context.Background(), &TrackEventRequest{
		EventName: "page_view",
		ContactID: &contactID,
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestEventsClient_TrackBatch(t *testing.T) {
	eventsClient, server := setupEventsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/events/track/batch", r.URL.Path)

		var req BatchTrackEventsRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Len(t, req.Events, 3)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(BatchTrackEventsResponse{
			Success: true,
			Results: []BatchTrackEventResult{
				{Success: true},
				{Success: true},
				{Success: false, Error: ptr("Invalid contact")},
			},
			TotalProcessed: 3,
			TotalFailed:    1,
		})
	})
	defer server.Close()

	email1 := "user1@example.com"
	email2 := "user2@example.com"
	email3 := "invalid"
	resp, err := eventsClient.TrackBatch(context.Background(), &BatchTrackEventsRequest{
		Events: []BatchTrackEventInput{
			{EventName: "signup", ContactEmail: &email1},
			{EventName: "signup", ContactEmail: &email2},
			{EventName: "signup", ContactEmail: &email3},
		},
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, 3, resp.TotalProcessed)
	assert.Equal(t, 1, resp.TotalFailed)
}

func TestEventsClient_ListOccurrences(t *testing.T) {
	eventsClient, server := setupEventsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/mail/events/occurrences")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListEventOccurrencesResponse{
			Occurrences: []EventOccurrence{
				{ID: "occ-1", EventID: "evt-1", ContactID: "contact-1"},
				{ID: "occ-2", EventID: "evt-1", ContactID: "contact-2"},
			},
			Total: 2,
		})
	})
	defer server.Close()

	resp, err := eventsClient.ListOccurrences(context.Background(), nil)

	require.NoError(t, err)
	assert.Len(t, resp.Occurrences, 2)
}

func TestEventsClient_ListOccurrences_WithFilters(t *testing.T) {
	eventsClient, server := setupEventsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.RawQuery, "eventId=evt-123")
		assert.Contains(t, r.URL.RawQuery, "contactId=contact-456")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListEventOccurrencesResponse{Occurrences: []EventOccurrence{}})
	})
	defer server.Close()

	eventID := "evt-123"
	contactID := "contact-456"
	resp, err := eventsClient.ListOccurrences(context.Background(), &ListEventOccurrencesRequest{
		EventID:   &eventID,
		ContactID: &contactID,
	})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestEventsClient_GetAnalytics(t *testing.T) {
	eventID := "evt-123"
	eventsClient, server := setupEventsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mail/events/analytics/"+eventID, r.URL.Path)

		receivedAt := time.Now()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(EventAnalyticsResponse{
			TotalReceived:  1000,
			LastReceivedAt: &receivedAt,
			UniqueContacts: 500,
			DailyCounts: []struct {
				Date  string `json:"date"`
				Count int    `json:"count"`
			}{
				{Date: "2024-01-01", Count: 50},
				{Date: "2024-01-02", Count: 75},
			},
		})
	})
	defer server.Close()

	resp, err := eventsClient.GetAnalytics(context.Background(), eventID)

	require.NoError(t, err)
	assert.Equal(t, 1000, resp.TotalReceived)
	assert.Equal(t, 500, resp.UniqueContacts)
	assert.Len(t, resp.DailyCounts, 2)
}
