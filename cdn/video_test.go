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

func TestClient_Transcode(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/video/transcode", r.URL.Path)

		var req TranscodeVideoRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "asset-123", req.AssetID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(TranscodeJob{
			ID:       "job-123",
			AssetID:  req.AssetID,
			Status:   TranscodeJobStatusPending,
			Progress: 0,
		})
	})
	defer server.Close()

	resp, err := cdnClient.Transcode(context.Background(), &TranscodeVideoRequest{
		AssetID: "asset-123",
	})

	require.NoError(t, err)
	assert.Equal(t, "job-123", resp.ID)
	assert.Equal(t, TranscodeJobStatusPending, resp.Status)
}

func TestClient_GetJob(t *testing.T) {
	jobID := "job-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cdn/video/jobs/"+jobID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(TranscodeJob{
			ID:       jobID,
			AssetID:  "asset-123",
			Status:   TranscodeJobStatusProcessing,
			Progress: 50,
		})
	})
	defer server.Close()

	resp, err := cdnClient.GetJob(context.Background(), jobID)

	require.NoError(t, err)
	assert.Equal(t, jobID, resp.ID)
	assert.Equal(t, TranscodeJobStatusProcessing, resp.Status)
	assert.Equal(t, 50, resp.Progress)
}

func TestClient_ListJobs(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/cdn/video/jobs")
		assert.Contains(t, r.URL.RawQuery, "projectSlug=my-project")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListJobsResponse{
			Jobs: []TranscodeJob{
				{ID: "job-1", Status: TranscodeJobStatusCompleted},
				{ID: "job-2", Status: TranscodeJobStatusProcessing},
			},
			Total: 2,
		})
	})
	defer server.Close()

	resp, err := cdnClient.ListJobs(context.Background(), &ListJobsRequest{
		ProjectSlug: "my-project",
	})

	require.NoError(t, err)
	assert.Len(t, resp.Jobs, 2)
}

func TestClient_CancelJob(t *testing.T) {
	jobID := "job-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/video/jobs/"+jobID+"/cancel", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SuccessResponse{Success: true})
	})
	defer server.Close()

	resp, err := cdnClient.CancelJob(context.Background(), jobID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestClient_GetStreamingURLs(t *testing.T) {
	assetID := "asset-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cdn/video/stream/"+assetID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(StreamingURLs{
			HLS:  "https://cdn.example.com/stream/video.m3u8",
			DASH: "https://cdn.example.com/stream/video.mpd",
		})
	})
	defer server.Close()

	resp, err := cdnClient.GetStreamingURLs(context.Background(), assetID)

	require.NoError(t, err)
	assert.Contains(t, resp.HLS, "m3u8")
	assert.Contains(t, resp.DASH, "mpd")
}

func TestClient_GetThumbnail(t *testing.T) {
	assetID := "asset-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/cdn/video/thumbnail/"+assetID)
		assert.Contains(t, r.URL.RawQuery, "timestamp=10")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ThumbnailResponse{
			URL:    "https://cdn.example.com/thumbnail.jpg",
			Width:  320,
			Height: 180,
		})
	})
	defer server.Close()

	resp, err := cdnClient.GetThumbnail(context.Background(), &ThumbnailRequest{
		AssetID:   assetID,
		Timestamp: 10.0,
	})

	require.NoError(t, err)
	assert.Contains(t, resp.URL, "thumbnail.jpg")
	assert.Equal(t, 320, resp.Width)
}

func TestClient_RegenerateThumbnail(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/video/thumbnail/regenerate", r.URL.Path)

		var req RegenerateThumbnailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "asset-123", req.AssetID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(RegenerateThumbnailResponse{
			Success:      true,
			ThumbnailURL: "https://cdn.example.com/new-thumbnail.jpg",
		})
	})
	defer server.Close()

	resp, err := cdnClient.RegenerateThumbnail(context.Background(), &RegenerateThumbnailRequest{
		AssetID: "asset-123",
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Contains(t, resp.ThumbnailURL, "new-thumbnail")
}

func TestClient_ListThumbnails(t *testing.T) {
	assetID := "asset-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cdn/video/"+assetID+"/thumbnails", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListThumbnailsResponse{
			Thumbnails: []VideoThumbnail{
				{URL: "https://cdn.example.com/thumb1.jpg", Timestamp: 0},
				{URL: "https://cdn.example.com/thumb2.jpg", Timestamp: 30},
			},
		})
	})
	defer server.Close()

	resp, err := cdnClient.ListThumbnails(context.Background(), assetID)

	require.NoError(t, err)
	assert.Len(t, resp.Thumbnails, 2)
}

