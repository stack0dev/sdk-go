package cdn

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stack0/sdk-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetPrivateUploadURL(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/private/upload", r.URL.Path)

		var req PrivateUploadURLRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "secret.pdf", req.Filename)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PrivateUploadURLResponse{
			UploadURL: "https://s3.amazonaws.com/private-bucket/presigned",
			FileID:    "file-123",
			ExpiresAt: time.Now().Add(15 * time.Minute),
		})
	})
	defer server.Close()

	resp, err := cdnClient.GetPrivateUploadURL(context.Background(), &PrivateUploadURLRequest{
		ProjectSlug: "my-project",
		Filename:    "secret.pdf",
		MimeType:    "application/pdf",
		Size:        1024000,
	})

	require.NoError(t, err)
	assert.Equal(t, "file-123", resp.FileID)
	assert.Contains(t, resp.UploadURL, "presigned")
}

func TestClient_ConfirmPrivateUpload(t *testing.T) {
	fileID := "file-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/private/upload/"+fileID+"/confirm", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PrivateFile{
			ID:       fileID,
			Filename: "secret.pdf",
			Status:   PrivateFileStatusReady,
		})
	})
	defer server.Close()

	resp, err := cdnClient.ConfirmPrivateUpload(context.Background(), fileID)

	require.NoError(t, err)
	assert.Equal(t, fileID, resp.ID)
	assert.Equal(t, PrivateFileStatusReady, resp.Status)
}

func TestClient_GetPrivateDownloadURL(t *testing.T) {
	fileID := "file-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/private/"+fileID+"/download", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PrivateDownloadURLResponse{
			DownloadURL: "https://s3.amazonaws.com/private-bucket/presigned-download",
			ExpiresAt:   time.Now().Add(1 * time.Hour),
		})
	})
	defer server.Close()

	resp, err := cdnClient.GetPrivateDownloadURL(context.Background(), &PrivateDownloadURLRequest{
		FileID: fileID,
	})

	require.NoError(t, err)
	assert.Contains(t, resp.DownloadURL, "presigned-download")
}

func TestClient_GetPrivateDownloadURL_WithExpiresIn(t *testing.T) {
	fileID := "file-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, float64(3600), body["expiresIn"])

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PrivateDownloadURLResponse{
			DownloadURL: "https://s3.amazonaws.com/private-bucket/presigned",
			ExpiresAt:   time.Now().Add(1 * time.Hour),
		})
	})
	defer server.Close()

	expiresIn := 3600
	resp, err := cdnClient.GetPrivateDownloadURL(context.Background(), &PrivateDownloadURLRequest{
		FileID:    fileID,
		ExpiresIn: &expiresIn,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.DownloadURL)
}

func TestClient_GetPrivateFile(t *testing.T) {
	fileID := "file-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cdn/private/"+fileID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PrivateFile{
			ID:       fileID,
			Filename: "confidential.docx",
			MimeType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			Size:     50000,
			Status:   PrivateFileStatusReady,
		})
	})
	defer server.Close()

	resp, err := cdnClient.GetPrivateFile(context.Background(), fileID)

	require.NoError(t, err)
	assert.Equal(t, fileID, resp.ID)
	assert.Equal(t, "confidential.docx", resp.Filename)
}

func TestClient_UpdatePrivateFile(t *testing.T) {
	fileID := "file-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/cdn/private/"+fileID, r.URL.Path)

		var req UpdatePrivateFileRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "new-name.pdf", *req.Filename)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PrivateFile{
			ID:       fileID,
			Filename: "new-name.pdf",
		})
	})
	defer server.Close()

	filename := "new-name.pdf"
	resp, err := cdnClient.UpdatePrivateFile(context.Background(), &UpdatePrivateFileRequest{
		FileID:   fileID,
		Filename: &filename,
	})

	require.NoError(t, err)
	assert.Equal(t, "new-name.pdf", resp.Filename)
}

func TestClient_DeletePrivateFile(t *testing.T) {
	fileID := "file-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/cdn/private/"+fileID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SuccessResponse{Success: true})
	})
	defer server.Close()

	resp, err := cdnClient.DeletePrivateFile(context.Background(), fileID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestClient_DeletePrivateFiles(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/private/delete", r.URL.Path)

		var body map[string][]string
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		assert.Len(t, body["fileIds"], 3)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DeleteAssetsResponse{
			Success:      true,
			DeletedCount: 3,
		})
	})
	defer server.Close()

	resp, err := cdnClient.DeletePrivateFiles(context.Background(), []string{"file-1", "file-2", "file-3"})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, 3, resp.DeletedCount)
}

