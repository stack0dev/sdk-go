package cdn

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

func setupCDNTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	httpClient := client.New("test-api-key", server.URL)
	return NewClient(httpClient, ""), server
}

func setupCDNTestClientWithCDNURL(t *testing.T, handler http.HandlerFunc, cdnURL string) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	httpClient := client.New("test-api-key", server.URL)
	return NewClient(httpClient, cdnURL), server
}

func TestClient_GetUploadURL(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/upload", r.URL.Path)

		var req UploadURLRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "my-project", req.ProjectSlug)
		assert.Equal(t, "image.jpg", req.Filename)
		assert.Equal(t, "image/jpeg", req.MimeType)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(UploadURLResponse{
			UploadURL: "https://s3.amazonaws.com/bucket/presigned-url",
			AssetID:   "asset-123",
			CDNURL:    "https://cdn.example.com/asset-123.jpg",
			ExpiresAt: time.Now().Add(15 * time.Minute),
		})
	})
	defer server.Close()

	resp, err := cdnClient.GetUploadURL(context.Background(), &UploadURLRequest{
		ProjectSlug: "my-project",
		Filename:    "image.jpg",
		MimeType:    "image/jpeg",
		Size:        1024,
	})

	require.NoError(t, err)
	assert.Equal(t, "asset-123", resp.AssetID)
	assert.Contains(t, resp.UploadURL, "presigned-url")
}

func TestClient_ConfirmUpload(t *testing.T) {
	assetID := "asset-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/upload/"+assetID+"/confirm", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Asset{
			ID:       assetID,
			Filename: "image.jpg",
			Status:   AssetStatusReady,
		})
	})
	defer server.Close()

	resp, err := cdnClient.ConfirmUpload(context.Background(), assetID)

	require.NoError(t, err)
	assert.Equal(t, assetID, resp.ID)
	assert.Equal(t, AssetStatusReady, resp.Status)
}

func TestClient_Get(t *testing.T) {
	assetID := "asset-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cdn/assets/"+assetID, r.URL.Path)

		width := 1920
		height := 1080
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Asset{
			ID:               assetID,
			Filename:         "photo.jpg",
			OriginalFilename: "photo.jpg",
			MimeType:         "image/jpeg",
			Size:             2048000,
			Type:             AssetTypeImage,
			Status:           AssetStatusReady,
			Width:            &width,
			Height:           &height,
			CDNURL:           "https://cdn.example.com/photo.jpg",
		})
	})
	defer server.Close()

	resp, err := cdnClient.Get(context.Background(), assetID)

	require.NoError(t, err)
	assert.Equal(t, assetID, resp.ID)
	assert.Equal(t, AssetTypeImage, resp.Type)
	assert.Equal(t, 1920, *resp.Width)
}

func TestClient_Update(t *testing.T) {
	assetID := "asset-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/cdn/assets/"+assetID, r.URL.Path)

		var req UpdateAssetRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "new-name.jpg", *req.Filename)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Asset{
			ID:       assetID,
			Filename: "new-name.jpg",
		})
	})
	defer server.Close()

	filename := "new-name.jpg"
	resp, err := cdnClient.Update(context.Background(), &UpdateAssetRequest{
		ID:       assetID,
		Filename: &filename,
	})

	require.NoError(t, err)
	assert.Equal(t, "new-name.jpg", resp.Filename)
}

func TestClient_Delete(t *testing.T) {
	assetID := "asset-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/cdn/assets/"+assetID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SuccessResponse{Success: true})
	})
	defer server.Close()

	resp, err := cdnClient.Delete(context.Background(), assetID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestClient_DeleteMany(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/assets/delete", r.URL.Path)

		var body map[string][]string
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		assert.Len(t, body["ids"], 3)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DeleteAssetsResponse{
			Success:      true,
			DeletedCount: 3,
		})
	})
	defer server.Close()

	resp, err := cdnClient.DeleteMany(context.Background(), []string{"asset-1", "asset-2", "asset-3"})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, 3, resp.DeletedCount)
}

