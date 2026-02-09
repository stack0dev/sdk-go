package cdn

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stack0/sdk-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateImport(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/imports", r.URL.Path)

		var req CreateImportRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "my-project", req.ProjectSlug)
		assert.Equal(t, "my-bucket", req.SourceBucket)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CreateImportResponse{
			ID:         "import-123",
			Status:     ImportJobStatusPending,
			TotalFiles: 100,
		})
	})
	defer server.Close()

	resp, err := cdnClient.CreateImport(context.Background(), &CreateImportRequest{
		ProjectSlug:  "my-project",
		SourceBucket: "my-bucket",
		SourcePrefix: "uploads/",
	})

	require.NoError(t, err)
	assert.Equal(t, "import-123", resp.ID)
	assert.Equal(t, ImportJobStatusPending, resp.Status)
	assert.Equal(t, 100, resp.TotalFiles)
}

func TestClient_GetImport(t *testing.T) {
	importID := "import-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cdn/imports/"+importID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ImportJob{
			ID:             importID,
			Status:         ImportJobStatusProcessing,
			TotalFiles:     100,
			ProcessedFiles: 50,
			SuccessfulFiles: 48,
			FailedFiles:    2,
			Progress:       50,
			CreatedAt:      time.Now(),
		})
	})
	defer server.Close()

	resp, err := cdnClient.GetImport(context.Background(), importID)

	require.NoError(t, err)
	assert.Equal(t, importID, resp.ID)
	assert.Equal(t, ImportJobStatusProcessing, resp.Status)
	assert.Equal(t, 50, resp.ProcessedFiles)
	assert.Equal(t, 50, resp.Progress)
}

func TestClient_ListImports(t *testing.T) {
	t.Run("basic request", func(t *testing.T) {
		cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Contains(t, r.URL.Path, "/cdn/imports")
			assert.Contains(t, r.URL.RawQuery, "projectSlug=my-project")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListImportsResponse{
				Imports: []ImportJob{
					{ID: "import-1", Status: ImportJobStatusCompleted},
					{ID: "import-2", Status: ImportJobStatusProcessing},
				},
				Total:   2,
				HasMore: false,
			})
		})
		defer server.Close()

		resp, err := cdnClient.ListImports(context.Background(), &ListImportsRequest{
			ProjectSlug: "my-project",
		})

		require.NoError(t, err)
		assert.Len(t, resp.Imports, 2)
		assert.False(t, resp.HasMore)
	})

	t.Run("with filters", func(t *testing.T) {
		cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "status=completed")
			assert.Contains(t, r.URL.RawQuery, "environment=production")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListImportsResponse{Imports: []ImportJob{}})
		})
		defer server.Close()

		status := ImportJobStatusCompleted
		env := types.EnvironmentProduction
		resp, err := cdnClient.ListImports(context.Background(), &ListImportsRequest{
			ProjectSlug: "my-project",
			Status:      &status,
			Environment: &env,
		})

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func TestClient_CancelImport(t *testing.T) {
	importID := "import-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/imports/"+importID+"/cancel", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CancelImportResponse{
			Success: true,
		})
	})
	defer server.Close()

	resp, err := cdnClient.CancelImport(context.Background(), importID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestClient_RetryImport(t *testing.T) {
	importID := "import-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/imports/"+importID+"/retry", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(RetryImportResponse{
			Success:      true,
			RetryingFiles: 5,
		})
	})
	defer server.Close()

	resp, err := cdnClient.RetryImport(context.Background(), importID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, 5, resp.RetryingFiles)
}

func TestClient_ListImportFiles(t *testing.T) {
	importID := "import-123"
	t.Run("basic request", func(t *testing.T) {
		cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Contains(t, r.URL.Path, "/cdn/imports/"+importID+"/files")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListImportFilesResponse{
				Files: []ImportFile{
					{ID: "file-1", SourceKey: "uploads/image1.jpg", Status: ImportFileStatusCompleted},
					{ID: "file-2", SourceKey: "uploads/image2.jpg", Status: ImportFileStatusFailed},
				},
				Total:   2,
				HasMore: false,
			})
		})
		defer server.Close()

		resp, err := cdnClient.ListImportFiles(context.Background(), &ListImportFilesRequest{
			ImportID: importID,
		})

		require.NoError(t, err)
		assert.Len(t, resp.Files, 2)
	})

	t.Run("with status filter", func(t *testing.T) {
		cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "status=failed")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ListImportFilesResponse{
				Files: []ImportFile{
					{ID: "file-2", Status: ImportFileStatusFailed, Error: ptr("File too large")},
				},
			})
		})
		defer server.Close()

		status := ImportFileStatusFailed
		resp, err := cdnClient.ListImportFiles(context.Background(), &ListImportFilesRequest{
			ImportID: importID,
			Status:   &status,
		})

		require.NoError(t, err)
		assert.Len(t, resp.Files, 1)
		assert.Equal(t, ImportFileStatusFailed, resp.Files[0].Status)
	})
}

func TestImportJobStatus_Constants(t *testing.T) {
	assert.Equal(t, ImportJobStatus("pending"), ImportJobStatusPending)
	assert.Equal(t, ImportJobStatus("processing"), ImportJobStatusProcessing)
	assert.Equal(t, ImportJobStatus("completed"), ImportJobStatusCompleted)
	assert.Equal(t, ImportJobStatus("failed"), ImportJobStatusFailed)
	assert.Equal(t, ImportJobStatus("cancelled"), ImportJobStatusCancelled)
}

func TestImportFileStatus_Constants(t *testing.T) {
	assert.Equal(t, ImportFileStatus("pending"), ImportFileStatusPending)
	assert.Equal(t, ImportFileStatus("processing"), ImportFileStatusProcessing)
	assert.Equal(t, ImportFileStatus("completed"), ImportFileStatusCompleted)
	assert.Equal(t, ImportFileStatus("failed"), ImportFileStatusFailed)
	assert.Equal(t, ImportFileStatus("skipped"), ImportFileStatusSkipped)
}

// Helper function
func ptr[T any](v T) *T {
	return &v
}
