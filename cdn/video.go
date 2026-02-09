package cdn

import (
	"context"
	"net/url"
	"strconv"
)

// Transcode starts a video transcoding job.
func (c *Client) Transcode(ctx context.Context, req *TranscodeVideoRequest) (*TranscodeJob, error) {
	var resp TranscodeJob
	if err := c.http.Post(ctx, "/cdn/video/transcode", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetJob retrieves a transcoding job by ID.
func (c *Client) GetJob(ctx context.Context, jobID string) (*TranscodeJob, error) {
	var resp TranscodeJob
	if err := c.http.Get(ctx, "/cdn/video/jobs/"+jobID, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListJobs lists transcoding jobs with filters.
func (c *Client) ListJobs(ctx context.Context, req *ListJobsRequest) (*ListJobsResponse, error) {
	params := url.Values{}
	params.Set("projectSlug", req.ProjectSlug)
	if req.AssetID != nil {
		params.Set("assetId", *req.AssetID)
	}
	if req.Status != nil {
		params.Set("status", string(*req.Status))
	}
	if req.Limit != nil {
		params.Set("limit", strconv.Itoa(*req.Limit))
	}
	if req.Offset != nil {
		params.Set("offset", strconv.Itoa(*req.Offset))
	}

	var resp ListJobsResponse
	if err := c.http.Get(ctx, "/cdn/video/jobs?"+params.Encode(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelJob cancels a pending or processing transcoding job.
func (c *Client) CancelJob(ctx context.Context, jobID string) (*SuccessResponse, error) {
	var resp SuccessResponse
	if err := c.http.Post(ctx, "/cdn/video/jobs/"+jobID+"/cancel", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetStreamingURLs gets streaming URLs for a transcoded video.
func (c *Client) GetStreamingURLs(ctx context.Context, assetID string) (*StreamingURLs, error) {
	var resp StreamingURLs
	if err := c.http.Get(ctx, "/cdn/video/stream/"+assetID, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetThumbnail generates a thumbnail from a video.
func (c *Client) GetThumbnail(ctx context.Context, req *ThumbnailRequest) (*ThumbnailResponse, error) {
	params := url.Values{}
	params.Set("timestamp", strconv.FormatFloat(req.Timestamp, 'f', -1, 64))
	if req.Width != nil {
		params.Set("width", strconv.Itoa(*req.Width))
	}
	if req.Format != nil {
		params.Set("format", *req.Format)
	}

	var resp ThumbnailResponse
	if err := c.http.Get(ctx, "/cdn/video/thumbnail/"+req.AssetID+"?"+params.Encode(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RegenerateThumbnail regenerates a thumbnail for a video.
func (c *Client) RegenerateThumbnail(ctx context.Context, req *RegenerateThumbnailRequest) (*RegenerateThumbnailResponse, error) {
	var resp RegenerateThumbnailResponse
	if err := c.http.Post(ctx, "/cdn/video/thumbnail/regenerate", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListThumbnails lists all thumbnails for a video asset.
func (c *Client) ListThumbnails(ctx context.Context, assetID string) (*ListThumbnailsResponse, error) {
	var resp ListThumbnailsResponse
	if err := c.http.Get(ctx, "/cdn/video/"+assetID+"/thumbnails", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ExtractAudio extracts audio from a video file.
func (c *Client) ExtractAudio(ctx context.Context, req *ExtractAudioRequest) (*ExtractAudioResponse, error) {
	var resp ExtractAudioResponse
	if err := c.http.Post(ctx, "/cdn/video/extract-audio", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GenerateGif generates an animated GIF from a video segment.
func (c *Client) GenerateGif(ctx context.Context, req *GenerateGifRequest) (*VideoGif, error) {
	var resp VideoGif
	if err := c.http.Post(ctx, "/cdn/video/gif", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetGif retrieves a specific GIF by ID.
func (c *Client) GetGif(ctx context.Context, gifID string) (*VideoGif, error) {
	var resp VideoGif
	if err := c.http.Get(ctx, "/cdn/video/gif/"+gifID, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListGifs lists all GIFs generated for a video asset.
func (c *Client) ListGifs(ctx context.Context, assetID string) ([]VideoGif, error) {
	var resp []VideoGif
	if err := c.http.Get(ctx, "/cdn/video/"+assetID+"/gifs", &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateMergeJob creates a merge job to combine videos/images.
func (c *Client) CreateMergeJob(ctx context.Context, req *CreateMergeJobRequest) (*MergeJob, error) {
	var resp MergeJob
	if err := c.http.Post(ctx, "/cdn/video/merge", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetMergeJob retrieves a merge job by ID with output details.
func (c *Client) GetMergeJob(ctx context.Context, jobID string) (*MergeJobWithOutput, error) {
	var resp MergeJobWithOutput
	if err := c.http.Get(ctx, "/cdn/video/merge/"+jobID, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListMergeJobs lists merge jobs with optional filters.
func (c *Client) ListMergeJobs(ctx context.Context, req *ListMergeJobsRequest) (*ListMergeJobsResponse, error) {
	params := url.Values{}
	params.Set("projectSlug", req.ProjectSlug)
	if req.Status != nil {
		params.Set("status", string(*req.Status))
	}
	if req.Limit != nil {
		params.Set("limit", strconv.Itoa(*req.Limit))
	}
	if req.Offset != nil {
		params.Set("offset", strconv.Itoa(*req.Offset))
	}

	var resp ListMergeJobsResponse
	if err := c.http.Get(ctx, "/cdn/video/merge?"+params.Encode(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelMergeJob cancels a pending or processing merge job.
func (c *Client) CancelMergeJob(ctx context.Context, jobID string) (*SuccessResponse, error) {
	var resp SuccessResponse
	if err := c.http.Post(ctx, "/cdn/video/merge/"+jobID+"/cancel", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
