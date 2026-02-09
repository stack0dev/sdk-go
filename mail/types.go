// Package mail provides the mail client for the Stack0 SDK.
package mail

import (
	"time"

	"github.com/stack0/sdk-go/types"
)

// EmailAddress represents an email address with optional name.
type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// Attachment represents an email attachment.
type Attachment struct {
	Filename    string `json:"filename"`
	Content     string `json:"content"` // Base64 encoded
	ContentType string `json:"contentType,omitempty"`
	Path        string `json:"path,omitempty"` // URL to file
}

// EmailStatus represents the status of an email.
type EmailStatus string

const (
	EmailStatusPending      EmailStatus = "pending"
	EmailStatusSent         EmailStatus = "sent"
	EmailStatusDelivered    EmailStatus = "delivered"
	EmailStatusBounced      EmailStatus = "bounced"
	EmailStatusFailed       EmailStatus = "failed"
	EmailStatusDeferred     EmailStatus = "deferred"
	EmailStatusOpened       EmailStatus = "opened"
	EmailStatusClicked      EmailStatus = "clicked"
	EmailStatusComplained   EmailStatus = "complained"
	EmailStatusUnsubscribed EmailStatus = "unsubscribed"
)

// SendEmailRequest is the request to send an email.
type SendEmailRequest struct {
	ProjectSlug       *string                `json:"projectSlug,omitempty"`
	Environment       *types.Environment     `json:"environment,omitempty"`
	From              interface{}            `json:"from"` // string or EmailAddress
	To                interface{}            `json:"to"`   // string, EmailAddress, or []interface{}
	CC                interface{}            `json:"cc,omitempty"`
	BCC               interface{}            `json:"bcc,omitempty"`
	ReplyTo           interface{}            `json:"replyTo,omitempty"`
	Subject           string                 `json:"subject"`
	HTML              *string                `json:"html,omitempty"`
	Text              *string                `json:"text,omitempty"`
	TemplateID        *string                `json:"templateId,omitempty"`
	TemplateVariables map[string]interface{} `json:"templateVariables,omitempty"`
	Tags              []string               `json:"tags,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	Attachments       []Attachment           `json:"attachments,omitempty"`
	Headers           map[string]string      `json:"headers,omitempty"`
	ScheduledAt       *time.Time             `json:"scheduledAt,omitempty"`
}

// SendEmailResponse is the response after sending an email.
type SendEmailResponse struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Subject   string    `json:"subject"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// SendBatchEmailRequest is the request to send batch emails.
type SendBatchEmailRequest struct {
	ProjectSlug *string            `json:"projectSlug,omitempty"`
	Emails      []SendEmailRequest `json:"emails"`
}

// BatchEmailResult represents the result of a single email in a batch.
type BatchEmailResult struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// SendBatchEmailResponse is the response after sending batch emails.
type SendBatchEmailResponse struct {
	Success bool               `json:"success"`
	Data    []BatchEmailResult `json:"data"`
}

// SendBroadcastEmailRequest is the request to send a broadcast email.
type SendBroadcastEmailRequest struct {
	ProjectSlug       *string                `json:"projectSlug,omitempty"`
	Environment       *types.Environment     `json:"environment,omitempty"`
	From              interface{}            `json:"from"`
	To                []interface{}          `json:"to"`
	Subject           string                 `json:"subject"`
	HTML              *string                `json:"html,omitempty"`
	Text              *string                `json:"text,omitempty"`
	TemplateID        *string                `json:"templateId,omitempty"`
	TemplateVariables map[string]interface{} `json:"templateVariables,omitempty"`
	Tags              []string               `json:"tags,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	ScheduledAt       *time.Time             `json:"scheduledAt,omitempty"`
}

// SendBroadcastEmailResponse is the response after sending a broadcast email.
type SendBroadcastEmailResponse struct {
	Success        bool               `json:"success"`
	Data           []BatchEmailResult `json:"data"`
	Count          int                `json:"count"`
	TotalRequested *int               `json:"totalRequested,omitempty"`
	LimitedByQuota *bool              `json:"limitedByQuota,omitempty"`
}

// GetEmailResponse is the response when getting an email.
type GetEmailResponse struct {
	ID                string                 `json:"id"`
	From              string                 `json:"from"`
	To                string                 `json:"to"`
	Subject           string                 `json:"subject"`
	Status            string                 `json:"status"`
	HTML              *string                `json:"html"`
	Text              *string                `json:"text"`
	Tags              []string               `json:"tags"`
	Metadata          map[string]interface{} `json:"metadata"`
	CreatedAt         time.Time              `json:"createdAt"`
	SentAt            *time.Time             `json:"sentAt"`
	DeliveredAt       *time.Time             `json:"deliveredAt"`
	OpenedAt          *time.Time             `json:"openedAt"`
	ClickedAt         *time.Time             `json:"clickedAt"`
	BouncedAt         *time.Time             `json:"bouncedAt"`
	ProviderMessageID *string                `json:"providerMessageId"`
}

// ListEmailsRequest is the request to list emails.
type ListEmailsRequest struct {
	ProjectSlug *string            `url:"projectSlug,omitempty"`
	Environment *types.Environment `url:"environment,omitempty"`
	Limit       *int               `url:"limit,omitempty"`
	Offset      *int               `url:"offset,omitempty"`
	Status      *EmailStatus       `url:"status,omitempty"`
	From        *string            `url:"from,omitempty"`
	To          *string            `url:"to,omitempty"`
	Subject     *string            `url:"subject,omitempty"`
	Tag         *string            `url:"tag,omitempty"`
	StartDate   *time.Time         `url:"startDate,omitempty"`
	EndDate     *time.Time         `url:"endDate,omitempty"`
	SortBy      *string            `url:"sortBy,omitempty"`
	SortOrder   *string            `url:"sortOrder,omitempty"`
}

// Email represents an email in a list.
type Email struct {
	ID                string                 `json:"id"`
	From              string                 `json:"from"`
	To                string                 `json:"to"`
	Subject           string                 `json:"subject"`
	Status            string                 `json:"status"`
	CC                *string                `json:"cc"`
	ReplyTo           *string                `json:"replyTo"`
	MessageID         *string                `json:"messageId"`
	Tags              []string               `json:"tags"`
	Metadata          map[string]interface{} `json:"metadata"`
	CreatedAt         time.Time              `json:"createdAt"`
	DeliveredAt       *time.Time             `json:"deliveredAt"`
	OpenedAt          *time.Time             `json:"openedAt"`
	ClickedAt         *time.Time             `json:"clickedAt"`
	BouncedAt         *time.Time             `json:"bouncedAt"`
	ProviderMessageID *string                `json:"providerMessageId"`
}

// ListEmailsResponse is the response when listing emails.
type ListEmailsResponse struct {
	Emails []Email `json:"emails"`
	Total  int     `json:"total"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
}

// ResendEmailResponse is the response when resending an email.
type ResendEmailResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
	} `json:"data"`
}

