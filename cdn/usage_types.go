package cdn

import "time"

// CdnEnvironment represents the CDN environment.
type CdnEnvironment string

const (
	CdnEnvironmentSandbox    CdnEnvironment = "sandbox"
	CdnEnvironmentProduction CdnEnvironment = "production"
)

// CdnUsageRequest is the request for getting usage.
type CdnUsageRequest struct {
	ProjectSlug *string        `json:"projectSlug,omitempty"`
	Environment *CdnEnvironment `json:"environment,omitempty"`
	PeriodStart *time.Time     `json:"periodStart,omitempty"`
	PeriodEnd   *time.Time     `json:"periodEnd,omitempty"`
}

// CdnUsageResponse is the response from getting usage.
type CdnUsageResponse struct {
	PeriodStart           time.Time `json:"periodStart"`
	PeriodEnd             time.Time `json:"periodEnd"`
	Requests              int64     `json:"requests"`
	BandwidthBytes        int64     `json:"bandwidthBytes"`
	BandwidthFormatted    string    `json:"bandwidthFormatted"`
	Transformations       int64     `json:"transformations"`
	StorageBytes          int64     `json:"storageBytes"`
	StorageFormatted      string    `json:"storageFormatted"`
	EstimatedCostCents    int       `json:"estimatedCostCents"`
	EstimatedCostFormatted string   `json:"estimatedCostFormatted"`
}

// CdnUsageHistoryRequest is the request for getting usage history.
type CdnUsageHistoryRequest struct {
	ProjectSlug *string        `json:"projectSlug,omitempty"`
	Environment *CdnEnvironment `json:"environment,omitempty"`
	Days        *int           `json:"days,omitempty"`
	Granularity *string        `json:"granularity,omitempty"`
}

// CdnUsageDataPoint represents a single data point in usage history.
type CdnUsageDataPoint struct {
	Timestamp       time.Time `json:"timestamp"`
	Requests        int64     `json:"requests"`
	BandwidthBytes  int64     `json:"bandwidthBytes"`
	Transformations int64     `json:"transformations"`
}

// CdnUsageHistoryTotals represents totals in usage history.
type CdnUsageHistoryTotals struct {
	Requests        int64 `json:"requests"`
	BandwidthBytes  int64 `json:"bandwidthBytes"`
	Transformations int64 `json:"transformations"`
}

// CdnUsageHistoryResponse is the response from getting usage history.
type CdnUsageHistoryResponse struct {
	Data   []CdnUsageDataPoint   `json:"data"`
	Totals CdnUsageHistoryTotals `json:"totals"`
}

// CdnStorageBreakdownRequest is the request for getting storage breakdown.
type CdnStorageBreakdownRequest struct {
	ProjectSlug *string        `json:"projectSlug,omitempty"`
	Environment *CdnEnvironment `json:"environment,omitempty"`
	GroupBy     *string        `json:"groupBy,omitempty"`
}

// CdnStorageBreakdownItem represents a storage breakdown item.
type CdnStorageBreakdownItem struct {
	Key           string  `json:"key"`
	Count         int     `json:"count"`
	SizeBytes     int64   `json:"sizeBytes"`
	SizeFormatted string  `json:"sizeFormatted"`
	Percentage    float64 `json:"percentage"`
}

// CdnStorageBreakdownTotal represents storage breakdown totals.
type CdnStorageBreakdownTotal struct {
	Count         int    `json:"count"`
	SizeBytes     int64  `json:"sizeBytes"`
	SizeFormatted string `json:"sizeFormatted"`
}

// CdnStorageBreakdownResponse is the response from getting storage breakdown.
type CdnStorageBreakdownResponse struct {
	Items []CdnStorageBreakdownItem `json:"items"`
	Total CdnStorageBreakdownTotal  `json:"total"`
}