func TestClient_List(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/cdn/assets")
		assert.Contains(t, r.URL.RawQuery, "projectSlug=my-project")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListAssetsResponse{
			Assets: []Asset{
				{ID: "asset-1", Filename: "image1.jpg"},
				{ID: "asset-2", Filename: "image2.jpg"},
			},
			Total:   2,
			HasMore: false,
		})
	})
	defer server.Close()

	resp, err := cdnClient.List(context.Background(), &ListAssetsRequest{
		ProjectSlug: "my-project",
	})

	require.NoError(t, err)
	assert.Len(t, resp.Assets, 2)
	assert.False(t, resp.HasMore)
}

func TestClient_List_WithFilters(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.RawQuery, "type=image")
		assert.Contains(t, r.URL.RawQuery, "status=ready")
		assert.Contains(t, r.URL.RawQuery, "folder=photos")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListAssetsResponse{Assets: []Asset{}})
	})
	defer server.Close()

	assetType := AssetTypeImage
	status := AssetStatusReady
	folder := "photos"
	resp, err := cdnClient.List(context.Background(), &ListAssetsRequest{
		ProjectSlug: "my-project",
		Type:        &assetType,
		Status:      &status,
		Folder:      &folder,
	})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestClient_Move(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/assets/move", r.URL.Path)

		var req MoveAssetsRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Len(t, req.AssetIDs, 2)
		assert.Equal(t, "new-folder", *req.Folder)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MoveAssetsResponse{
			Success:    true,
			MovedCount: 2,
		})
	})
	defer server.Close()

	folder := "new-folder"
	resp, err := cdnClient.Move(context.Background(), &MoveAssetsRequest{
		AssetIDs: []string{"asset-1", "asset-2"},
		Folder:   &folder,
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, 2, resp.MovedCount)
}

func TestClient_GetTransformURL(t *testing.T) {
	t.Run("with full URL", func(t *testing.T) {
		cdnClient, _ := setupCDNTestClient(t, nil)

		width := 800
		quality := 80
		url, err := cdnClient.GetTransformURL(
			"https://cdn.example.com/image.jpg",
			&TransformOptions{
				Width:   &width,
				Quality: &quality,
			},
		)

		require.NoError(t, err)
		assert.Contains(t, url, "https://cdn.example.com/image.jpg")
		assert.Contains(t, url, "w=")
		assert.Contains(t, url, "q=80")
	})

	t.Run("with S3 key and CDN URL configured", func(t *testing.T) {
		cdnClient, _ := setupCDNTestClientWithCDNURL(t, nil, "https://cdn.example.com")

		width := 640
		url, err := cdnClient.GetTransformURL(
			"uploads/image.jpg",
			&TransformOptions{
				Width: &width,
			},
		)

		require.NoError(t, err)
		assert.Contains(t, url, "https://cdn.example.com/uploads/image.jpg")
		assert.Contains(t, url, "w=")
	})

	t.Run("error without CDN URL", func(t *testing.T) {
		cdnClient, _ := setupCDNTestClient(t, nil)

		_, err := cdnClient.GetTransformURL(
			"uploads/image.jpg",
			&TransformOptions{},
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "cdnURL")
	})

	t.Run("with all transform options", func(t *testing.T) {
		cdnClient, _ := setupCDNTestClient(t, nil)

		width := 800
		height := 600
		quality := 85
		blur := 5
		rotate := 90
		format := "webp"
		fit := "cover"
		url, err := cdnClient.GetTransformURL(
			"https://cdn.example.com/image.jpg",
			&TransformOptions{
				Width:     &width,
				Height:    &height,
				Quality:   &quality,
				Blur:      &blur,
				Rotate:    &rotate,
				Format:    &format,
				Fit:       &fit,
				Grayscale: true,
				Flip:      true,
				Flop:      true,
			},
		)

		require.NoError(t, err)
		assert.Contains(t, url, "w=")
		assert.Contains(t, url, "h=600")
		assert.Contains(t, url, "q=85")
		assert.Contains(t, url, "blur=5")
		assert.Contains(t, url, "rotate=90")
		assert.Contains(t, url, "f=webp")
		assert.Contains(t, url, "fit=cover")
		assert.Contains(t, url, "grayscale=true")
		assert.Contains(t, url, "flip=y")
		assert.Contains(t, url, "flop=x")
	})
}

