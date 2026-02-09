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

func setupDomainsTestClient(t *testing.T, handler http.HandlerFunc) (*DomainsClient, *httptest.Server) {
	server := httptest.NewServer(handler)
	httpClient := client.New("test-api-key", server.URL)
	return NewDomainsClient(httpClient), server
}

func TestDomainsClient_List(t *testing.T) {
	domainsClient, server := setupDomainsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/mail/domains")
		assert.Contains(t, r.URL.RawQuery, "projectSlug=my-project")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Domain{
			{
				ID:        "domain-123",
				Domain:    "example.com",
				Status:    DomainStatusVerified,
				IsDefault: true,
				CreatedAt: time.Now(),
			},
		})
	})
	defer server.Close()

	resp, err := domainsClient.List(context.Background(), &ListDomainsRequest{
		ProjectSlug: "my-project",
	})

	require.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "example.com", resp[0].Domain)
	assert.Equal(t, DomainStatusVerified, resp[0].Status)
}

func TestDomainsClient_List_WithEnvironment(t *testing.T) {
	domainsClient, server := setupDomainsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.RawQuery, "environment=production")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Domain{})
	})
	defer server.Close()

	env := types.EnvironmentProduction
	resp, err := domainsClient.List(context.Background(), &ListDomainsRequest{
		ProjectSlug: "my-project",
		Environment: &env,
	})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestDomainsClient_Add(t *testing.T) {
	domainsClient, server := setupDomainsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/domains", r.URL.Path)

		var req AddDomainRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "newdomain.com", req.Domain)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AddDomainResponse{
			Domain: &Domain{
				ID:     "domain-456",
				Domain: "newdomain.com",
				Status: DomainStatusPending,
			},
			DNSRecords: struct {
				Domain                string      `json:"domain"`
				DKIMRecords           []DNSRecord `json:"dkimRecords"`
				SPFRecord             DNSRecord   `json:"spfRecord"`
				DMARCRecord           DNSRecord   `json:"dmarcRecord"`
				VerificationToken     string      `json:"verificationToken"`
				SESVerificationRecord *DNSRecord  `json:"sesVerificationRecord,omitempty"`
			}{
				Domain: "newdomain.com",
				DKIMRecords: []DNSRecord{
					{Type: "CNAME", Name: "dkim._domainkey", Value: "dkim.example.com"},
				},
				SPFRecord: DNSRecord{Type: "TXT", Name: "@", Value: "v=spf1 include:amazonses.com ~all"},
				DMARCRecord: DNSRecord{Type: "TXT", Name: "_dmarc", Value: "v=DMARC1; p=none"},
				VerificationToken: "abc123",
			},
		})
	})
	defer server.Close()

	resp, err := domainsClient.Add(context.Background(), &AddDomainRequest{
		Domain: "newdomain.com",
	})

	require.NoError(t, err)
	assert.Equal(t, "domain-456", resp.Domain.ID)
	assert.Equal(t, DomainStatusPending, resp.Domain.Status)
	assert.Equal(t, "abc123", resp.DNSRecords.VerificationToken)
}

func TestDomainsClient_GetDNSRecords(t *testing.T) {
	domainID := "domain-123"
	domainsClient, server := setupDomainsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mail/domains/"+domainID+"/dns", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetDNSRecordsResponse{
			Domain: "example.com",
			DKIMRecords: []DNSRecord{
				{Type: "CNAME", Name: "dkim._domainkey", Value: "dkim.example.com"},
			},
			Status: DomainStatusVerified,
		})
	})
	defer server.Close()

	resp, err := domainsClient.GetDNSRecords(context.Background(), domainID)

	require.NoError(t, err)
	assert.Equal(t, "example.com", resp.Domain)
	assert.Len(t, resp.DKIMRecords, 1)
	assert.Equal(t, DomainStatusVerified, resp.Status)
}

func TestDomainsClient_Verify(t *testing.T) {
	domainID := "domain-123"
	domainsClient, server := setupDomainsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/domains/"+domainID+"/verify", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(VerifyDomainResponse{
			Verified: true,
			Message:  "Domain verified successfully",
		})
	})
	defer server.Close()

	resp, err := domainsClient.Verify(context.Background(), domainID)

	require.NoError(t, err)
	assert.True(t, resp.Verified)
	assert.Equal(t, "Domain verified successfully", resp.Message)
}

func TestDomainsClient_Verify_Failed(t *testing.T) {
	domainID := "domain-123"
	domainsClient, server := setupDomainsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(VerifyDomainResponse{
			Verified: false,
			Message:  "DNS records not found",
		})
	})
	defer server.Close()

	resp, err := domainsClient.Verify(context.Background(), domainID)

	require.NoError(t, err)
	assert.False(t, resp.Verified)
	assert.Contains(t, resp.Message, "DNS records not found")
}

func TestDomainsClient_Delete(t *testing.T) {
	domainID := "domain-123"
	domainsClient, server := setupDomainsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/mail/domains/"+domainID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DeleteDomainResponse{Success: true})
	})
	defer server.Close()

	resp, err := domainsClient.Delete(context.Background(), domainID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestDomainsClient_SetDefault(t *testing.T) {
	domainID := "domain-123"
	domainsClient, server := setupDomainsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/domains/"+domainID+"/default", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Domain{
			ID:        domainID,
			Domain:    "example.com",
			IsDefault: true,
		})
	})
	defer server.Close()

	resp, err := domainsClient.SetDefault(context.Background(), domainID)

	require.NoError(t, err)
	assert.True(t, resp.IsDefault)
}

func TestDomainStatus_Constants(t *testing.T) {
	assert.Equal(t, DomainStatus("pending"), DomainStatusPending)
	assert.Equal(t, DomainStatus("verified"), DomainStatusVerified)
	assert.Equal(t, DomainStatus("failed"), DomainStatusFailed)
}