// CancelEmailResponse is the response when canceling an email.
type CancelEmailResponse struct {
	Success bool `json:"success"`
}

// EmailAnalyticsResponse contains email analytics data.
type EmailAnalyticsResponse struct {
	Total        int     `json:"total"`
	Sent         int     `json:"sent"`
	Delivered    int     `json:"delivered"`
	Bounced      int     `json:"bounced"`
	Failed       int     `json:"failed"`
	DeliveryRate float64 `json:"deliveryRate"`
	OpenRate     float64 `json:"openRate"`
	ClickRate    float64 `json:"clickRate"`
}

// TimeSeriesAnalyticsRequest is the request for time series analytics.
type TimeSeriesAnalyticsRequest struct {
	Days *int `url:"days,omitempty"`
}

// TimeSeriesDataPoint represents a single data point in time series analytics.
type TimeSeriesDataPoint struct {
	Date      string `json:"date"`
	Sent      int    `json:"sent"`
	Delivered int    `json:"delivered"`
	Opened    int    `json:"opened"`
	Clicked   int    `json:"clicked"`
	Bounced   int    `json:"bounced"`
	Failed    int    `json:"failed"`
}

// TimeSeriesAnalyticsResponse contains time series analytics data.
type TimeSeriesAnalyticsResponse struct {
	Data []TimeSeriesDataPoint `json:"data"`
}

// HourlyAnalyticsDataPoint represents a single hourly data point.
type HourlyAnalyticsDataPoint struct {
	Hour      int `json:"hour"`
	Sent      int `json:"sent"`
	Delivered int `json:"delivered"`
	Opened    int `json:"opened"`
	Clicked   int `json:"clicked"`
}

// HourlyAnalyticsResponse contains hourly analytics data.
type HourlyAnalyticsResponse struct {
	Data []HourlyAnalyticsDataPoint `json:"data"`
}

// ListSendersRequest is the request to list senders.
type ListSendersRequest struct {
	ProjectSlug *string            `url:"projectSlug,omitempty"`
	Environment *types.Environment `url:"environment,omitempty"`
	Search      *string            `url:"search,omitempty"`
}

// Sender represents a unique sender with statistics.
type Sender struct {
	From      string `json:"from"`
	Total     int    `json:"total"`
	Sent      int    `json:"sent"`
	Delivered int    `json:"delivered"`
	Bounced   int    `json:"bounced"`
	Failed    int    `json:"failed"`
}

// ListSendersResponse is the response when listing senders.
type ListSendersResponse struct {
	Senders []Sender `json:"senders"`
}

// DomainStatus represents the verification status of a domain.
type DomainStatus string

const (
	DomainStatusPending  DomainStatus = "pending"
	DomainStatusVerified DomainStatus = "verified"
	DomainStatusFailed   DomainStatus = "failed"
)

// DNSRecord represents a DNS record for domain verification.
type DNSRecord struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Priority *int   `json:"priority,omitempty"`
}

// Domain represents a mail domain.
type Domain struct {
	ID                    string       `json:"id"`
	OrganizationID        string       `json:"organizationId"`
	Domain                string       `json:"domain"`
	Status                DomainStatus `json:"status"`
	DKIMRecord            []DNSRecord  `json:"dkimRecord"`
	SPFRecord             *DNSRecord   `json:"spfRecord"`
	DMARCRecord           *DNSRecord   `json:"dmarcRecord"`
	VerificationToken     *string      `json:"verificationToken"`
	SESVerificationRecord *DNSRecord   `json:"sesVerificationRecord"`
	IsDefault             bool         `json:"isDefault"`
	VerifiedAt            *time.Time   `json:"verifiedAt"`
	LastCheckedAt         *time.Time   `json:"lastCheckedAt"`
	CreatedAt             time.Time    `json:"createdAt"`
	UpdatedAt             *time.Time   `json:"updatedAt"`
}

// ListDomainsRequest is the request to list domains.
type ListDomainsRequest struct {
	ProjectSlug string             `url:"projectSlug"`
	Environment *types.Environment `url:"environment,omitempty"`
}

// AddDomainRequest is the request to add a domain.
type AddDomainRequest struct {
	Domain string `json:"domain"`
}

// AddDomainResponse is the response when adding a domain.
type AddDomainResponse struct {
	Domain     *Domain `json:"domain,omitempty"`
	DNSRecords struct {
		Domain                string      `json:"domain"`
		DKIMRecords           []DNSRecord `json:"dkimRecords"`
		SPFRecord             DNSRecord   `json:"spfRecord"`
		DMARCRecord           DNSRecord   `json:"dmarcRecord"`
		VerificationToken     string      `json:"verificationToken"`
		SESVerificationRecord *DNSRecord  `json:"sesVerificationRecord,omitempty"`
	} `json:"dnsRecords"`
}