func TestClient_ListPrivateFiles(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/cdn/private")
		assert.Contains(t, r.URL.RawQuery, "projectSlug=my-project")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListPrivateFilesResponse{
			Files: []PrivateFile{
				{ID: "file-1", Filename: "doc1.pdf"},
				{ID: "file-2", Filename: "doc2.pdf"},
			},
			Total:   2,
			HasMore: false,
		})
	})
	defer server.Close()

	resp, err := cdnClient.ListPrivateFiles(context.Background(), &ListPrivateFilesRequest{
		ProjectSlug: "my-project",
	})

	require.NoError(t, err)
	assert.Len(t, resp.Files, 2)
	assert.False(t, resp.HasMore)
}

func TestClient_MovePrivateFiles(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/private/move", r.URL.Path)

		var req MovePrivateFilesRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Len(t, req.FileIDs, 2)
		assert.Equal(t, "new-folder", *req.Folder)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MovePrivateFilesResponse{
			Success:    true,
			MovedCount: 2,
		})
	})
	defer server.Close()

	folder := "new-folder"
	resp, err := cdnClient.MovePrivateFiles(context.Background(), &MovePrivateFilesRequest{
		FileIDs: []string{"file-1", "file-2"},
		Folder:  &folder,
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, 2, resp.MovedCount)
}

func TestClient_CreateBundle(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/bundles", r.URL.Path)

		var req CreateBundleRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "My Bundle", req.Name)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CreateBundleResponse{
			ID:     "bundle-123",
			Status: BundleStatusPending,
		})
	})
	defer server.Close()

	resp, err := cdnClient.CreateBundle(context.Background(), &CreateBundleRequest{
		ProjectSlug: "my-project",
		Name:        "My Bundle",
		AssetIDs:    []string{"asset-1", "asset-2"},
	})

	require.NoError(t, err)
	assert.Equal(t, "bundle-123", resp.ID)
	assert.Equal(t, BundleStatusPending, resp.Status)
}

func TestClient_GetBundle(t *testing.T) {
	bundleID := "bundle-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cdn/bundles/"+bundleID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DownloadBundle{
			ID:     bundleID,
			Name:   "My Bundle",
			Status: BundleStatusReady,
		})
	})
	defer server.Close()

	resp, err := cdnClient.GetBundle(context.Background(), bundleID)

	require.NoError(t, err)
	assert.Equal(t, bundleID, resp.ID)
	assert.Equal(t, BundleStatusReady, resp.Status)
}

func TestClient_ListBundles(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/cdn/bundles")
		assert.Contains(t, r.URL.RawQuery, "projectSlug=my-project")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListBundlesResponse{
			Bundles: []DownloadBundle{
				{ID: "bundle-1", Name: "Bundle 1"},
				{ID: "bundle-2", Name: "Bundle 2"},
			},
			Total: 2,
		})
	})
	defer server.Close()

	resp, err := cdnClient.ListBundles(context.Background(), &ListBundlesRequest{
		ProjectSlug: "my-project",
	})

	require.NoError(t, err)
	assert.Len(t, resp.Bundles, 2)
}

func TestClient_GetBundleDownloadURL(t *testing.T) {
	bundleID := "bundle-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/bundles/"+bundleID+"/download", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(BundleDownloadURLResponse{
			DownloadURL: "https://s3.amazonaws.com/bundles/bundle.zip",
			ExpiresAt:   time.Now().Add(1 * time.Hour),
		})
	})
	defer server.Close()

	resp, err := cdnClient.GetBundleDownloadURL(context.Background(), &BundleDownloadURLRequest{
		BundleID: bundleID,
	})

	require.NoError(t, err)
	assert.Contains(t, resp.DownloadURL, "bundle.zip")
}

func TestClient_DeleteBundle(t *testing.T) {
	bundleID := "bundle-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/cdn/bundles/"+bundleID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SuccessResponse{Success: true})
	})
	defer server.Close()

	resp, err := cdnClient.DeleteBundle(context.Background(), bundleID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}