func TestClient_GetFolderTree(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/cdn/folders/tree")
		assert.Contains(t, r.URL.RawQuery, "projectSlug=my-project")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"tree": []FolderTreeNode{
				{ID: "folder-1", Name: "Images", Path: "/Images", AssetCount: 10},
				{ID: "folder-2", Name: "Documents", Path: "/Documents", AssetCount: 5},
			},
		})
	})
	defer server.Close()

	resp, err := cdnClient.GetFolderTree(context.Background(), &GetFolderTreeRequest{
		ProjectSlug: "my-project",
	})

	require.NoError(t, err)
	assert.Len(t, resp, 2)
}

func TestClient_CreateFolder(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/folders", r.URL.Path)

		var req CreateFolderRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "New Folder", req.Name)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Folder{
			ID:   "folder-new",
			Name: req.Name,
			Path: "/" + req.Name,
		})
	})
	defer server.Close()

	resp, err := cdnClient.CreateFolder(context.Background(), &CreateFolderRequest{
		ProjectSlug: "my-project",
		Name:        "New Folder",
	})

	require.NoError(t, err)
	assert.Equal(t, "folder-new", resp.ID)
}

func TestClient_GetFolder(t *testing.T) {
	folderID := "folder-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cdn/folders/"+folderID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Folder{
			ID:         folderID,
			Name:       "Images",
			Path:       "/Images",
			AssetCount: 25,
		})
	})
	defer server.Close()

	resp, err := cdnClient.GetFolder(context.Background(), folderID)

	require.NoError(t, err)
	assert.Equal(t, folderID, resp.ID)
	assert.Equal(t, 25, resp.AssetCount)
}

func TestClient_DeleteFolder(t *testing.T) {
	folderID := "folder-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Contains(t, r.URL.Path, "/cdn/folders/"+folderID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SuccessResponse{Success: true})
	})
	defer server.Close()

	resp, err := cdnClient.DeleteFolder(context.Background(), folderID, false)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestAssetStatus_Constants(t *testing.T) {
	assert.Equal(t, AssetStatus("pending"), AssetStatusPending)
	assert.Equal(t, AssetStatus("processing"), AssetStatusProcessing)
	assert.Equal(t, AssetStatus("ready"), AssetStatusReady)
	assert.Equal(t, AssetStatus("failed"), AssetStatusFailed)
	assert.Equal(t, AssetStatus("deleted"), AssetStatusDeleted)
}

func TestAssetType_Constants(t *testing.T) {
	assert.Equal(t, AssetType("image"), AssetTypeImage)
	assert.Equal(t, AssetType("video"), AssetTypeVideo)
	assert.Equal(t, AssetType("audio"), AssetTypeAudio)
	assert.Equal(t, AssetType("document"), AssetTypeDocument)
	assert.Equal(t, AssetType("other"), AssetTypeOther)
}

func TestAllowedWidths(t *testing.T) {
	expected := []int{256, 384, 640, 750, 828, 1080, 1200, 1920, 2048, 3840}
	assert.Equal(t, expected, AllowedWidths)
}

func TestClient_GetNearestWidth(t *testing.T) {
	cdnClient, _ := setupCDNTestClient(t, nil)

	testCases := []struct {
		input    int
		expected int
	}{
		{100, 256},
		{300, 256},
		{400, 384},
		{600, 640},
		{700, 750},
		{800, 828},
		{1000, 1080},
		{1100, 1080},
		{1500, 1200},
		{1800, 1920},
		{2000, 2048},
		{3000, 2048},
		{4000, 3840},
	}

	for _, tc := range testCases {
		result := cdnClient.getNearestWidth(tc.input)
		assert.Equal(t, tc.expected, result, "getNearestWidth(%d) should be %d", tc.input, tc.expected)
	}
}