// GetDNSRecordsResponse is the response when getting DNS records.
type GetDNSRecordsResponse struct {
	Domain                string       `json:"domain"`
	DKIMRecords           []DNSRecord  `json:"dkimRecords"`
	SPFRecord             *DNSRecord   `json:"spfRecord"`
	DMARCRecord           *DNSRecord   `json:"dmarcRecord"`
	SESVerificationRecord *DNSRecord   `json:"sesVerificationRecord"`
	Status                DomainStatus `json:"status"`
	VerifiedAt            *time.Time   `json:"verifiedAt"`
	VerificationDetails   *struct {
		DomainVerified     bool   `json:"domainVerified"`
		DKIMVerified       bool   `json:"dkimVerified"`
		VerificationStatus string `json:"verificationStatus"`
		DKIMStatus         string `json:"dkimStatus"`
	} `json:"verificationDetails,omitempty"`
}

// VerifyDomainResponse is the response when verifying a domain.
type VerifyDomainResponse struct {
	Verified bool   `json:"verified"`
	Message  string `json:"message"`
}

// DeleteDomainResponse is the response when deleting a domain.
type DeleteDomainResponse struct {
	Success bool `json:"success"`
}

// Template represents an email template.
type Template struct {
	ID              string                 `json:"id"`
	OrganizationID  string                 `json:"organizationId"`
	Environment     types.Environment      `json:"environment"`
	CreatedByUserID *string                `json:"createdByUserId"`
	Name            string                 `json:"name"`
	Slug            string                 `json:"slug"`
	Description     *string                `json:"description"`
	Subject         string                 `json:"subject"`
	PreviewText     *string                `json:"previewText"`
	HTML            string                 `json:"html"`
	Text            *string                `json:"text"`
	MailyJSON       map[string]interface{} `json:"mailyJson"`
	VariablesSchema map[string]interface{} `json:"variablesSchema"`
	IsActive        bool                   `json:"isActive"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       *time.Time             `json:"updatedAt"`
}

// CreateTemplateRequest is the request to create a template.
type CreateTemplateRequest struct {
	Environment     *types.Environment     `json:"environment,omitempty"`
	Name            string                 `json:"name"`
	Slug            string                 `json:"slug"`
	Description     *string                `json:"description,omitempty"`
	Subject         string                 `json:"subject"`
	PreviewText     *string                `json:"previewText,omitempty"`
	HTML            string                 `json:"html"`
	Text            *string                `json:"text,omitempty"`
	MailyJSON       map[string]interface{} `json:"mailyJson,omitempty"`
	VariablesSchema map[string]interface{} `json:"variablesSchema,omitempty"`
	IsActive        *bool                  `json:"isActive,omitempty"`
}

// UpdateTemplateRequest is the request to update a template.
type UpdateTemplateRequest struct {
	ID              string
	Name            *string                `json:"name,omitempty"`
	Slug            *string                `json:"slug,omitempty"`
	Description     *string                `json:"description,omitempty"`
	Subject         *string                `json:"subject,omitempty"`
	PreviewText     *string                `json:"previewText,omitempty"`
	HTML            *string                `json:"html,omitempty"`
	Text            *string                `json:"text,omitempty"`
	MailyJSON       map[string]interface{} `json:"mailyJson,omitempty"`
	VariablesSchema map[string]interface{} `json:"variablesSchema,omitempty"`
	IsActive        *bool                  `json:"isActive,omitempty"`
}

// ListTemplatesRequest is the request to list templates.
type ListTemplatesRequest struct {
	Environment *types.Environment `url:"environment,omitempty"`
	Limit       *int               `url:"limit,omitempty"`
	Offset      *int               `url:"offset,omitempty"`
	IsActive    *bool              `url:"isActive,omitempty"`
	Search      *string            `url:"search,omitempty"`
}

// ListTemplatesResponse is the response when listing templates.
type ListTemplatesResponse struct {
	Templates []Template `json:"templates"`
	Total     int        `json:"total"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

// DeleteTemplateResponse is the response when deleting a template.
type DeleteTemplateResponse struct {
	Success bool `json:"success"`
}

// PreviewTemplateRequest is the request to preview a template.
type PreviewTemplateRequest struct {
	ID        string
	Variables map[string]interface{} `json:"variables"`
}

// PreviewTemplateResponse is the response when previewing a template.
type PreviewTemplateResponse struct {
	Subject string  `json:"subject"`
	HTML    string  `json:"html"`
	Text    *string `json:"text"`
}

// Audience represents a contact audience.
type Audience struct {
	ID                   string     `json:"id"`
	OrganizationID       string     `json:"organizationId"`
	ProjectID            *string    `json:"projectId"`
	Environment          string     `json:"environment"`
	Name                 string     `json:"name"`
	Description          *string    `json:"description"`
	TotalContacts        int        `json:"totalContacts"`
	SubscribedContacts   int        `json:"subscribedContacts"`
	UnsubscribedContacts int        `json:"unsubscribedContacts"`
	CreatedByUserID      *string    `json:"createdByUserId"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            *time.Time `json:"updatedAt"`
}

// CreateAudienceRequest is the request to create an audience.
type CreateAudienceRequest struct {
	Environment *types.Environment `json:"environment,omitempty"`
	Name        string             `json:"name"`
	Description *string            `json:"description,omitempty"`
}

// UpdateAudienceRequest is the request to update an audience.
type UpdateAudienceRequest struct {
	ID          string
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ListAudiencesRequest is the request to list audiences.
type ListAudiencesRequest struct {
	Environment *types.Environment `url:"environment,omitempty"`
	Limit       *int               `url:"limit,omitempty"`
	Offset      *int               `url:"offset,omitempty"`
	Search      *string            `url:"search,omitempty"`
}

// ListAudiencesResponse is the response when listing audiences.
type ListAudiencesResponse struct {
	Audiences []Audience `json:"audiences"`
	Total     int        `json:"total"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

// DeleteAudienceResponse is the response when deleting an audience.
type DeleteAudienceResponse struct {
	Success bool `json:"success"`
}

// AddContactsToAudienceRequest is the request to add contacts to an audience.
type AddContactsToAudienceRequest struct {
	ID         string
	ContactIDs []string `json:"contactIds"`
}

// AddContactsToAudienceResponse is the response when adding contacts to an audience.
type AddContactsToAudienceResponse struct {
	Success bool `json:"success"`
	Added   int  `json:"added"`
}

// RemoveContactsFromAudienceRequest is the request to remove contacts from an audience.
type RemoveContactsFromAudienceRequest struct {
	ID         string
	ContactIDs []string `json:"contactIds"`
}

// RemoveContactsFromAudienceResponse is the response when removing contacts from an audience.
type RemoveContactsFromAudienceResponse struct {
	Success bool `json:"success"`
	Removed int  `json:"removed"`
}

// ContactStatus represents the status of a contact.
type ContactStatus string

const (
	ContactStatusSubscribed   ContactStatus = "subscribed"
	ContactStatusUnsubscribed ContactStatus = "unsubscribed"
	ContactStatusBounced      ContactStatus = "bounced"
	ContactStatusComplained   ContactStatus = "complained"
)

// MailContact represents a mail contact.
type MailContact struct {
	ID             string                 `json:"id"`
	OrganizationID string                 `json:"organizationId"`
	ProjectID      *string                `json:"projectId"`
	Environment    string                 `json:"environment"`
	Email          string                 `json:"email"`
	FirstName      *string                `json:"firstName"`
	LastName       *string                `json:"lastName"`
	Metadata       map[string]interface{} `json:"metadata"`
	Status         string                 `json:"status"`
	SubscribedAt   *time.Time             `json:"subscribedAt"`
	UnsubscribedAt *time.Time             `json:"unsubscribedAt"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      *time.Time             `json:"updatedAt"`
}

// AudienceContact represents a contact within an audience.
type AudienceContact struct {
	MailContact
	AddedAt *time.Time `json:"addedAt"`
}

// CreateContactRequest is the request to create a contact.
type CreateContactRequest struct {
	Environment *types.Environment     `json:"environment,omitempty"`
	Email       string                 `json:"email"`
	FirstName   *string                `json:"firstName,omitempty"`
	LastName    *string                `json:"lastName,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateContactRequest is the request to update a contact.
type UpdateContactRequest struct {
	ID        string
	Email     *string                `json:"email,omitempty"`
	FirstName *string                `json:"firstName,omitempty"`
	LastName  *string                `json:"lastName,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Status    *ContactStatus         `json:"status,omitempty"`
}

// ListContactsRequest is the request to list contacts.
type ListContactsRequest struct {
	Environment *types.Environment `url:"environment,omitempty"`
	Limit       *int               `url:"limit,omitempty"`
	Offset      *int               `url:"offset,omitempty"`
	Search      *string            `url:"search,omitempty"`
	Status      *ContactStatus     `url:"status,omitempty"`
}

// ListContactsResponse is the response when listing contacts.
type ListContactsResponse struct {
	Contacts []MailContact `json:"contacts"`
	Total    int           `json:"total"`
	Limit    int           `json:"limit"`
	Offset   int           `json:"offset"`
}

// ListAudienceContactsRequest is the request to list contacts in an audience.
type ListAudienceContactsRequest struct {
	ID          string
	Environment *types.Environment `url:"environment,omitempty"`
	Limit       *int               `url:"limit,omitempty"`
	Offset      *int               `url:"offset,omitempty"`
	Search      *string            `url:"search,omitempty"`
	Status      *ContactStatus     `url:"status,omitempty"`
}

// ListAudienceContactsResponse is the response when listing audience contacts.
type ListAudienceContactsResponse struct {
	Contacts []AudienceContact `json:"contacts"`
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

// DeleteContactResponse is the response when deleting a contact.
type DeleteContactResponse struct {
	Success bool `json:"success"`
}

// ImportContactInput represents a single contact to import.
type ImportContactInput struct {
	Email     string                 `json:"email"`
	FirstName *string                `json:"firstName,omitempty"`
	LastName  *string                `json:"lastName,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ImportContactsRequest is the request to import contacts.
type ImportContactsRequest struct {
	Environment *types.Environment   `json:"environment,omitempty"`
	AudienceID  *string              `json:"audienceId,omitempty"`
	Contacts    []ImportContactInput `json:"contacts"`
}

// ImportContactError represents an error during contact import.
type ImportContactError struct {
	Email string `json:"email"`
	Error string `json:"error"`
}

// ImportContactsResponse is the response when importing contacts.
type ImportContactsResponse struct {
	Success  bool                 `json:"success"`
	Imported int                  `json:"imported"`
	Skipped  int                  `json:"skipped"`
	Errors   []ImportContactError `json:"errors"`
}

// CampaignStatus represents the status of a campaign.
type CampaignStatus string

const (
	CampaignStatusDraft     CampaignStatus = "draft"
	CampaignStatusScheduled CampaignStatus = "scheduled"
	CampaignStatusSending   CampaignStatus = "sending"
	CampaignStatusSent      CampaignStatus = "sent"
	CampaignStatusPaused    CampaignStatus = "paused"
	CampaignStatusCancelled CampaignStatus = "cancelled"
	CampaignStatusFailed    CampaignStatus = "failed"
)

// Campaign represents an email campaign.
type Campaign struct {
	ID              string                 `json:"id"`
	OrganizationID  string                 `json:"organizationId"`
	ProjectID       *string                `json:"projectId"`
	Environment     string                 `json:"environment"`
	Name            string                 `json:"name"`
	Subject         string                 `json:"subject"`
	PreviewText     *string                `json:"previewText"`
	FromEmail       string                 `json:"fromEmail"`
	FromName        *string                `json:"fromName"`
	ReplyTo         *string                `json:"replyTo"`
	TemplateID      *string                `json:"templateId"`
	HTML            *string                `json:"html"`
	Text            *string                `json:"text"`
	AudienceID      *string                `json:"audienceId"`
	Status          string                 `json:"status"`
	ScheduledAt     *time.Time             `json:"scheduledAt"`
	SentAt          *time.Time             `json:"sentAt"`
	CompletedAt     *time.Time             `json:"completedAt"`
	TotalRecipients int                    `json:"totalRecipients"`
	SentCount       int                    `json:"sentCount"`
	DeliveredCount  int                    `json:"deliveredCount"`
	OpenedCount     int                    `json:"openedCount"`
	ClickedCount    int                    `json:"clickedCount"`
	BouncedCount    int                    `json:"bouncedCount"`
	FailedCount     int                    `json:"failedCount"`
	Tags            []string               `json:"tags"`
	Metadata        map[string]interface{} `json:"metadata"`
	CreatedByUserID *string                `json:"createdByUserId"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       *time.Time             `json:"updatedAt"`
}

// CreateCampaignRequest is the request to create a campaign.
type CreateCampaignRequest struct {
	Environment *types.Environment `json:"environment,omitempty"`
	Name        string             `json:"name"`
	Subject     string             `json:"subject"`
	PreviewText *string            `json:"previewText,omitempty"`
	FromEmail   string             `json:"fromEmail"`
	FromName    *string            `json:"fromName,omitempty"`
	ReplyTo     *string            `json:"replyTo,omitempty"`
	TemplateID  *string            `json:"templateId,omitempty"`
	HTML        *string            `json:"html,omitempty"`
	Text        *string            `json:"text,omitempty"`
	AudienceID  *string            `json:"audienceId,omitempty"`
	ScheduledAt *time.Time         `json:"scheduledAt,omitempty"`
	Tags        []string           `json:"tags,omitempty"`
}

// UpdateCampaignRequest is the request to update a campaign.
type UpdateCampaignRequest struct {
	ID          string
	Name        *string    `json:"name,omitempty"`
	Subject     *string    `json:"subject,omitempty"`
	PreviewText *string    `json:"previewText,omitempty"`
	FromEmail   *string    `json:"fromEmail,omitempty"`
	FromName    *string    `json:"fromName,omitempty"`
	ReplyTo     *string    `json:"replyTo,omitempty"`
	TemplateID  *string    `json:"templateId,omitempty"`
	HTML        *string    `json:"html,omitempty"`
	Text        *string    `json:"text,omitempty"`
	AudienceID  *string    `json:"audienceId,omitempty"`
	ScheduledAt *time.Time `json:"scheduledAt,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
}

// ListCampaignsRequest is the request to list campaigns.
type ListCampaignsRequest struct {
	Environment *types.Environment `url:"environment,omitempty"`
	Limit       *int               `url:"limit,omitempty"`
	Offset      *int               `url:"offset,omitempty"`
	Search      *string            `url:"search,omitempty"`
	Status      *CampaignStatus    `url:"status,omitempty"`
}

// ListCampaignsResponse is the response when listing campaigns.
type ListCampaignsResponse struct {
	Campaigns []Campaign `json:"campaigns"`
	Total     int        `json:"total"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

// DeleteCampaignResponse is the response when deleting a campaign.
type DeleteCampaignResponse struct {
	Success bool `json:"success"`
}

// SendCampaignRequest is the request to send a campaign.
type SendCampaignRequest struct {
	ID          string
	SendNow     *bool      `json:"sendNow,omitempty"`
	ScheduledAt *time.Time `json:"scheduledAt,omitempty"`
}

// SendCampaignResponse is the response when sending a campaign.
type SendCampaignResponse struct {
	Success         bool `json:"success"`
	SentCount       int  `json:"sentCount"`
	FailedCount     int  `json:"failedCount"`
	TotalRecipients int  `json:"totalRecipients"`
}

// PauseCampaignResponse is the response when pausing a campaign.
type PauseCampaignResponse struct {
	Success bool `json:"success"`
}

// CancelCampaignResponse is the response when canceling a campaign.
type CancelCampaignResponse struct {
	Success bool `json:"success"`
}

// CampaignStatsResponse contains campaign statistics.
type CampaignStatsResponse struct {
	Total        int     `json:"total"`
	Sent         int     `json:"sent"`
	Delivered    int     `json:"delivered"`
	Opened       int     `json:"opened"`
	Clicked      int     `json:"clicked"`
	Bounced      int     `json:"bounced"`
	Failed       int     `json:"failed"`
	DeliveryRate float64 `json:"deliveryRate"`
	OpenRate     float64 `json:"openRate"`
	ClickRate    float64 `json:"clickRate"`
	BounceRate   float64 `json:"bounceRate"`
}

// SequenceStatus represents the status of a sequence.
type SequenceStatus string

const (
	SequenceStatusDraft    SequenceStatus = "draft"
	SequenceStatusActive   SequenceStatus = "active"
	SequenceStatusPaused   SequenceStatus = "paused"
	SequenceStatusArchived SequenceStatus = "archived"
)

// SequenceTriggerType represents the type of trigger for a sequence.
type SequenceTriggerType string

const (
	SequenceTriggerManual       SequenceTriggerType = "manual"
	SequenceTriggerEventReceived SequenceTriggerType = "event_received"
	SequenceTriggerContactAdded  SequenceTriggerType = "contact_added"
	SequenceTriggerAPI          SequenceTriggerType = "api"
	SequenceTriggerScheduled    SequenceTriggerType = "scheduled"
)

// SequenceTriggerFrequency represents how often a trigger can fire.
type SequenceTriggerFrequency string

const (
	SequenceTriggerOnce   SequenceTriggerFrequency = "once"
	SequenceTriggerAlways SequenceTriggerFrequency = "always"
)

// SequenceNodeType represents the type of a sequence node.
type SequenceNodeType string

const (
	SequenceNodeTrigger       SequenceNodeType = "trigger"
	SequenceNodeEmail         SequenceNodeType = "email"
	SequenceNodeTimer         SequenceNodeType = "timer"
	SequenceNodeFilter        SequenceNodeType = "filter"
	SequenceNodeBranch        SequenceNodeType = "branch"
	SequenceNodeExperiment    SequenceNodeType = "experiment"
	SequenceNodeExit          SequenceNodeType = "exit"
	SequenceNodeAddToList     SequenceNodeType = "add_to_list"
	SequenceNodeUpdateContact SequenceNodeType = "update_contact"
)

// ConnectionType represents the type of connection between nodes.
type ConnectionType string

const (
	ConnectionDefault ConnectionType = "default"
	ConnectionYes     ConnectionType = "yes"
	ConnectionNo      ConnectionType = "no"
	ConnectionBranch  ConnectionType = "branch"
	ConnectionVariant ConnectionType = "variant"
)

// Sequence represents an email sequence.
type Sequence struct {
	ID               string                   `json:"id"`
	OrganizationID   string                   `json:"organizationId"`
	Environment      string                   `json:"environment"`
	Name             string                   `json:"name"`
	Description      *string                  `json:"description"`
	TriggerType      SequenceTriggerType      `json:"triggerType"`
	TriggerFrequency SequenceTriggerFrequency `json:"triggerFrequency"`
	TriggerConfig    map[string]interface{}   `json:"triggerConfig"`
	AudienceFilterID *string                  `json:"audienceFilterId"`
	Status           SequenceStatus           `json:"status"`
	TotalEntered     int                      `json:"totalEntered"`
	TotalCompleted   int                      `json:"totalCompleted"`
	TotalActive      int                      `json:"totalActive"`
	PublishedAt      *time.Time               `json:"publishedAt"`
	PausedAt         *time.Time               `json:"pausedAt"`
	ArchivedAt       *time.Time               `json:"archivedAt"`
	CreatedByUserID  *string                  `json:"createdByUserId"`
	CreatedAt        time.Time                `json:"createdAt"`
	UpdatedAt        *time.Time               `json:"updatedAt"`
}

// SequenceNode represents a node in a sequence.
type SequenceNode struct {
	ID        string                 `json:"id"`
	LoopID    string                 `json:"loopId"`
	NodeType  SequenceNodeType       `json:"nodeType"`
	Name      string                 `json:"name"`
	PositionX float64                `json:"positionX"`
	PositionY float64                `json:"positionY"`
	SortOrder int                    `json:"sortOrder"`
	Config    map[string]interface{} `json:"config"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt *time.Time             `json:"updatedAt"`
}

// SequenceConnection represents a connection between nodes.
type SequenceConnection struct {
	ID             string         `json:"id"`
	LoopID         string         `json:"loopId"`
	SourceNodeID   string         `json:"sourceNodeId"`
	TargetNodeID   string         `json:"targetNodeId"`
	ConnectionType ConnectionType `json:"connectionType"`
	Label          *string        `json:"label"`
	CreatedAt      time.Time      `json:"createdAt"`
}

// SequenceWithNodes represents a sequence with its nodes and connections.
type SequenceWithNodes struct {
	Sequence
	Nodes       []SequenceNode       `json:"nodes"`
	Connections []SequenceConnection `json:"connections"`
}

// CreateSequenceRequest is the request to create a sequence.
type CreateSequenceRequest struct {
	Environment      *types.Environment       `json:"environment,omitempty"`
	Name             string                   `json:"name"`
	Description      *string                  `json:"description,omitempty"`
	TriggerType      SequenceTriggerType      `json:"triggerType"`
	TriggerFrequency *SequenceTriggerFrequency `json:"triggerFrequency,omitempty"`
	TriggerConfig    map[string]interface{}   `json:"triggerConfig,omitempty"`
	AudienceFilterID *string                  `json:"audienceFilterId,omitempty"`
}

// UpdateSequenceRequest is the request to update a sequence.
type UpdateSequenceRequest struct {
	ID               string
	Name             *string                   `json:"name,omitempty"`
	Description      *string                   `json:"description,omitempty"`
	TriggerType      *SequenceTriggerType      `json:"triggerType,omitempty"`
	TriggerFrequency *SequenceTriggerFrequency `json:"triggerFrequency,omitempty"`
	TriggerConfig    map[string]interface{}   `json:"triggerConfig,omitempty"`
	AudienceFilterID *string                  `json:"audienceFilterId,omitempty"`
}

// ListSequencesRequest is the request to list sequences.
type ListSequencesRequest struct {
	Environment *types.Environment   `url:"environment,omitempty"`
	Limit       *int                 `url:"limit,omitempty"`
	Offset      *int                 `url:"offset,omitempty"`
	Search      *string              `url:"search,omitempty"`
	Status      *SequenceStatus      `url:"status,omitempty"`
	TriggerType *SequenceTriggerType `url:"triggerType,omitempty"`
}

// ListSequencesResponse is the response when listing sequences.
type ListSequencesResponse struct {
	Sequences []Sequence `json:"sequences"`
	Total     int        `json:"total"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

// DeleteSequenceResponse is the response when deleting a sequence.
type DeleteSequenceResponse struct {
	Success bool `json:"success"`
}

// PublishSequenceResponse is the response when publishing a sequence.
type PublishSequenceResponse struct {
	Success bool `json:"success"`
}

// PauseSequenceResponse is the response when pausing a sequence.
type PauseSequenceResponse struct {
	Success bool `json:"success"`
}

// ResumeSequenceResponse is the response when resuming a sequence.
type ResumeSequenceResponse struct {
	Success bool `json:"success"`
}

// ArchiveSequenceResponse is the response when archiving a sequence.
type ArchiveSequenceResponse struct {
	Success bool `json:"success"`
}

// CreateNodeRequest is the request to create a node.
type CreateNodeRequest struct {
	ID        string // sequence ID
	NodeType  SequenceNodeType       `json:"nodeType"`
	Name      string                 `json:"name"`
	PositionX float64                `json:"positionX"`
	PositionY float64                `json:"positionY"`
	SortOrder *int                   `json:"sortOrder,omitempty"`
	Config    map[string]interface{} `json:"config,omitempty"`
}

// UpdateNodeRequest is the request to update a node.
type UpdateNodeRequest struct {
	ID        string // sequence ID
	NodeID    string
	Name      *string                `json:"name,omitempty"`
	PositionX *float64               `json:"positionX,omitempty"`
	PositionY *float64               `json:"positionY,omitempty"`
	SortOrder *int                   `json:"sortOrder,omitempty"`
	Config    map[string]interface{} `json:"config,omitempty"`
}

// UpdateNodePositionRequest is the request to update a node's position.
type UpdateNodePositionRequest struct {
	ID        string // sequence ID
	NodeID    string
	PositionX float64 `json:"positionX"`
	PositionY float64 `json:"positionY"`
}

// DeleteNodeResponse is the response when deleting a node.
type DeleteNodeResponse struct {
	Success bool `json:"success"`
}

// SetNodeEmailRequest is the request to set email content for a node.
type SetNodeEmailRequest struct {
	NodeID      string
	Subject     *string                `json:"subject,omitempty"`
	PreviewText *string                `json:"previewText,omitempty"`
	HTML        *string                `json:"html,omitempty"`
	Text        *string                `json:"text,omitempty"`
	TemplateID  *string                `json:"templateId,omitempty"`
	MailyJSON   map[string]interface{} `json:"mailyJson,omitempty"`
	FromEmail   *string                `json:"fromEmail,omitempty"`
	FromName    *string                `json:"fromName,omitempty"`
	ReplyTo     *string                `json:"replyTo,omitempty"`
}

// SetNodeTimerRequest is the request to set timer configuration for a node.
type SetNodeTimerRequest struct {
	NodeID            string
	DelayAmount       int    `json:"delayAmount"`
	DelayUnit         string `json:"delayUnit"` // minutes, hours, days, weeks
	WaitUntilTime     *string `json:"waitUntilTime,omitempty"`
	WaitUntilTimezone *string `json:"waitUntilTimezone,omitempty"`
}

// SetNodeFilterRequest is the request to set filter configuration for a node.
type SetNodeFilterRequest struct {
	NodeID         string
	Conditions     map[string]interface{} `json:"conditions"`
	NonMatchAction *string                `json:"nonMatchAction,omitempty"` // stop, continue
}

// BranchCondition represents a branch condition.
type BranchCondition struct {
	Name       string                 `json:"name"`
	Conditions map[string]interface{} `json:"conditions"`
}

// SetNodeBranchRequest is the request to set branch configuration for a node.
type SetNodeBranchRequest struct {
	NodeID           string
	Branches         []BranchCondition `json:"branches"`
	HasDefaultBranch *bool             `json:"hasDefaultBranch,omitempty"`
}

// ExperimentVariant represents an experiment variant.
type ExperimentVariant struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
}

// SetNodeExperimentRequest is the request to set experiment configuration for a node.
type SetNodeExperimentRequest struct {
	NodeID     string
	SampleSize *int                `json:"sampleSize,omitempty"`
	Variants   []ExperimentVariant `json:"variants"`
}

// CreateConnectionRequest is the request to create a connection.
type CreateConnectionRequest struct {
	ID             string // sequence ID
	SourceNodeID   string          `json:"sourceNodeId"`
	TargetNodeID   string          `json:"targetNodeId"`
	ConnectionType *ConnectionType `json:"connectionType,omitempty"`
	Label          *string         `json:"label,omitempty"`
}

// DeleteConnectionResponse is the response when deleting a connection.
type DeleteConnectionResponse struct {
	Success bool `json:"success"`
}

// SequenceEntryStatus represents the status of an entry in a sequence.
type SequenceEntryStatus string

const (
	SequenceEntryStatusActive    SequenceEntryStatus = "active"
	SequenceEntryStatusPaused    SequenceEntryStatus = "paused"
	SequenceEntryStatusCompleted SequenceEntryStatus = "completed"
	SequenceEntryStatusStopped   SequenceEntryStatus = "stopped"
	SequenceEntryStatusFailed    SequenceEntryStatus = "failed"
)

// SequenceEntry represents a contact's entry in a sequence.
type SequenceEntry struct {
	ID            string              `json:"id"`
	LoopID        string              `json:"loopId"`
	ContactID     string              `json:"contactId"`
	CurrentNodeID *string             `json:"currentNodeId"`
	Status        SequenceEntryStatus `json:"status"`
	EnteredAt     time.Time           `json:"enteredAt"`
	ExitedAt      *time.Time          `json:"exitedAt"`
	ExitReason    *string             `json:"exitReason"`
	Contact       *MailContact        `json:"contact,omitempty"`
}

// ListSequenceEntriesRequest is the request to list sequence entries.
type ListSequenceEntriesRequest struct {
	ID     string // sequence ID
	Status *SequenceEntryStatus `url:"status,omitempty"`
	Limit  *int                 `url:"limit,omitempty"`
	Offset *int                 `url:"offset,omitempty"`
}

// ListSequenceEntriesResponse is the response when listing sequence entries.
type ListSequenceEntriesResponse struct {
	Entries []SequenceEntry `json:"entries"`
	Total   int             `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
}

// AddContactToSequenceRequest is the request to add a contact to a sequence.
type AddContactToSequenceRequest struct {
	ID        string // sequence ID
	ContactID string `json:"contactId"`
}

// RemoveContactFromSequenceRequest is the request to remove a contact from a sequence.
type RemoveContactFromSequenceRequest struct {
	ID      string // sequence ID
	EntryID string `json:"entryId"`
	Reason  *string `json:"reason,omitempty"`
}

// RemoveContactFromSequenceResponse is the response when removing a contact from a sequence.
type RemoveContactFromSequenceResponse struct {
	Success bool `json:"success"`
}

// SequenceAnalyticsResponse contains sequence analytics.
type SequenceAnalyticsResponse struct {
	Sequence struct {
		TotalEntered   int `json:"totalEntered"`
		TotalCompleted int `json:"totalCompleted"`
		TotalActive    int `json:"totalActive"`
	} `json:"sequence"`
	StatusBreakdown map[string]int `json:"statusBreakdown"`
	NodeAnalytics   []struct {
		NodeID          string `json:"nodeId"`
		Entered         int    `json:"entered"`
		Exited          int    `json:"exited"`
		EmailsSent      int    `json:"emailsSent"`
		EmailsDelivered int    `json:"emailsDelivered"`
		EmailsOpened    int    `json:"emailsOpened"`
		EmailsClicked   int    `json:"emailsClicked"`
		EmailsBounced   int    `json:"emailsBounced"`
		Passed          int    `json:"passed"`
		Filtered        int    `json:"filtered"`
	} `json:"nodeAnalytics"`
}

// EventProperty represents a property in an event schema.
type EventProperty struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"` // string, number, boolean, date, object, array
	Description *string `json:"description,omitempty"`
	Required    *bool   `json:"required,omitempty"`
}

// EventPropertiesSchema represents the schema for event properties.
type EventPropertiesSchema struct {
	Properties []EventProperty `json:"properties"`
}

// MailEvent represents an event definition.
type MailEvent struct {
	ID               string                 `json:"id"`
	OrganizationID   string                 `json:"organizationId"`
	ProjectID        *string                `json:"projectId"`
	Environment      string                 `json:"environment"`
	Name             string                 `json:"name"`
	Description      *string                `json:"description"`
	PropertiesSchema *EventPropertiesSchema `json:"propertiesSchema"`
	TotalReceived    int                    `json:"totalReceived"`
	LastReceivedAt   *time.Time             `json:"lastReceivedAt"`
	CreatedAt        time.Time              `json:"createdAt"`
	UpdatedAt        *time.Time             `json:"updatedAt"`
}

// CreateEventRequest is the request to create an event.
type CreateEventRequest struct {
	ProjectSlug      *string                `json:"projectSlug,omitempty"`
	Environment      *types.Environment     `json:"environment,omitempty"`
	Name             string                 `json:"name"`
	Description      *string                `json:"description,omitempty"`
	PropertiesSchema *EventPropertiesSchema `json:"propertiesSchema,omitempty"`
}

// UpdateEventRequest is the request to update an event.
type UpdateEventRequest struct {
	ID               string
	Name             *string                `json:"name,omitempty"`
	Description      *string                `json:"description,omitempty"`
	PropertiesSchema *EventPropertiesSchema `json:"propertiesSchema,omitempty"`
}

// ListEventsRequest is the request to list events.
type ListEventsRequest struct {
	ProjectSlug *string            `url:"projectSlug,omitempty"`
	Environment *types.Environment `url:"environment,omitempty"`
	Limit       *int               `url:"limit,omitempty"`
	Offset      *int               `url:"offset,omitempty"`
	Search      *string            `url:"search,omitempty"`
}

// ListEventsResponse is the response when listing events.
type ListEventsResponse struct {
	Events []MailEvent `json:"events"`
	Total  int         `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}

// DeleteEventResponse is the response when deleting an event.
type DeleteEventResponse struct {
	Success bool `json:"success"`
}

// TrackEventRequest is the request to track an event.
type TrackEventRequest struct {
	Environment  *types.Environment     `json:"environment,omitempty"`
	EventName    string                 `json:"eventName"`
	ContactID    *string                `json:"contactId,omitempty"`
	ContactEmail *string                `json:"contactEmail,omitempty"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
}

// TrackEventResponse is the response when tracking an event.
type TrackEventResponse struct {
	Success           bool    `json:"success"`
	EventOccurrenceID *string `json:"eventOccurrenceId,omitempty"`
	Error             *string `json:"error,omitempty"`
}

// BatchTrackEventInput represents a single event in a batch track request.
type BatchTrackEventInput struct {
	EventName    string                 `json:"eventName"`
	ContactID    *string                `json:"contactId,omitempty"`
	ContactEmail *string                `json:"contactEmail,omitempty"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
	Timestamp    *time.Time             `json:"timestamp,omitempty"`
}

// BatchTrackEventsRequest is the request to track multiple events.
type BatchTrackEventsRequest struct {
	Environment *types.Environment     `json:"environment,omitempty"`
	Events      []BatchTrackEventInput `json:"events"`
}

// BatchTrackEventResult represents the result of tracking a single event in a batch.
type BatchTrackEventResult struct {
	Success           bool    `json:"success"`
	EventOccurrenceID *string `json:"eventOccurrenceId,omitempty"`
	Error             *string `json:"error,omitempty"`
}

// BatchTrackEventsResponse is the response when tracking multiple events.
type BatchTrackEventsResponse struct {
	Success        bool                    `json:"success"`
	Results        []BatchTrackEventResult `json:"results"`
	TotalProcessed int                     `json:"totalProcessed"`
	TotalFailed    int                     `json:"totalFailed"`
}

// EventOccurrence represents a single occurrence of an event.
type EventOccurrence struct {
	ID          string                 `json:"id"`
	EventID     string                 `json:"eventId"`
	ContactID   string                 `json:"contactId"`
	Properties  map[string]interface{} `json:"properties"`
	Processed   bool                   `json:"processed"`
	ProcessedAt *time.Time             `json:"processedAt"`
	CreatedAt   time.Time              `json:"createdAt"`
}

// ListEventOccurrencesRequest is the request to list event occurrences.
type ListEventOccurrencesRequest struct {
	EventID   *string    `url:"eventId,omitempty"`
	ContactID *string    `url:"contactId,omitempty"`
	StartDate *time.Time `url:"startDate,omitempty"`
	EndDate   *time.Time `url:"endDate,omitempty"`
	Limit     *int       `url:"limit,omitempty"`
	Offset    *int       `url:"offset,omitempty"`
}

// ListEventOccurrencesResponse is the response when listing event occurrences.
type ListEventOccurrencesResponse struct {
	Occurrences []EventOccurrence `json:"occurrences"`
	Total       int               `json:"total"`
	Limit       int               `json:"limit"`
	Offset      int               `json:"offset"`
}

// EventAnalyticsResponse contains event analytics.
type EventAnalyticsResponse struct {
	TotalReceived  int        `json:"totalReceived"`
	LastReceivedAt *time.Time `json:"lastReceivedAt"`
	UniqueContacts int        `json:"uniqueContacts"`
	DailyCounts    []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	} `json:"dailyCounts"`
}
