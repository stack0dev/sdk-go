package cdn

import "time"

// VideoQuality represents video output quality.
type VideoQuality string

const (
	VideoQuality360p  VideoQuality = "360p"
	VideoQuality480p  VideoQuality = "480p"
	VideoQuality720p  VideoQuality = "720p"
	VideoQuality1080p VideoQuality = "1080p"
	VideoQuality1440p VideoQuality = "1440p"
	VideoQuality2160p VideoQuality = "2160p"
)

// VideoCodec represents video codec options.
type VideoCodec string

const (
	VideoCodecH264 VideoCodec = "h264"
	VideoCodecH265 VideoCodec = "h265"
)

// VideoOutputFormat represents video output format.
type VideoOutputFormat string

const (
	VideoOutputHLS VideoOutputFormat = "hls"
	VideoOutputMP4 VideoOutputFormat = "mp4"
)

// TranscodingStatus represents the status of a transcoding job.
type TranscodingStatus string

const (
	TranscodingPending   TranscodingStatus = "pending"
	TranscodingQueued    TranscodingStatus = "queued"
	TranscodingProcess   TranscodingStatus = "processing"
	TranscodingCompleted TranscodingStatus = "completed"
	TranscodingFailed    TranscodingStatus = "failed"
	TranscodingCancelled TranscodingStatus = "cancelled"
)

// VideoVariant represents a video output variant.
type VideoVariant struct {
	Quality VideoQuality `json:"quality"`
	Codec   *VideoCodec  `json:"codec,omitempty"`
	Bitrate *int         `json:"bitrate,omitempty"`
}

// WatermarkOptions represents video watermark options.
type WatermarkOptions struct {
	Type         string  `json:"type"`
	ImageAssetID *string `json:"imageAssetId,omitempty"`
	Text         *string `json:"text,omitempty"`
	Position     string  `json:"position"`
	Opacity      *int    `json:"opacity,omitempty"`
}

// TrimOptions represents video trim options.
type TrimOptions struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

// TranscodeVideoRequest is the request for transcoding a video.
type TranscodeVideoRequest struct {
	ProjectSlug  string            `json:"projectSlug"`
	AssetID      string            `json:"assetId"`
	OutputFormat VideoOutputFormat `json:"outputFormat"`
	Variants     []VideoVariant    `json:"variants"`
	Watermark    *WatermarkOptions `json:"watermark,omitempty"`
	Trim         *TrimOptions      `json:"trim,omitempty"`
	WebhookURL   *string           `json:"webhookUrl,omitempty"`
}

// TranscodeJob represents a video transcoding job.
type TranscodeJob struct {
	ID              string            `json:"id"`
	AssetID         string            `json:"assetId"`
	Status          TranscodingStatus `json:"status"`
	OutputFormat    VideoOutputFormat `json:"outputFormat"`
	Variants        []VideoVariant    `json:"variants"`
	Progress        *int              `json:"progress,omitempty"`
	Error           *string           `json:"error,omitempty"`
	MediaConvertJob *string           `json:"mediaConvertJobId,omitempty"`
	CreatedAt       time.Time         `json:"createdAt"`
	StartedAt       *time.Time        `json:"startedAt,omitempty"`
	CompletedAt     *time.Time        `json:"completedAt,omitempty"`
}

// ListJobsRequest is the request for listing transcoding jobs.
type ListJobsRequest struct {
	ProjectSlug string             `json:"projectSlug"`
	AssetID     *string            `json:"assetId,omitempty"`
	Status      *TranscodingStatus `json:"status,omitempty"`
	Limit       *int               `json:"limit,omitempty"`
	Offset      *int               `json:"offset,omitempty"`
}

// ListJobsResponse is the response from listing transcoding jobs.
type ListJobsResponse struct {
	Jobs    []TranscodeJob `json:"jobs"`
	Total   int            `json:"total"`
	HasMore bool           `json:"hasMore"`
}

// MP4URL represents an MP4 URL with quality.
type MP4URL struct {
	Quality VideoQuality `json:"quality"`
	URL     string       `json:"url"`
}

