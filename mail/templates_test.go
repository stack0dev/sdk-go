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

func setupTemplatesTestClient(t *testing.T, handler http.HandlerFunc) (*TemplatesClient, *httptest.Server) {
	server := httptest.NewServer(handler)
	httpClient := client.New("test-api-key", server.URL)
	return NewTemplatesClient(httpClient), server
}

func TestTemplatesClient_List(t *testing.T) {
	t.Run("without filters", func(t *testing.T) {
		templatesClient, server := setupTemplatesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/mail/templates", r.URL.Path)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListTemplatesResponse{
				Templates: []Template{
					{ID: "tpl-1", Name: "Welcome", Slug: "welcome"},
					{ID: "tpl-2", Name: "Reset Password", Slug: "reset-password"},
				},
				Total:  2,
				Limit:  10,
				Offset: 0,
			})
		})
		defer server.Close()

		resp, err := templatesClient.List(context.Background(), nil)

		require.NoError(t, err)
		assert.Len(t, resp.Templates, 2)
		assert.Equal(t, 2, resp.Total)
	})

	t.Run("with filters", func(t *testing.T) {
		templatesClient, server := setupTemplatesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "limit=5")
			assert.Contains(t, r.URL.RawQuery, "isActive=true")
			assert.Contains(t, r.URL.RawQuery, "search=welcome")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListTemplatesResponse{Templates: []Template{}})
		})
		defer server.Close()

		limit := 5
		isActive := true
		search := "welcome"
		resp, err := templatesClient.List(context.Background(), &ListTemplatesRequest{
			Limit:    &limit,
			IsActive: &isActive,
			Search:   &search,
		})

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func TestTemplatesClient_Get(t *testing.T) {
	templateID := "tpl-123"
	templatesClient, server := setupTemplatesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mail/templates/"+templateID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Template{
			ID:          templateID,
			Name:        "Welcome Email",
			Slug:        "welcome",
			Subject:     "Welcome to our service!",
			HTML:        "<h1>Welcome</h1>",
			IsActive:    true,
			Environment: types.EnvironmentProduction,
			CreatedAt:   time.Now(),
		})
	})
	defer server.Close()

	resp, err := templatesClient.Get(context.Background(), templateID)

	require.NoError(t, err)
	assert.Equal(t, templateID, resp.ID)
	assert.Equal(t, "Welcome Email", resp.Name)
	assert.Equal(t, "welcome", resp.Slug)
	assert.True(t, resp.IsActive)
}

func TestTemplatesClient_GetBySlug(t *testing.T) {
	slug := "welcome"
	templatesClient, server := setupTemplatesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mail/templates/slug/"+slug, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Template{
			ID:   "tpl-123",
			Name: "Welcome Email",
			Slug: slug,
		})
	})
	defer server.Close()

	resp, err := templatesClient.GetBySlug(context.Background(), slug)

	require.NoError(t, err)
	assert.Equal(t, slug, resp.Slug)
}

func TestTemplatesClient_Create(t *testing.T) {
	templatesClient, server := setupTemplatesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/templates", r.URL.Path)

		var req CreateTemplateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "New Template", req.Name)
		assert.Equal(t, "new-template", req.Slug)
		assert.Equal(t, "Subject Line", req.Subject)
		assert.Equal(t, "<p>Content</p>", req.HTML)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Template{
			ID:      "tpl-new",
			Name:    req.Name,
			Slug:    req.Slug,
			Subject: req.Subject,
			HTML:    req.HTML,
		})
	})
	defer server.Close()

	resp, err := templatesClient.Create(context.Background(), &CreateTemplateRequest{
		Name:    "New Template",
		Slug:    "new-template",
		Subject: "Subject Line",
		HTML:    "<p>Content</p>",
	})

	require.NoError(t, err)
	assert.Equal(t, "tpl-new", resp.ID)
	assert.Equal(t, "New Template", resp.Name)
}

func TestTemplatesClient_Update(t *testing.T) {
	templateID := "tpl-123"
	templatesClient, server := setupTemplatesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/mail/templates/"+templateID, r.URL.Path)

		var req UpdateTemplateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", *req.Name)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Template{
			ID:   templateID,
			Name: "Updated Name",
		})
	})
	defer server.Close()

	name := "Updated Name"
	resp, err := templatesClient.Update(context.Background(), &UpdateTemplateRequest{
		ID:   templateID,
		Name: &name,
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated Name", resp.Name)
}

func TestTemplatesClient_Delete(t *testing.T) {
	templateID := "tpl-123"
	templatesClient, server := setupTemplatesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/mail/templates/"+templateID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DeleteTemplateResponse{Success: true})
	})
	defer server.Close()

	resp, err := templatesClient.Delete(context.Background(), templateID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestTemplatesClient_Preview(t *testing.T) {
	templateID := "tpl-123"
	templatesClient, server := setupTemplatesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/templates/"+templateID+"/preview", r.URL.Path)

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		variables := body["variables"].(map[string]interface{})
		assert.Equal(t, "John", variables["name"])

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PreviewTemplateResponse{
			Subject: "Welcome John!",
			HTML:    "<h1>Hello John!</h1>",
		})
	})
	defer server.Close()

	resp, err := templatesClient.Preview(context.Background(), &PreviewTemplateRequest{
		ID: templateID,
		Variables: map[string]interface{}{
			"name": "John",
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "Welcome John!", resp.Subject)
	assert.Contains(t, resp.HTML, "John")
}

func TestTemplatesClient_Preview_WithTextOutput(t *testing.T) {
	templateID := "tpl-123"
	templatesClient, server := setupTemplatesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		text := "Hello World!"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PreviewTemplateResponse{
			Subject: "Subject",
			HTML:    "<p>Hello World!</p>",
			Text:    &text,
		})
	})
	defer server.Close()

	resp, err := templatesClient.Preview(context.Background(), &PreviewTemplateRequest{
		ID:        templateID,
		Variables: map[string]interface{}{},
	})

	require.NoError(t, err)
	assert.NotNil(t, resp.Text)
	assert.Equal(t, "Hello World!", *resp.Text)
}
