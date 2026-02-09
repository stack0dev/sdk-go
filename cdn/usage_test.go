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

func TestClient_GetUsage(t *testing.T) {
	t.Run("without filters", func(t *testing.T) {
		cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/cdn/usage", r.URL.Path)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CdnUsageResponse{
				PeriodStart:         time.Now().AddDate(0, -1, 0),
				PeriodEnd:           time.Now(),
				StorageUsedBytes:    5000000000,
				BandwidthUsedBytes:  50000000000,
				RequestCount:        100000,
				TransformationCount: 50000,
			})
		})
		defer server.Close()

		resp, err := cdnClient.GetUsage(context.Background(), nil)

		require.NoError(t, err)
		assert.Equal(t, int64(5000000000), resp.StorageUsedBytes)
		assert.Equal(t, int64(50000000000), resp.BandwidthUsedBytes)
		assert.Equal(t, 100000, resp.RequestCount)
	})

	t.Run("with environment filter", func(t *testing.T) {
		cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "environment=production")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CdnUsageResponse{})
		})
		defer server.Close()

		env := types.EnvironmentProduction
		resp, err := cdnClient.GetUsage(context.Background(), &CdnUsageRequest{
			Environment: &env,
		})

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("with date range", func(t *testing.T) {
		cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "periodStart=")
			assert.Contains(t, r.URL.RawQuery, "periodEnd=")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CdnUsageResponse{})
		})
		defer server.Close()

		start := time.Now().AddDate(0, 0, -7)
		end := time.Now()
		resp, err := cdnClient.GetUsage(context.Background(), &CdnUsageRequest{
			PeriodStart: &start,
			PeriodEnd:   &end,
		})

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func TestClient_GetUsageHistory(t *testing.T) {
	t.Run("default request", func(t *testing.T) {
		cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Contains(t, r.URL.Path, "/cdn/usage/history")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CdnUsageHistoryResponse{
				Data: []CdnUsageHistoryDataPoint{
					{Date: "2024-01-01", StorageBytes: 1000000, BandwidthBytes: 5000000, Requests: 1000},
					{Date: "2024-01-02", StorageBytes: 1100000, BandwidthBytes: 5500000, Requests: 1100},
				},
			})
		})
		defer server.Close()

		resp, err := cdnClient.GetUsageHistory(context.Background(), nil)

		require.NoError(t, err)
		assert.Len(t, resp.Data, 2)
	})

	t.Run("with days parameter", func(t *testing.T) {
		cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "days=30")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CdnUsageHistoryResponse{Data: []CdnUsageHistoryDataPoint{}})
		})
		defer server.Close()

		days := 30
		resp, err := cdnClient.GetUsageHistory(context.Background(), &CdnUsageHistoryRequest{
			Days: &days,
		})

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("with granularity", func(t *testing.T) {
		cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "granularity=hourly")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CdnUsageHistoryResponse{Data: []CdnUsageHistoryDataPoint{}})
		})
		defer server.Close()

		granularity := "hourly"
		resp, err := cdnClient.GetUsageHistory(context.Background(), &CdnUsageHistoryRequest{
			Granularity: &granularity,
		})

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func TestClient_GetStorageBreakdown(t *testing.T) {
	t.Run("by type", func(t *testing.T) {
		cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Contains(t, r.URL.Path, "/cdn/usage/storage-breakdown")
			assert.Contains(t, r.URL.RawQuery, "groupBy=type")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CdnStorageBreakdownResponse{
				Breakdown: []CdnStorageBreakdownItem{
					{Category: "image", Size: 2000000000, Count: 500},
					{Category: "video", Size: 10000000000, Count: 50},
					{Category: "document", Size: 500000000, Count: 200},
				},
			})
		})
		defer server.Close()

		groupBy := "type"
		resp, err := cdnClient.GetStorageBreakdown(context.Background(), &CdnStorageBreakdownRequest{
			GroupBy: &groupBy,
		})

		require.NoError(t, err)
		assert.Len(t, resp.Breakdown, 3)
	})

	t.Run("by folder", func(t *testing.T) {
		cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "groupBy=folder")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CdnStorageBreakdownResponse{
				Breakdown: []CdnStorageBreakdownItem{
					{Category: "/images", Size: 3000000000, Count: 300},
					{Category: "/documents", Size: 1000000000, Count: 150},
				},
			})
		})
		defer server.Close()

		groupBy := "folder"
		resp, err := cdnClient.GetStorageBreakdown(context.Background(), &CdnStorageBreakdownRequest{
			GroupBy: &groupBy,
		})

		require.NoError(t, err)
		assert.Len(t, resp.Breakdown, 2)
	})

	t.Run("with project filter", func(t *testing.T) {
		cdnClient, server := setupCDNTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "projectSlug=my-project")

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CdnStorageBreakdownResponse{Breakdown: []CdnStorageBreakdownItem{}})
		})
		defer server.Close()

		projectSlug := "my-project"
		resp, err := cdnClient.GetStorageBreakdown(context.Background(), &CdnStorageBreakdownRequest{
			ProjectSlug: &projectSlug,
		})

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})
}
