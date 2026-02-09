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

func setupContactsTestClient(t *testing.T, handler http.HandlerFunc) (*ContactsClient, *httptest.Server) {
	server := httptest.NewServer(handler)
	httpClient := client.New("test-api-key", server.URL)
	return NewContactsClient(httpClient), server
}

func TestContactsClient_List(t *testing.T) {
	t.Run("without filters", func(t *testing.T) {
		contactsClient, server := setupContactsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/mail/contacts", r.URL.Path)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListContactsResponse{
				Contacts: []MailContact{
					{ID: "contact-1", Email: "user1@example.com", Status: "subscribed"},
					{ID: "contact-2", Email: "user2@example.com", Status: "subscribed"},
				},
				Total: 2,
			})
		})
		defer server.Close()

		resp, err := contactsClient.List(context.Background(), nil)

		require.NoError(t, err)
		assert.Len(t, resp.Contacts, 2)
		assert.Equal(t, 2, resp.Total)
	})

	t.Run("with filters", func(t *testing.T) {
		contactsClient, server := setupContactsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "environment=production")
			assert.Contains(t, r.URL.RawQuery, "status=subscribed")
			assert.Contains(t, r.URL.RawQuery, "search=john")
			assert.Contains(t, r.URL.RawQuery, "limit=25")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListContactsResponse{Contacts: []MailContact{}})
		})
		defer server.Close()

		env := types.EnvironmentProduction
		status := ContactStatusSubscribed
		search := "john"
		limit := 25
		resp, err := contactsClient.List(context.Background(), &ListContactsRequest{
			Environment: &env,
			Status:      &status,
			Search:      &search,
			Limit:       &limit,
		})

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func TestContactsClient_Get(t *testing.T) {
	contactID := "contact-123"
	contactsClient, server := setupContactsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mail/contacts/"+contactID, r.URL.Path)

		firstName := "John"
		lastName := "Doe"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MailContact{
			ID:        contactID,
			Email:     "john@example.com",
			FirstName: &firstName,
			LastName:  &lastName,
			Status:    "subscribed",
			CreatedAt: time.Now(),
		})
	})
	defer server.Close()

	resp, err := contactsClient.Get(context.Background(), contactID)

	require.NoError(t, err)
	assert.Equal(t, contactID, resp.ID)
	assert.Equal(t, "john@example.com", resp.Email)
	assert.Equal(t, "John", *resp.FirstName)
	assert.Equal(t, "Doe", *resp.LastName)
}

func TestContactsClient_Create(t *testing.T) {
	contactsClient, server := setupContactsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/contacts", r.URL.Path)

		var req CreateContactRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "newuser@example.com", req.Email)
		assert.Equal(t, "Jane", *req.FirstName)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MailContact{
			ID:        "contact-new",
			Email:     req.Email,
			FirstName: req.FirstName,
			Status:    "subscribed",
		})
	})
	defer server.Close()

	firstName := "Jane"
	resp, err := contactsClient.Create(context.Background(), &CreateContactRequest{
		Email:     "newuser@example.com",
		FirstName: &firstName,
	})

	require.NoError(t, err)
	assert.Equal(t, "contact-new", resp.ID)
	assert.Equal(t, "newuser@example.com", resp.Email)
}

func TestContactsClient_Update(t *testing.T) {
	contactID := "contact-123"
	contactsClient, server := setupContactsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/mail/contacts/"+contactID, r.URL.Path)

		var req UpdateContactRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "UpdatedFirst", *req.FirstName)

		firstName := "UpdatedFirst"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MailContact{
			ID:        contactID,
			Email:     "user@example.com",
			FirstName: &firstName,
		})
	})
	defer server.Close()

	firstName := "UpdatedFirst"
	resp, err := contactsClient.Update(context.Background(), &UpdateContactRequest{
		ID:        contactID,
		FirstName: &firstName,
	})

	require.NoError(t, err)
	assert.Equal(t, "UpdatedFirst", *resp.FirstName)
}

func TestContactsClient_Update_Status(t *testing.T) {
	contactID := "contact-123"
	contactsClient, server := setupContactsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var req UpdateContactRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, ContactStatusUnsubscribed, *req.Status)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MailContact{
			ID:     contactID,
			Status: string(ContactStatusUnsubscribed),
		})
	})
	defer server.Close()

	status := ContactStatusUnsubscribed
	resp, err := contactsClient.Update(context.Background(), &UpdateContactRequest{
		ID:     contactID,
		Status: &status,
	})

	require.NoError(t, err)
	assert.Equal(t, string(ContactStatusUnsubscribed), resp.Status)
}

func TestContactsClient_Delete(t *testing.T) {
	contactID := "contact-123"
	contactsClient, server := setupContactsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/mail/contacts/"+contactID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DeleteContactResponse{Success: true})
	})
	defer server.Close()

	resp, err := contactsClient.Delete(context.Background(), contactID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestContactsClient_Import(t *testing.T) {
	contactsClient, server := setupContactsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/contacts/import", r.URL.Path)

		var req ImportContactsRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Len(t, req.Contacts, 3)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ImportContactsResponse{
			Success:  true,
			Imported: 2,
			Skipped:  1,
			Errors: []ImportContactError{
				{Email: "invalid@", Error: "Invalid email format"},
			},
		})
	})
	defer server.Close()

	firstName1 := "John"
	firstName2 := "Jane"
	resp, err := contactsClient.Import(context.Background(), &ImportContactsRequest{
		Contacts: []ImportContactInput{
			{Email: "john@example.com", FirstName: &firstName1},
			{Email: "jane@example.com", FirstName: &firstName2},
			{Email: "invalid@"},
		},
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, 2, resp.Imported)
	assert.Equal(t, 1, resp.Skipped)
	assert.Len(t, resp.Errors, 1)
}

func TestContactsClient_Import_WithAudience(t *testing.T) {
	audienceID := "aud-123"
	contactsClient, server := setupContactsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var req ImportContactsRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, audienceID, *req.AudienceID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ImportContactsResponse{
			Success:  true,
			Imported: 1,
		})
	})
	defer server.Close()

	resp, err := contactsClient.Import(context.Background(), &ImportContactsRequest{
		AudienceID: &audienceID,
		Contacts: []ImportContactInput{
			{Email: "user@example.com"},
		},
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestContactStatus_Constants(t *testing.T) {
	assert.Equal(t, ContactStatus("subscribed"), ContactStatusSubscribed)
	assert.Equal(t, ContactStatus("unsubscribed"), ContactStatusUnsubscribed)
	assert.Equal(t, ContactStatus("bounced"), ContactStatusBounced)
	assert.Equal(t, ContactStatus("complained"), ContactStatusComplained)
}