func TestClient_ExtractAudio(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/video/extract-audio", r.URL.Path)

		var req ExtractAudioRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "asset-123", req.AssetID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ExtractAudioResponse{
			Success:  true,
			AudioURL: "https://cdn.example.com/audio.mp3",
		})
	})
	defer server.Close()

	resp, err := cdnClient.ExtractAudio(context.Background(), &ExtractAudioRequest{
		AssetID: "asset-123",
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Contains(t, resp.AudioURL, "audio.mp3")
}

func TestClient_GenerateGif(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/video/gif", r.URL.Path)

		var req GenerateGifRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "asset-123", req.AssetID)
		assert.Equal(t, 5.0, *req.StartTime)
		assert.Equal(t, 3.0, *req.Duration)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(VideoGif{
			ID:        "gif-123",
			AssetID:   req.AssetID,
			URL:       "https://cdn.example.com/video.gif",
			StartTime: 5.0,
			Duration:  3.0,
		})
	})
	defer server.Close()

	startTime := 5.0
	duration := 3.0
	resp, err := cdnClient.GenerateGif(context.Background(), &GenerateGifRequest{
		AssetID:   "asset-123",
		StartTime: &startTime,
		Duration:  &duration,
	})

	require.NoError(t, err)
	assert.Equal(t, "gif-123", resp.ID)
	assert.Equal(t, 5.0, resp.StartTime)
}

func TestClient_GetGif(t *testing.T) {
	gifID := "gif-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cdn/video/gif/"+gifID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(VideoGif{
			ID:      gifID,
			AssetID: "asset-123",
			URL:     "https://cdn.example.com/video.gif",
		})
	})
	defer server.Close()

	resp, err := cdnClient.GetGif(context.Background(), gifID)

	require.NoError(t, err)
	assert.Equal(t, gifID, resp.ID)
}

func TestClient_ListGifs(t *testing.T) {
	assetID := "asset-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cdn/video/"+assetID+"/gifs", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]VideoGif{
			{ID: "gif-1", AssetID: assetID},
			{ID: "gif-2", AssetID: assetID},
		})
	})
	defer server.Close()

	resp, err := cdnClient.ListGifs(context.Background(), assetID)

	require.NoError(t, err)
	assert.Len(t, resp, 2)
}

func TestClient_CreateMergeJob(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/video/merge", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MergeJob{
			ID:     "merge-123",
			Status: MergeJobStatusPending,
		})
	})
	defer server.Close()

	resp, err := cdnClient.CreateMergeJob(context.Background(), &CreateMergeJobRequest{
		ProjectSlug: "my-project",
		Inputs:      []MergeInput{{AssetID: "asset-1"}, {AssetID: "asset-2"}},
	})

	require.NoError(t, err)
	assert.Equal(t, "merge-123", resp.ID)
	assert.Equal(t, MergeJobStatusPending, resp.Status)
}

func TestClient_GetMergeJob(t *testing.T) {
	jobID := "merge-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cdn/video/merge/"+jobID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MergeJobWithOutput{
			MergeJob: MergeJob{
				ID:     jobID,
				Status: MergeJobStatusCompleted,
			},
			OutputAsset: &Asset{
				ID:       "output-asset",
				Filename: "merged.mp4",
			},
		})
	})
	defer server.Close()

	resp, err := cdnClient.GetMergeJob(context.Background(), jobID)

	require.NoError(t, err)
	assert.Equal(t, jobID, resp.ID)
	assert.Equal(t, MergeJobStatusCompleted, resp.Status)
	assert.NotNil(t, resp.OutputAsset)
}

func TestClient_ListMergeJobs(t *testing.T) {
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/cdn/video/merge")
		assert.Contains(t, r.URL.RawQuery, "projectSlug=my-project")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListMergeJobsResponse{
			Jobs:  []MergeJob{{ID: "merge-1"}, {ID: "merge-2"}},
			Total: 2,
		})
	})
	defer server.Close()

	resp, err := cdnClient.ListMergeJobs(context.Background(), &ListMergeJobsRequest{
		ProjectSlug: "my-project",
	})

	require.NoError(t, err)
	assert.Len(t, resp.Jobs, 2)
}

func TestClient_CancelMergeJob(t *testing.T) {
	jobID := "merge-123"
	cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/cdn/video/merge/"+jobID+"/cancel", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SuccessResponse{Success: true})
	})
	defer server.Close()

	resp, err := cdnClient.CancelMergeJob(context.Background(), jobID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}
