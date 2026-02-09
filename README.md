# Stack0 Go SDK

The official Go SDK for [Stack0](https://stack0.dev) -- a modular platform for email, CDN, screenshots, and web data extraction.

## Installation

```bash
go get github.com/stack0dev/sdk-go
```

Requires Go 1.21 or later.

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	stack0 "github.com/stack0dev/sdk-go"
	"github.com/stack0dev/sdk-go/mail"
)

func main() {
	client := stack0.New("stack0_api_key")

	ctx := context.Background()

	resp, err := client.Mail.Send(ctx, &mail.SendEmailRequest{
		From:    &mail.EmailAddress{Email: "noreply@example.com", Name: "My App"},
		To:      []mail.EmailAddress{{Email: "user@example.com"}},
		Subject: "Hello from Stack0",
		HTML:    ptr("<h1>Welcome</h1><p>Thanks for signing up.</p>"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Email sent: %s (status: %s)\n", resp.ID, resp.Status)
}

func ptr[T any](v T) *T { return &v }
```

## Client Configuration

```go
// Default configuration (https://api.stack0.dev)
client := stack0.New("stack0_api_key")

// Custom API base URL
client := stack0.New("stack0_api_key", stack0.WithBaseURL("https://custom.api.example.com"))

// Custom CDN URL for client-side image transform URLs
client := stack0.New("stack0_api_key", stack0.WithCDNURL("https://cdn.example.com"))

// Multiple options
client := stack0.New("stack0_api_key",
	stack0.WithBaseURL("https://custom.api.example.com"),
	stack0.WithCDNURL("https://cdn.example.com"),
)
```

The client exposes four service modules:

| Property             | Description                          |
|----------------------|--------------------------------------|
| `client.Mail`        | Email sending, templates, campaigns  |
| `client.CDN`         | Asset upload, transforms, folders    |
| `client.Screenshots` | Webpage screenshot capture           |
| `client.Extraction`  | AI-powered web data extraction       |

---

## Mail

The Mail module handles transactional email, bulk sends, templates, audiences, contacts, campaigns, sequences, and events.

### Sending Email

```go
import "github.com/stack0dev/sdk-go/mail"

// Single email
resp, err := client.Mail.Send(ctx, &mail.SendEmailRequest{
	From:    "noreply@example.com",
	To:      []mail.EmailAddress{{Email: "user@example.com"}},
	Subject: "Order Confirmation",
	HTML:    ptr("<p>Your order has been confirmed.</p>"),
	Tags:    []string{"transactional", "orders"},
})

// With a template
resp, err := client.Mail.Send(ctx, &mail.SendEmailRequest{
	From:       "noreply@example.com",
	To:         []mail.EmailAddress{{Email: "user@example.com"}},
	Subject:    "Welcome",
	TemplateID: ptr("tmpl_abc123"),
	TemplateVariables: map[string]interface{}{
		"name":    "Alice",
		"company": "Acme Corp",
	},
})

// With attachments
resp, err := client.Mail.Send(ctx, &mail.SendEmailRequest{
	From:    "noreply@example.com",
	To:      []mail.EmailAddress{{Email: "user@example.com"}},
	Subject: "Your Invoice",
	HTML:    ptr("<p>Invoice attached.</p>"),
	Attachments: []mail.Attachment{
		{Filename: "invoice.pdf", Path: "https://example.com/invoice.pdf"},
	},
})

// Scheduled send
scheduledTime := time.Now().Add(2 * time.Hour)
resp, err := client.Mail.Send(ctx, &mail.SendEmailRequest{
	From:        "noreply@example.com",
	To:          []mail.EmailAddress{{Email: "user@example.com"}},
	Subject:     "Reminder",
	HTML:        ptr("<p>Don't forget your appointment.</p>"),
	ScheduledAt: &scheduledTime,
})
```

### Batch and Broadcast

```go
// Batch: send different emails to different recipients
batchResp, err := client.Mail.SendBatch(ctx, &mail.SendBatchEmailRequest{
	Emails: []mail.SendEmailRequest{
		{From: "noreply@example.com", To: "alice@example.com", Subject: "Hi Alice", HTML: ptr("<p>Hello</p>")},
		{From: "noreply@example.com", To: "bob@example.com", Subject: "Hi Bob", HTML: ptr("<p>Hello</p>")},
	},
})

// Broadcast: same email to many recipients
broadcastResp, err := client.Mail.SendBroadcast(ctx, &mail.SendBroadcastEmailRequest{
	From:    "noreply@example.com",
	To:      []interface{}{"alice@example.com", "bob@example.com", "carol@example.com"},
	Subject: "Product Update",
	HTML:    ptr("<p>Check out our latest features.</p>"),
})
```

### Managing Emails

```go
// Get email by ID
email, err := client.Mail.Get(ctx, "email_abc123")

// List emails with filters
emails, err := client.Mail.List(ctx, &mail.ListEmailsRequest{
	Status: ptr(mail.EmailStatusDelivered),
	Limit:  ptr(20),
})

// Resend a failed email
resendResp, err := client.Mail.Resend(ctx, "email_abc123")

// Cancel a scheduled email
cancelResp, err := client.Mail.Cancel(ctx, "email_abc123")
```

### Analytics

```go
analytics, err := client.Mail.GetAnalytics(ctx)
fmt.Printf("Delivery rate: %.2f%%\n", analytics.DeliveryRate*100)

timeSeries, err := client.Mail.GetTimeSeriesAnalytics(ctx, ptr(30))
hourly, err := client.Mail.GetHourlyAnalytics(ctx)

senders, err := client.Mail.ListSenders(ctx, &mail.ListSendersRequest{
	Search: ptr("noreply"),
})
```

### Domains

```go
// Add a domain
addResp, err := client.Mail.Domains.Add(ctx, &mail.AddDomainRequest{
	Domain: "mail.example.com",
})

// Get DNS records to configure
records, err := client.Mail.Domains.GetDNSRecords(ctx, "domain_id")

// Verify domain after configuring DNS
verifyResp, err := client.Mail.Domains.Verify(ctx, "domain_id")

// List all domains
domains, err := client.Mail.Domains.List(ctx, &mail.ListDomainsRequest{
	ProjectSlug: "my-project",
})

// Set default domain
client.Mail.Domains.SetDefault(ctx, "domain_id")

// Delete a domain
client.Mail.Domains.Delete(ctx, "domain_id")
```

### Templates

```go
// Create a template
tmpl, err := client.Mail.Templates.Create(ctx, &mail.CreateTemplateRequest{
	Name:    "Welcome Email",
	Slug:    "welcome-email",
	Subject: "Welcome, {{name}}!",
	HTML:    "<h1>Welcome, {{name}}</h1><p>We're glad to have you.</p>",
})

// List templates
templates, err := client.Mail.Templates.List(ctx, &mail.ListTemplatesRequest{
	Search:   ptr("welcome"),
	IsActive: ptr(true),
})

// Get by ID or slug
tmpl, err := client.Mail.Templates.Get(ctx, "tmpl_id")
tmpl, err := client.Mail.Templates.GetBySlug(ctx, "welcome-email")

// Preview with variables
preview, err := client.Mail.Templates.Preview(ctx, &mail.PreviewTemplateRequest{
	ID:        "tmpl_id",
	Variables: map[string]interface{}{"name": "Alice"},
})

// Update
client.Mail.Templates.Update(ctx, &mail.UpdateTemplateRequest{
	ID:      "tmpl_id",
	Subject: ptr("Updated Subject"),
})

// Delete
client.Mail.Templates.Delete(ctx, "tmpl_id")
```

### Audiences and Contacts

```go
// Create an audience
audience, err := client.Mail.Audiences.Create(ctx, &mail.CreateAudienceRequest{
	Name:        "Newsletter Subscribers",
	Description: ptr("Users opted in to the weekly newsletter"),
})

// Create a contact
contact, err := client.Mail.Contacts.Create(ctx, &mail.CreateContactRequest{
	Email:     "alice@example.com",
	FirstName: ptr("Alice"),
	LastName:  ptr("Smith"),
	Metadata:  map[string]interface{}{"plan": "pro"},
})

// Add contacts to an audience
client.Mail.Audiences.AddContacts(ctx, &mail.AddContactsToAudienceRequest{
	ID:         "audience_id",
	ContactIDs: []string{"contact_id_1", "contact_id_2"},
})

// Import contacts in bulk
importResp, err := client.Mail.Contacts.Import(ctx, &mail.ImportContactsRequest{
	AudienceID: ptr("audience_id"),
	Contacts: []mail.ImportContactInput{
		{Email: "user1@example.com", FirstName: ptr("User"), LastName: ptr("One")},
		{Email: "user2@example.com", FirstName: ptr("User"), LastName: ptr("Two")},
	},
})
fmt.Printf("Imported: %d, Skipped: %d\n", importResp.Imported, importResp.Skipped)

// List contacts in an audience
contacts, err := client.Mail.Audiences.ListContacts(ctx, &mail.ListAudienceContactsRequest{
	ID:    "audience_id",
	Limit: ptr(50),
})
```

### Campaigns

```go
// Create a campaign
campaign, err := client.Mail.Campaigns.Create(ctx, &mail.CreateCampaignRequest{
	Name:       "Product Launch",
	Subject:    "Introducing Our New Product",
	FromEmail:  "marketing@example.com",
	FromName:   ptr("Marketing Team"),
	AudienceID: ptr("audience_id"),
	TemplateID: ptr("tmpl_id"),
})

// Send immediately
sendResp, err := client.Mail.Campaigns.Send(ctx, &mail.SendCampaignRequest{
	ID:      "campaign_id",
	SendNow: ptr(true),
})

// Get campaign stats
stats, err := client.Mail.Campaigns.GetStats(ctx, "campaign_id")
fmt.Printf("Open rate: %.2f%%\n", stats.OpenRate*100)

// Pause, cancel, duplicate
client.Mail.Campaigns.Pause(ctx, "campaign_id")
client.Mail.Campaigns.Cancel(ctx, "campaign_id")
client.Mail.Campaigns.Duplicate(ctx, "campaign_id")
```

### Sequences

Sequences are visual automation flows with nodes (email, timer, filter, branch, experiment) and connections.

```go
// Create a sequence
seq, err := client.Mail.Sequences.Create(ctx, &mail.CreateSequenceRequest{
	Name:        "Onboarding Flow",
	TriggerType: mail.SequenceTriggerContactAdded,
})

// Add nodes
emailNode, err := client.Mail.Sequences.CreateNode(ctx, &mail.CreateNodeRequest{
	ID:        "seq_id",
	NodeType:  mail.SequenceNodeEmail,
	Name:      "Welcome Email",
	PositionX: 100,
	PositionY: 200,
})

timerNode, err := client.Mail.Sequences.CreateNode(ctx, &mail.CreateNodeRequest{
	ID:        "seq_id",
	NodeType:  mail.SequenceNodeTimer,
	Name:      "Wait 3 Days",
	PositionX: 100,
	PositionY: 400,
})

// Connect nodes
client.Mail.Sequences.CreateConnection(ctx, &mail.CreateConnectionRequest{
	ID:           "seq_id",
	SourceNodeID: emailNode.ID,
	TargetNodeID: timerNode.ID,
})

// Configure node email content
client.Mail.Sequences.SetNodeEmail(ctx, &mail.SetNodeEmailRequest{
	NodeID:  emailNode.ID,
	Subject: ptr("Welcome aboard!"),
	HTML:    ptr("<p>Thanks for joining.</p>"),
})

// Configure timer
client.Mail.Sequences.SetNodeTimer(ctx, &mail.SetNodeTimerRequest{
	NodeID:      timerNode.ID,
	DelayAmount: 3,
	DelayUnit:   "days",
})

// Publish the sequence
client.Mail.Sequences.Publish(ctx, "seq_id")

// Add a contact
client.Mail.Sequences.AddContact(ctx, &mail.AddContactToSequenceRequest{
	ID:        "seq_id",
	ContactID: "contact_id",
})

// Get analytics
analytics, err := client.Mail.Sequences.GetAnalytics(ctx, "seq_id")
```

### Events

Events allow tracking user actions that can trigger sequences.

```go
// Create an event definition
event, err := client.Mail.Events.Create(ctx, &mail.CreateEventRequest{
	Name:        "purchase_completed",
	Description: ptr("Fired when a user completes a purchase"),
})

// Track an event
trackResp, err := client.Mail.Events.Track(ctx, &mail.TrackEventRequest{
	EventName:    "purchase_completed",
	ContactEmail: ptr("alice@example.com"),
	Properties:   map[string]interface{}{"amount": 99.99, "product": "Pro Plan"},
})

// Batch track
client.Mail.Events.TrackBatch(ctx, &mail.BatchTrackEventsRequest{
	Events: []mail.BatchTrackEventInput{
		{EventName: "page_viewed", ContactEmail: ptr("alice@example.com"), Properties: map[string]interface{}{"page": "/pricing"}},
		{EventName: "page_viewed", ContactEmail: ptr("bob@example.com"), Properties: map[string]interface{}{"page": "/docs"}},
	},
})

// Get event analytics
eventAnalytics, err := client.Mail.Events.GetAnalytics(ctx, "event_id")
```

### Mail Method Reference

**Mail (direct)**

| Method                     | Description                        |
|----------------------------|------------------------------------|
| `Send`                     | Send a single email                |
| `SendBatch`                | Send multiple emails               |
| `SendBroadcast`            | Broadcast to many recipients       |
| `Get`                      | Get email by ID                    |
| `List`                     | List emails with filters           |
| `Resend`                   | Resend an email                    |
| `Cancel`                   | Cancel a scheduled email           |
| `GetAnalytics`             | Overall email analytics            |
| `GetTimeSeriesAnalytics`   | Time series analytics              |
| `GetHourlyAnalytics`       | Hourly send analytics              |
| `ListSenders`              | List unique senders with stats     |

**Mail.Domains**

| Method          | Description                     |
|-----------------|---------------------------------|
| `List`          | List domains                    |
| `Add`           | Add a new domain                |
| `GetDNSRecords` | Get DNS records for setup       |
| `Verify`        | Verify domain DNS configuration |
| `Delete`        | Remove a domain                 |
| `SetDefault`    | Set as default sending domain   |

**Mail.Templates**

| Method       | Description                      |
|--------------|----------------------------------|
| `List`       | List templates                   |
| `Get`        | Get template by ID               |
| `GetBySlug`  | Get template by slug             |
| `Create`     | Create a new template            |
| `Update`     | Update a template                |
| `Delete`     | Delete a template                |
| `Preview`    | Preview with template variables  |

**Mail.Audiences**

| Method           | Description                        |
|------------------|------------------------------------|
| `List`           | List audiences                     |
| `Get`            | Get audience by ID                 |
| `Create`         | Create an audience                 |
| `Update`         | Update an audience                 |
| `Delete`         | Delete an audience                 |
| `ListContacts`   | List contacts in an audience       |
| `AddContacts`    | Add contacts to an audience        |
| `RemoveContacts` | Remove contacts from an audience   |

**Mail.Contacts**

| Method    | Description               |
|-----------|---------------------------|
| `List`    | List contacts             |
| `Get`     | Get contact by ID         |
| `Create`  | Create a contact          |
| `Update`  | Update a contact          |
| `Delete`  | Delete a contact          |
| `Import`  | Bulk import contacts      |

**Mail.Campaigns**

| Method      | Description                  |
|-------------|------------------------------|
| `List`      | List campaigns               |
| `Get`       | Get campaign by ID           |
| `Create`    | Create a campaign            |
| `Update`    | Update a campaign            |
| `Delete`    | Delete a campaign            |
| `Send`      | Send or schedule a campaign  |
| `Pause`     | Pause a sending campaign     |
| `Cancel`    | Cancel a campaign            |
| `Duplicate` | Duplicate a campaign         |
| `GetStats`  | Get campaign statistics      |

**Mail.Sequences**

| Method             | Description                              |
|--------------------|------------------------------------------|
| `List`             | List sequences                           |
| `Get`              | Get sequence with nodes and connections  |
| `Create`           | Create a sequence                        |
| `Update`           | Update sequence settings                 |
| `Delete`           | Delete a sequence                        |
| `Publish`          | Publish a draft sequence                 |
| `Pause`            | Pause an active sequence                 |
| `Resume`           | Resume a paused sequence                 |
| `Archive`          | Archive a sequence                       |
| `Duplicate`        | Duplicate a sequence                     |
| `CreateNode`       | Add a node to the sequence               |
| `UpdateNode`       | Update a node                            |
| `DeleteNode`       | Remove a node                            |
| `SetNodeEmail`     | Set email content for a node             |
| `SetNodeTimer`     | Set timer delay for a node               |
| `SetNodeFilter`    | Set filter conditions for a node         |
| `SetNodeBranch`    | Set branching conditions for a node      |
| `SetNodeExperiment`| Set A/B experiment config for a node     |
| `CreateConnection` | Connect two nodes                        |
| `DeleteConnection` | Remove a connection                      |
| `ListEntries`      | List contacts in the sequence            |
| `AddContact`       | Add a contact to the sequence            |
| `RemoveContact`    | Remove a contact from the sequence       |
| `GetAnalytics`     | Get sequence performance analytics       |

**Mail.Events**

| Method             | Description                       |
|--------------------|-----------------------------------|
| `List`             | List event definitions            |
| `Get`              | Get event by ID                   |
| `Create`           | Create an event definition        |
| `Update`           | Update an event definition        |
| `Delete`           | Delete an event definition        |
| `Track`            | Track a single event occurrence   |
| `TrackBatch`       | Track multiple events at once     |
| `ListOccurrences`  | List event occurrences            |
| `GetAnalytics`     | Get event analytics               |

---

## CDN

The CDN module handles asset upload, image transforms, folders, private files, video processing, S3 imports, download bundles, and usage analytics.

### Uploading Assets

```go
import "github.com/stack0dev/sdk-go/cdn"

// 1. Get a presigned upload URL
upload, err := client.CDN.GetUploadURL(ctx, &cdn.UploadURLRequest{
	ProjectSlug: "my-project",
	Filename:    "hero.png",
	MimeType:    "image/png",
	Size:        1024000,
	Folder:      ptr("images/heroes"),
})

// 2. Upload the file to upload.UploadURL using an HTTP PUT
// (use your preferred HTTP client)

// 3. Confirm the upload
asset, err := client.CDN.ConfirmUpload(ctx, upload.AssetID)
fmt.Printf("Asset ready at: %s\n", asset.CDNURL)
```

### Image Transforms

Generate transformed image URLs client-side -- no API call required. Requires `WithCDNURL` to be set, or pass a full asset URL.

```go
client := stack0.New("key", stack0.WithCDNURL("https://cdn.example.com"))

// Resize and convert format
url, err := client.CDN.GetTransformURL("path/to/image.png", &cdn.TransformOptions{
	Width:   ptr(800),
	Height:  ptr(600),
	Format:  ptr("webp"),
	Quality: ptr(85),
})

// From a full URL
url, err := client.CDN.GetTransformURL("https://cdn.example.com/img.jpg", &cdn.TransformOptions{
	Width:     ptr(400),
	Grayscale: true,
	Blur:      ptr(5),
})
```

Available transform options:

| Option       | Type   | Description                               |
|--------------|--------|-------------------------------------------|
| `Width`      | `*int` | Target width (snapped to nearest allowed) |
| `Height`     | `*int` | Target height                             |
| `Format`     | `*string` | Output format (webp, jpeg, png, avif)  |
| `Quality`    | `*int` | Compression quality (1-100)               |
| `Fit`        | `*string` | Resize fit mode                        |
| `Crop`       | `*string` | Crop mode                              |
| `Blur`       | `*int` | Blur radius                               |
| `Sharpen`    | `*int` | Sharpen amount                            |
| `Brightness` | `*int` | Brightness adjustment                     |
| `Saturation` | `*int` | Saturation adjustment                     |
| `Grayscale`  | `bool` | Convert to grayscale                      |
| `Rotate`     | `*int` | Rotation in degrees                       |
| `Flip`       | `bool` | Flip vertically                           |
| `Flop`       | `bool` | Flip horizontally                         |

### Managing Assets

```go
// Get asset by ID
asset, err := client.CDN.Get(ctx, "asset_id")

// List assets
assets, err := client.CDN.List(ctx, &cdn.ListAssetsRequest{
	ProjectSlug: "my-project",
	Type:        ptr(cdn.AssetTypeImage),
	Search:      ptr("hero"),
	Limit:       ptr(20),
})

// Update metadata
asset, err := client.CDN.Update(ctx, &cdn.UpdateAssetRequest{
	ID:   "asset_id",
	Alt:  ptr("Hero banner image"),
	Tags: []string{"banner", "homepage"},
})

// Move assets to a folder
client.CDN.Move(ctx, &cdn.MoveAssetsRequest{
	AssetIDs: []string{"asset_1", "asset_2"},
	Folder:   ptr("archive"),
})

// Delete
client.CDN.Delete(ctx, "asset_id")
client.CDN.DeleteMany(ctx, []string{"asset_1", "asset_2"})
```

### Folders

```go
// Create a folder
folder, err := client.CDN.CreateFolder(ctx, &cdn.CreateFolderRequest{
	ProjectSlug: "my-project",
	Name:        "marketing",
	ParentID:    ptr("parent_folder_id"),
})

// Get folder tree
tree, err := client.CDN.GetFolderTree(ctx, &cdn.GetFolderTreeRequest{
	ProjectSlug: "my-project",
	MaxDepth:    ptr(3),
})

// Get folder by path
folder, err := client.CDN.GetFolderByPath(ctx, "images/marketing")

// List, move, delete
folders, err := client.CDN.ListFolders(ctx, &cdn.ListFoldersRequest{Search: ptr("market")})
client.CDN.MoveFolder(ctx, &cdn.MoveFolderRequest{ID: "folder_id", NewParentID: ptr("new_parent_id")})
client.CDN.DeleteFolder(ctx, "folder_id", true) // true = delete contents
```

### Private Files

Private files are stored securely and accessed via time-limited presigned URLs.

```go
// Upload a private file
upload, err := client.CDN.GetPrivateUploadURL(ctx, &cdn.PrivateUploadURLRequest{
	ProjectSlug: "my-project",
	Filename:    "contract.pdf",
	MimeType:    "application/pdf",
	Size:        2048000,
})
// Upload the file, then confirm
file, err := client.CDN.ConfirmPrivateUpload(ctx, upload.FileID)

// Get a temporary download URL
download, err := client.CDN.GetPrivateDownloadURL(ctx, &cdn.PrivateDownloadURLRequest{
	FileID:    "file_id",
	ExpiresIn: ptr(3600), // seconds
})
fmt.Println(download.DownloadURL)

// List, update, delete
files, err := client.CDN.ListPrivateFiles(ctx, &cdn.ListPrivateFilesRequest{ProjectSlug: "my-project"})
client.CDN.UpdatePrivateFile(ctx, &cdn.UpdatePrivateFileRequest{FileID: "file_id", Filename: ptr("renamed.pdf")})
client.CDN.DeletePrivateFile(ctx, "file_id")
```

### Video Processing

```go
// Transcode a video
job, err := client.CDN.Transcode(ctx, &cdn.TranscodeVideoRequest{...})

// Get streaming URLs (HLS/DASH)
urls, err := client.CDN.GetStreamingURLs(ctx, "asset_id")

// Generate a thumbnail
thumb, err := client.CDN.GetThumbnail(ctx, &cdn.ThumbnailRequest{
	AssetID:   "asset_id",
	Timestamp: 5.0, // seconds
	Width:     ptr(320),
})

// Generate a GIF from a video segment
gif, err := client.CDN.GenerateGif(ctx, &cdn.GenerateGifRequest{...})

// Merge videos
mergeJob, err := client.CDN.CreateMergeJob(ctx, &cdn.CreateMergeJobRequest{...})

// Extract audio
audio, err := client.CDN.ExtractAudio(ctx, &cdn.ExtractAudioRequest{...})
```

### S3 Imports

```go
// Bulk import from S3
importJob, err := client.CDN.CreateImport(ctx, &cdn.CreateImportRequest{...})

// Monitor progress
job, err := client.CDN.GetImport(ctx, "import_id")

// List files in an import
files, err := client.CDN.ListImportFiles(ctx, &cdn.ListImportFilesRequest{ImportID: "import_id"})

// Cancel or retry
client.CDN.CancelImport(ctx, "import_id")
client.CDN.RetryImport(ctx, "import_id")
```

### Download Bundles

```go
bundle, err := client.CDN.CreateBundle(ctx, &cdn.CreateBundleRequest{...})
download, err := client.CDN.GetBundleDownloadURL(ctx, &cdn.BundleDownloadURLRequest{BundleID: "bundle_id"})
```

### CDN Usage

```go
usage, err := client.CDN.GetUsage(ctx, &cdn.CdnUsageRequest{ProjectSlug: ptr("my-project")})
history, err := client.CDN.GetUsageHistory(ctx, &cdn.CdnUsageHistoryRequest{Days: ptr(30)})
breakdown, err := client.CDN.GetStorageBreakdown(ctx, &cdn.CdnStorageBreakdownRequest{GroupBy: ptr("type")})
```

---

## Screenshots

Capture webpage screenshots with full control over viewport, format, blocking, and rendering options.

### Basic Capture

```go
import "github.com/stack0dev/sdk-go/screenshots"

// Capture and wait for the result (polls until complete)
screenshot, err := client.Screenshots.CaptureAndWait(ctx, &screenshots.CreateScreenshotRequest{
	URL:    "https://example.com",
	Format: ptr(screenshots.ScreenshotFormatWebP),
}, nil)
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Screenshot: %s (%dx%d)\n", *screenshot.ImageURL, *screenshot.ImageWidth, *screenshot.ImageHeight)
```

### Advanced Options

```go
screenshot, err := client.Screenshots.CaptureAndWait(ctx, &screenshots.CreateScreenshotRequest{
	URL:                "https://example.com",
	Format:             ptr(screenshots.ScreenshotFormatPNG),
	Quality:            ptr(90),
	FullPage:           ptr(true),
	DeviceType:         ptr(screenshots.DeviceTypeMobile),
	ViewportWidth:      ptr(375),
	ViewportHeight:     ptr(812),
	BlockAds:           ptr(true),
	BlockCookieBanners: ptr(true),
	BlockChatWidgets:   ptr(true),
	DarkMode:           ptr(true),
	WaitForSelector:    ptr("#main-content"),
	WaitForTimeout:     ptr(2000),
	CustomCSS:          ptr("body { font-size: 16px; }"),
	HideSelectors:      []string{".cookie-banner", ".popup"},
	Clip:               &screenshots.Clip{X: 0, Y: 0, Width: 1200, Height: 630},
}, &screenshots.CaptureAndWaitOptions{
	PollInterval: 500 * time.Millisecond,
	Timeout:      30 * time.Second,
})
```

### Async Capture

```go
// Start capture without waiting
resp, err := client.Screenshots.Capture(ctx, &screenshots.CreateScreenshotRequest{
	URL:        "https://example.com",
	WebhookURL: ptr("https://myapp.com/webhook/screenshot"),
})

// Check status later
screenshot, err := client.Screenshots.Get(ctx, &screenshots.GetScreenshotRequest{
	ID: resp.ID,
})
```

### Batch Screenshots

```go
// Capture multiple URLs
batchResp, err := client.Screenshots.Batch(ctx, &screenshots.CreateBatchScreenshotsRequest{
	URLs: []string{
		"https://example.com",
		"https://example.com/about",
		"https://example.com/pricing",
	},
	Config: &screenshots.BatchScreenshotConfig{
		Format:   ptr(screenshots.ScreenshotFormatWebP),
		FullPage: ptr(true),
	},
})

// Wait for all to complete
job, err := client.Screenshots.BatchAndWait(ctx, &screenshots.CreateBatchScreenshotsRequest{
	URLs: []string{"https://example.com", "https://example.com/blog"},
}, &screenshots.CaptureAndWaitOptions{
	Timeout: 5 * time.Minute,
})
fmt.Printf("Processed: %d, Successful: %d\n", job.ProcessedURLs, job.SuccessfulURLs)
```

### Scheduled Screenshots

```go
import "github.com/stack0dev/sdk-go/types"

// Create a recurring screenshot schedule
schedResp, err := client.Screenshots.CreateSchedule(ctx, &screenshots.CreateScreenshotScheduleRequest{
	Name:          "Homepage Monitor",
	URL:           "https://example.com",
	Frequency:     ptr(types.ScheduleFrequencyDaily),
	DetectChanges: ptr(true),
	WebhookURL:    ptr("https://myapp.com/webhook/change-detected"),
})

// List schedules
schedules, err := client.Screenshots.ListSchedules(ctx, &screenshots.ListSchedulesRequest{
	IsActive: ptr(true),
})

// Toggle on/off
client.Screenshots.ToggleSchedule(ctx, &screenshots.GetScheduleRequest{ID: "sched_id"})
```

### Screenshots Method Reference

| Method             | Description                              |
|--------------------|------------------------------------------|
| `Capture`          | Start async screenshot capture           |
| `CaptureAndWait`   | Capture and poll until complete          |
| `Get`              | Get screenshot by ID                     |
| `List`             | List screenshots                         |
| `Delete`           | Delete a screenshot                      |
| `Batch`            | Create batch screenshot job              |
| `BatchAndWait`     | Create batch and poll until complete     |
| `GetBatchJob`      | Get batch job status                     |
| `ListBatchJobs`    | List batch jobs                          |
| `CancelBatchJob`   | Cancel a batch job                       |
| `CreateSchedule`   | Create a recurring schedule              |
| `UpdateSchedule`   | Update a schedule                        |
| `GetSchedule`      | Get schedule by ID                       |
| `ListSchedules`    | List schedules                           |
| `DeleteSchedule`   | Delete a schedule                        |
| `ToggleSchedule`   | Toggle schedule active/inactive          |

---

## Extraction

Extract structured data from any webpage using AI. Supports schema-based extraction, markdown conversion, and raw HTML retrieval.

### Basic Extraction

```go
import "github.com/stack0dev/sdk-go/extraction"

// Extract and wait for results
result, err := client.Extraction.ExtractAndWait(ctx, &extraction.CreateExtractionRequest{
	URL:  "https://example.com/product",
	Mode: ptr(extraction.ExtractionModeAuto),
}, nil)
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Extracted data: %v\n", result.ExtractedData)
```

### Schema-Based Extraction

```go
result, err := client.Extraction.ExtractAndWait(ctx, &extraction.CreateExtractionRequest{
	URL:  "https://example.com/product",
	Mode: ptr(extraction.ExtractionModeSchema),
	Schema: map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"title":       map[string]interface{}{"type": "string"},
			"price":       map[string]interface{}{"type": "number"},
			"description": map[string]interface{}{"type": "string"},
			"inStock":     map[string]interface{}{"type": "boolean"},
		},
	},
	Prompt: ptr("Extract the main product information from this page"),
}, nil)
```

### Markdown Extraction

```go
result, err := client.Extraction.ExtractAndWait(ctx, &extraction.CreateExtractionRequest{
	URL:             "https://example.com/blog/post",
	Mode:            ptr(extraction.ExtractionModeMarkdown),
	IncludeLinks:    ptr(true),
	IncludeImages:   ptr(true),
	IncludeMetadata: ptr(true),
}, nil)

fmt.Println(*result.Markdown)
fmt.Printf("Page title: %s\n", *result.PageMetadata.Title)
```

### Batch Extraction

```go
batchResp, err := client.Extraction.Batch(ctx, &extraction.CreateBatchExtractionsRequest{
	URLs: []string{
		"https://example.com/page1",
		"https://example.com/page2",
		"https://example.com/page3",
	},
	Config: &extraction.BatchExtractionConfig{
		Mode: ptr(extraction.ExtractionModeMarkdown),
	},
})

// Wait for completion
job, err := client.Extraction.BatchAndWait(ctx, &extraction.CreateBatchExtractionsRequest{
	URLs: []string{"https://example.com/page1", "https://example.com/page2"},
}, &extraction.ExtractAndWaitOptions{
	Timeout: 5 * time.Minute,
})
```

### Scheduled Extraction

```go
schedResp, err := client.Extraction.CreateSchedule(ctx, &extraction.CreateExtractionScheduleRequest{
	Name:          "Price Monitor",
	URL:           "https://competitor.com/pricing",
	Frequency:     ptr(types.ScheduleFrequencyDaily),
	DetectChanges: ptr(true),
	Config: &extraction.BatchExtractionConfig{
		Mode: ptr(extraction.ExtractionModeSchema),
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"plans": map[string]interface{}{"type": "array"},
			},
		},
	},
})
```

### Usage Statistics

```go
usage, err := client.Extraction.GetUsage(ctx, &extraction.GetUsageRequest{})
fmt.Printf("Extractions this period: %d (credits: %d)\n", usage.ExtractionsTotal, usage.ExtractionCreditsUsed)