// ThumbnailInfo represents a video thumbnail.
type ThumbnailInfo struct {
	URL       string `json:"url"`
	Timestamp int    `json:"timestamp"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

// StreamingURLs represents the streaming URLs for a video.
type StreamingURLs struct {
	HLSURL     *string         `json:"hlsUrl,omitempty"`
	MP4URLs    []MP4URL        `json:"mp4Urls"`
	Thumbnails []ThumbnailInfo `json:"thumbnails"`
}

// ThumbnailRequest is the request for generating a thumbnail.
type ThumbnailRequest struct {
	AssetID   string  `json:"assetId"`
	Timestamp float64 `json:"timestamp"`
	Width     *int    `json:"width,omitempty"`
	Format    *string `json:"format,omitempty"`
}

// ThumbnailResponse is the response from generating a thumbnail.
type ThumbnailResponse struct {
	URL       string `json:"url"`
	Timestamp int    `json:"timestamp"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

// RegenerateThumbnailRequest is the request for regenerating a thumbnail.
type RegenerateThumbnailRequest struct {
	AssetID   string  `json:"assetId"`
	Timestamp float64 `json:"timestamp"`
	Width     *int    `json:"width,omitempty"`
	Format    *string `json:"format,omitempty"`
}

// RegenerateThumbnailResponse is the response from regenerating a thumbnail.
type RegenerateThumbnailResponse struct {
	ID        *string `json:"id,omitempty"`
	AssetID   string  `json:"assetId"`
	Timestamp int     `json:"timestamp"`
	URL       *string `json:"url,omitempty"`
	Width     *int    `json:"width,omitempty"`
	Height    *int    `json:"height,omitempty"`
	Format    string  `json:"format"`
	Status    *string `json:"status,omitempty"`
}

// ExtractAudioRequest is the request for extracting audio.
type ExtractAudioRequest struct {
	ProjectSlug string  `json:"projectSlug"`
	AssetID     string  `json:"assetId"`
	Format      string  `json:"format"`
	Bitrate     *int    `json:"bitrate,omitempty"`
}

// ExtractAudioResponse is the response from extracting audio.
type ExtractAudioResponse struct {
	JobID  string            `json:"jobId"`
	Status TranscodingStatus `json:"status"`
}

// GifStatus represents the status of a GIF generation.
type GifStatus string

const (
	GifPending    GifStatus = "pending"
	GifProcessing GifStatus = "processing"
	GifCompleted  GifStatus = "completed"
	GifFailed     GifStatus = "failed"
)

// GenerateGifRequest is the request for generating a GIF.
type GenerateGifRequest struct {
	ProjectSlug     string   `json:"projectSlug"`
	AssetID         string   `json:"assetId"`
	StartTime       *float64 `json:"startTime,omitempty"`
	Duration        *float64 `json:"duration,omitempty"`
	Width           *int     `json:"width,omitempty"`
	FPS             *int     `json:"fps,omitempty"`
	OptimizePalette *bool    `json:"optimizePalette,omitempty"`
}

// VideoGif represents a generated GIF.
type VideoGif struct {
	ID           string     `json:"id"`
	AssetID      string     `json:"assetId"`
	StartTime    float64    `json:"startTime"`
	Duration     float64    `json:"duration"`
	FPS          int        `json:"fps"`
	URL          *string    `json:"url,omitempty"`
	Width        *int       `json:"width,omitempty"`
	Height       *int       `json:"height,omitempty"`
	SizeBytes    *int64     `json:"sizeBytes,omitempty"`
	FrameCount   *int       `json:"frameCount,omitempty"`
	Status       GifStatus  `json:"status"`
	ErrorMessage *string    `json:"errorMessage,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`
}

// ListGifsRequest is the request for listing GIFs.
type ListGifsRequest struct {
	AssetID string `json:"assetId"`
}

// VideoThumbnail represents a video thumbnail.
type VideoThumbnail struct {
	ID        string  `json:"id"`
	AssetID   string  `json:"assetId"`
	Timestamp int     `json:"timestamp"`
	URL       string  `json:"url"`
	Width     *int    `json:"width,omitempty"`
	Height    *int    `json:"height,omitempty"`
	Format    string  `json:"format"`
}

// ListThumbnailsResponse is the response from listing thumbnails.
type ListThumbnailsResponse struct {
	Thumbnails []VideoThumbnail `json:"thumbnails"`
}

// MergeOutputFormat represents merge output format.
type MergeOutputFormat string

const (
	MergeOutputMP4  MergeOutputFormat = "mp4"
	MergeOutputWebM MergeOutputFormat = "webm"
)

// TextOverlayShadow represents text overlay shadow config.
type TextOverlayShadow struct {
	Color   *string `json:"color,omitempty"`
	OffsetX *int    `json:"offsetX,omitempty"`
	OffsetY *int    `json:"offsetY,omitempty"`
}

// TextOverlayStroke represents text overlay stroke config.
type TextOverlayStroke struct {
	Color *string `json:"color,omitempty"`
	Width *int    `json:"width,omitempty"`
}

// TextOverlay represents text overlay config for videos.
type TextOverlay struct {
	Text            string             `json:"text"`
	Position        *string            `json:"position,omitempty"`
	FontSize        *int               `json:"fontSize,omitempty"`
	FontFamily      *string            `json:"fontFamily,omitempty"`
	FontWeight      *string            `json:"fontWeight,omitempty"`
	Color           *string            `json:"color,omitempty"`
	BackgroundColor *string            `json:"backgroundColor,omitempty"`
	Padding         *int               `json:"padding,omitempty"`
	MaxWidth        *int               `json:"maxWidth,omitempty"`
	Shadow          *TextOverlayShadow `json:"shadow,omitempty"`
	Stroke          *TextOverlayStroke `json:"stroke,omitempty"`
}

// MergeInputItem represents an input item for merge operations.
type MergeInputItem struct {
	AssetID     string       `json:"assetId"`
	Duration    *float64     `json:"duration,omitempty"`
	StartTime   *float64     `json:"startTime,omitempty"`
	EndTime     *float64     `json:"endTime,omitempty"`
	TextOverlay *TextOverlay `json:"textOverlay,omitempty"`
}

// AudioTrackInput represents an audio track for merge.
type AudioTrackInput struct {
	AssetID string   `json:"assetId"`
	Loop    bool     `json:"loop,omitempty"`
	FadeIn  *float64 `json:"fadeIn,omitempty"`
	FadeOut *float64 `json:"fadeOut,omitempty"`
}

// MergeOutputConfig represents output config for merge jobs.
type MergeOutputConfig struct {
	Format   *MergeOutputFormat `json:"format,omitempty"`
	Quality  *VideoQuality      `json:"quality,omitempty"`
	Filename *string            `json:"filename,omitempty"`
}

// CreateMergeJobRequest is the request for creating a merge job.
type CreateMergeJobRequest struct {
	ProjectSlug string             `json:"projectSlug"`
	Inputs      []MergeInputItem   `json:"inputs"`
	AudioTrack  *AudioTrackInput   `json:"audioTrack,omitempty"`
	Output      *MergeOutputConfig `json:"output,omitempty"`
	WebhookURL  *string            `json:"webhookUrl,omitempty"`
}

// MergeJob represents a video merge job.
type MergeJob struct {
	ID                   string            `json:"id"`
	OrganizationID       string            `json:"organizationId"`
	ProjectID            string            `json:"projectId"`
	Environment          string            `json:"environment"`
	Inputs               []MergeInputItem  `json:"inputs"`
	AudioTrackAssetID    *string           `json:"audioTrackAssetId,omitempty"`
	OutputFormat         MergeOutputFormat `json:"outputFormat"`
	OutputQuality        VideoQuality      `json:"outputQuality"`
	OutputFilename       *string           `json:"outputFilename,omitempty"`
	OutputAssetID        *string           `json:"outputAssetId,omitempty"`
	Status               TranscodingStatus `json:"status"`
	Progress             *int              `json:"progress,omitempty"`
	ErrorMessage         *string           `json:"errorMessage,omitempty"`
	StartedAt            *time.Time        `json:"startedAt,omitempty"`
	CompletedAt          *time.Time        `json:"completedAt,omitempty"`
	TotalDurationSeconds *float64          `json:"totalDurationSeconds,omitempty"`
	WebhookURL           *string           `json:"webhookUrl,omitempty"`
	CreatedAt            time.Time         `json:"createdAt"`
	UpdatedAt            *time.Time        `json:"updatedAt,omitempty"`
}

// MergeJobOutputAsset represents the output asset from a merge job.
type MergeJobOutputAsset struct {
	ID        string   `json:"id"`
	CDNURL    string   `json:"cdnUrl"`
	DirectURL string   `json:"directUrl"`
	Filename  string   `json:"filename"`
	Size      int64    `json:"size"`
	Duration  *float64 `json:"duration,omitempty"`
}

// MergeJobWithOutput represents a merge job with output details.
type MergeJobWithOutput struct {
	MergeJob
	OutputAsset *MergeJobOutputAsset `json:"outputAsset,omitempty"`
}

// ListMergeJobsRequest is the request for listing merge jobs.
type ListMergeJobsRequest struct {
	ProjectSlug string             `json:"projectSlug"`
	Status      *TranscodingStatus `json:"status,omitempty"`
	Limit       *int               `json:"limit,omitempty"`
	Offset      *int               `json:"offset,omitempty"`
}

// ListMergeJobsResponse is the response from listing merge jobs.
type ListMergeJobsResponse struct {
	Jobs    []MergeJob `json:"jobs"`
	Total   int        `json:"total"`
	HasMore bool       `json:"hasMore"`
}
