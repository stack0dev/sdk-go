package cdn

import (
	"context"
	"net/url"
	"strconv"
)

// GetUsage gets current usage stats for the billing period.
func (c *Client) GetUsage(ctx context.Context, req *CdnUsageRequest) (*CdnUsageResponse, error) {
	params := url.Values{}
	if req != nil {
		if req.ProjectSlug != nil {
			params.Set("projectSlug", *req.ProjectSlug)
		}
		if req.Environment != nil {
			params.Set("environment", string(*req.Environment))
		}
		if req.PeriodStart != nil {
			params.Set("periodStart", req.PeriodStart.Format("2006-01-02T15:04:05Z07:00"))
		}
		if req.PeriodEnd != nil {
			params.Set("periodEnd", req.PeriodEnd.Format("2006-01-02T15:04:05Z07:00"))
		}
	}

	path := "/cdn/usage"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp CdnUsageResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetUsageHistory gets usage history (time series data for charts).
func (c *Client) GetUsageHistory(ctx context.Context, req *CdnUsageHistoryRequest) (*CdnUsageHistoryResponse, error) {
	params := url.Values{}
	if req != nil {
		if req.ProjectSlug != nil {
			params.Set("projectSlug", *req.ProjectSlug)
		}
		if req.Environment != nil {
			params.Set("environment", string(*req.Environment))
		}
		if req.Days != nil {
			params.Set("days", strconv.Itoa(*req.Days))
		}
		if req.Granularity != nil {
			params.Set("granularity", *req.Granularity)
		}
	}

	path := "/cdn/usage/history"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp CdnUsageHistoryResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetStorageBreakdown gets storage breakdown by type or folder.
func (c *Client) GetStorageBreakdown(ctx context.Context, req *CdnStorageBreakdownRequest) (*CdnStorageBreakdownResponse, error) {
	params := url.Values{}
	if req != nil {
		if req.ProjectSlug != nil {
			params.Set("projectSlug", *req.ProjectSlug)
		}
		if req.Environment != nil {
			params.Set("environment", string(*req.Environment))
		}
		if req.GroupBy != nil {
			params.Set("groupBy", *req.GroupBy)
		}
	}

	path := "/cdn/usage/storage-breakdown"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp CdnStorageBreakdownResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