daily, err := client.Extraction.GetUsageDaily(ctx, &extraction.GetUsageRequest{})
for _, day := range daily.Days {
	fmt.Printf("%s: %d extractions\n", day.Date, day.Extractions)
}
```

### Extraction Method Reference

| Method             | Description                              |
|--------------------|------------------------------------------|
| `Extract`          | Start async extraction                   |
| `ExtractAndWait`   | Extract and poll until complete          |
| `Get`              | Get extraction by ID                     |
| `List`             | List extractions                         |
| `Delete`           | Delete an extraction                     |
| `Batch`            | Create batch extraction job              |
| `BatchAndWait`     | Create batch and poll until complete     |
| `GetBatchJob`      | Get batch job status                     |
| `ListBatchJobs`    | List batch jobs                          |
| `CancelBatchJob`   | Cancel a batch job                       |
| `CreateSchedule`   | Create a recurring schedule              |
| `UpdateSchedule`   | Update a schedule                        |
| `GetSchedule`      | Get schedule by ID                       |
| `ListSchedules`    | List schedules                           |
| `DeleteSchedule`   | Delete a schedule                        |
| `ToggleSchedule`   | Toggle schedule active/inactive          |
| `GetUsage`         | Get usage statistics                     |
| `GetUsageDaily`    | Get daily usage breakdown                |

---

## Error Handling

All methods return idiomatic Go errors. API errors are returned as `*types.APIError`, and polling timeouts as `*types.TimeoutError`.

```go
import (
	"errors"
	"github.com/stack0dev/sdk-go/types"
)

