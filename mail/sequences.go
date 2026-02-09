package mail

import (
	"context"
	"net/url"
	"strconv"

	"github.com/stack0/sdk-go/client"
)

// SequencesClient handles sequence operations.
type SequencesClient struct {
	http *client.HTTPClient
}

// NewSequencesClient creates a new sequences client.
func NewSequencesClient(http *client.HTTPClient) *SequencesClient {
	return &SequencesClient{http: http}
}

// List lists all sequences.
func (c *SequencesClient) List(ctx context.Context, req *ListSequencesRequest) (*ListSequencesResponse, error) {
	params := url.Values{}
	if req != nil {
		if req.Environment != nil {
			params.Set("environment", string(*req.Environment))
		}
		if req.Limit != nil {
			params.Set("limit", strconv.Itoa(*req.Limit))
		}
		if req.Offset != nil {
			params.Set("offset", strconv.Itoa(*req.Offset))
		}
		if req.Search != nil {
			params.Set("search", *req.Search)
		}
		if req.Status != nil {
			params.Set("status", string(*req.Status))
		}
		if req.TriggerType != nil {
			params.Set("triggerType", string(*req.TriggerType))
		}
	}

	path := "/mail/sequences"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListSequencesResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a sequence by ID with nodes and connections.
func (c *SequencesClient) Get(ctx context.Context, id string) (*SequenceWithNodes, error) {
	var resp SequenceWithNodes
	if err := c.http.Get(ctx, "/mail/sequences/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new sequence.
func (c *SequencesClient) Create(ctx context.Context, req *CreateSequenceRequest) (*Sequence, error) {
	var resp Sequence
	if err := c.http.Post(ctx, "/mail/sequences", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates a sequence.
func (c *SequencesClient) Update(ctx context.Context, req *UpdateSequenceRequest) (*Sequence, error) {
	var resp Sequence
	if err := c.http.Put(ctx, "/mail/sequences/"+req.ID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes a sequence.
func (c *SequencesClient) Delete(ctx context.Context, id string) (*DeleteSequenceResponse, error) {
	var resp DeleteSequenceResponse
	if err := c.http.Delete(ctx, "/mail/sequences/"+id, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Publish publishes (activates) a sequence.
func (c *SequencesClient) Publish(ctx context.Context, id string) (*PublishSequenceResponse, error) {
	var resp PublishSequenceResponse
	if err := c.http.Post(ctx, "/mail/sequences/"+id+"/publish", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Pause pauses an active sequence.
func (c *SequencesClient) Pause(ctx context.Context, id string) (*PauseSequenceResponse, error) {
	var resp PauseSequenceResponse
	if err := c.http.Post(ctx, "/mail/sequences/"+id+"/pause", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Resume resumes a paused sequence.
func (c *SequencesClient) Resume(ctx context.Context, id string) (*ResumeSequenceResponse, error) {
	var resp ResumeSequenceResponse
	if err := c.http.Post(ctx, "/mail/sequences/"+id+"/resume", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Archive archives a sequence.
func (c *SequencesClient) Archive(ctx context.Context, id string) (*ArchiveSequenceResponse, error) {
	var resp ArchiveSequenceResponse
	if err := c.http.Post(ctx, "/mail/sequences/"+id+"/archive", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Duplicate duplicates a sequence.
func (c *SequencesClient) Duplicate(ctx context.Context, id string, name *string) (*Sequence, error) {
	body := map[string]interface{}{}
	if name != nil {
		body["name"] = *name
	}
	var resp Sequence
	if err := c.http.Post(ctx, "/mail/sequences/"+id+"/duplicate", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateNode creates a new node in a sequence.
func (c *SequencesClient) CreateNode(ctx context.Context, req *CreateNodeRequest) (*SequenceNode, error) {
	var resp SequenceNode
	if err := c.http.Post(ctx, "/mail/sequences/"+req.ID+"/nodes", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateNode updates a node.
func (c *SequencesClient) UpdateNode(ctx context.Context, req *UpdateNodeRequest) (*SequenceNode, error) {
	var resp SequenceNode
	if err := c.http.Put(ctx, "/mail/sequences/"+req.ID+"/nodes/"+req.NodeID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateNodePosition updates a node's position.
func (c *SequencesClient) UpdateNodePosition(ctx context.Context, req *UpdateNodePositionRequest) (*SequenceNode, error) {
	body := map[string]interface{}{
		"positionX": req.PositionX,
		"positionY": req.PositionY,
	}
	var resp SequenceNode
	if err := c.http.Put(ctx, "/mail/sequences/"+req.ID+"/nodes/"+req.NodeID+"/position", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteNode deletes a node.
func (c *SequencesClient) DeleteNode(ctx context.Context, sequenceID, nodeID string) (*DeleteNodeResponse, error) {
	var resp DeleteNodeResponse
	if err := c.http.Delete(ctx, "/mail/sequences/"+sequenceID+"/nodes/"+nodeID, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetNodeEmail sets email content for a node.
func (c *SequencesClient) SetNodeEmail(ctx context.Context, sequenceID string, req *SetNodeEmailRequest) (*SequenceNode, error) {
	var resp SequenceNode
	if err := c.http.Put(ctx, "/mail/sequences/"+sequenceID+"/nodes/"+req.NodeID+"/email", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetNodeTimer sets timer configuration for a node.
func (c *SequencesClient) SetNodeTimer(ctx context.Context, sequenceID string, req *SetNodeTimerRequest) (*SequenceNode, error) {
	var resp SequenceNode
	if err := c.http.Put(ctx, "/mail/sequences/"+sequenceID+"/nodes/"+req.NodeID+"/timer", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetNodeFilter sets filter configuration for a node.
func (c *SequencesClient) SetNodeFilter(ctx context.Context, sequenceID string, req *SetNodeFilterRequest) (*SequenceNode, error) {
	var resp SequenceNode
	if err := c.http.Put(ctx, "/mail/sequences/"+sequenceID+"/nodes/"+req.NodeID+"/filter", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetNodeBranch sets branch configuration for a node.
func (c *SequencesClient) SetNodeBranch(ctx context.Context, sequenceID string, req *SetNodeBranchRequest) (*SequenceNode, error) {
	var resp SequenceNode
	if err := c.http.Put(ctx, "/mail/sequences/"+sequenceID+"/nodes/"+req.NodeID+"/branch", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetNodeExperiment sets experiment configuration for a node.
func (c *SequencesClient) SetNodeExperiment(ctx context.Context, sequenceID string, req *SetNodeExperimentRequest) (*SequenceNode, error) {
	var resp SequenceNode
	if err := c.http.Put(ctx, "/mail/sequences/"+sequenceID+"/nodes/"+req.NodeID+"/experiment", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateConnection creates a connection between nodes.
func (c *SequencesClient) CreateConnection(ctx context.Context, req *CreateConnectionRequest) (*SequenceConnection, error) {
	var resp SequenceConnection
	if err := c.http.Post(ctx, "/mail/sequences/"+req.ID+"/connections", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteConnection deletes a connection.
func (c *SequencesClient) DeleteConnection(ctx context.Context, sequenceID, connectionID string) (*DeleteConnectionResponse, error) {
	var resp DeleteConnectionResponse
	if err := c.http.Delete(ctx, "/mail/sequences/"+sequenceID+"/connections/"+connectionID, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListEntries lists contacts in a sequence.
func (c *SequencesClient) ListEntries(ctx context.Context, req *ListSequenceEntriesRequest) (*ListSequenceEntriesResponse, error) {
	params := url.Values{}
	if req.Limit != nil {
		params.Set("limit", strconv.Itoa(*req.Limit))
	}
	if req.Offset != nil {
		params.Set("offset", strconv.Itoa(*req.Offset))
	}
	if req.Status != nil {
		params.Set("status", string(*req.Status))
	}

	path := "/mail/sequences/" + req.ID + "/entries"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListSequenceEntriesResponse
	if err := c.http.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AddContact adds a contact to a sequence.
func (c *SequencesClient) AddContact(ctx context.Context, req *AddContactToSequenceRequest) (*SequenceEntry, error) {
	body := map[string]interface{}{"contactId": req.ContactID}
	var resp SequenceEntry
	if err := c.http.Post(ctx, "/mail/sequences/"+req.ID+"/add-contact", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RemoveContact removes a contact from a sequence.
func (c *SequencesClient) RemoveContact(ctx context.Context, req *RemoveContactFromSequenceRequest) (*RemoveContactFromSequenceResponse, error) {
	body := map[string]interface{}{"entryId": req.EntryID}
	if req.Reason != nil {
		body["reason"] = *req.Reason
	}
	var resp RemoveContactFromSequenceResponse
	if err := c.http.Post(ctx, "/mail/sequences/"+req.ID+"/remove-contact", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAnalytics retrieves sequence analytics.
func (c *SequencesClient) GetAnalytics(ctx context.Context, id string) (*SequenceAnalyticsResponse, error) {
	var resp SequenceAnalyticsResponse
	if err := c.http.Get(ctx, "/mail/sequences/"+id+"/analytics", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