resp, err := client.Mail.Send(ctx, &mail.SendEmailRequest{...})
if err != nil {
	var apiErr *types.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("API error %d: %s (code: %s)\n", apiErr.StatusCode, apiErr.Message, apiErr.Code)

		switch apiErr.StatusCode {
		case 400:
			// Bad request -- check your parameters
		case 401:
			// Invalid API key
		case 403:
			// Insufficient permissions
		case 404:
			// Resource not found
		case 429:
			// Rate limited -- back off and retry
		}
		return
	}

	var timeoutErr *types.TimeoutError
	if errors.As(err, &timeoutErr) {
		fmt.Println("Operation timed out:", timeoutErr.Message)
		return
	}

	// Network error, context cancellation, etc.
	fmt.Println("Unexpected error:", err)
}
```

### Context Cancellation

All methods accept a `context.Context` as the first parameter, providing cancellation and deadline support.

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

screenshot, err := client.Screenshots.CaptureAndWait(ctx, &screenshots.CreateScreenshotRequest{
	URL: "https://example.com",
}, nil)
if err != nil {
	// err may be context.DeadlineExceeded
}
```

---

## Shared Types

The `types` package provides constants and types shared across modules.

```go
import "github.com/stack0dev/sdk-go/types"
```

**Environments**

| Constant                     | Value          |
|------------------------------|----------------|
| `types.EnvironmentSandbox`   | `"sandbox"`    |
| `types.EnvironmentProduction`| `"production"` |

**Batch Job Status**

| Constant                          | Value          |
|-----------------------------------|----------------|
| `types.BatchJobStatusPending`     | `"pending"`    |
| `types.BatchJobStatusProcessing`  | `"processing"` |
| `types.BatchJobStatusCompleted`   | `"completed"`  |
| `types.BatchJobStatusFailed`      | `"failed"`     |
| `types.BatchJobStatusCancelled`   | `"cancelled"`  |

**Schedule Frequency**

| Constant                            | Value       |
|-------------------------------------|-------------|
| `types.ScheduleFrequencyHourly`     | `"hourly"`  |
| `types.ScheduleFrequencyDaily`      | `"daily"`   |
| `types.ScheduleFrequencyWeekly`     | `"weekly"`  |
| `types.ScheduleFrequencyMonthly`    | `"monthly"` |

---

## Full API Reference

For complete API documentation, request/response schemas, and guides, visit the [Stack0 Documentation](https://stack0.dev/docs).

## License

MIT License. See [LICENSE](LICENSE) for details.
