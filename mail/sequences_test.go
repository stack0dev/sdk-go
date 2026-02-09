package mail

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

func setupSequencesTestClient(t *testing.T, handler http.HandlerFunc) (*SequencesClient, *httptest.Server) {
	server := httptest.NewServer(handler)
	httpClient := client.New("test-api-key", server.URL)
	return NewSequencesClient(httpClient), server
}

func TestSequencesClient_List(t *testing.T) {
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/mail/sequences")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListSequencesResponse{
			Sequences: []Sequence{
				{ID: "seq-1", Name: "Welcome Series", Status: SequenceStatusActive},
				{ID: "seq-2", Name: "Onboarding", Status: SequenceStatusDraft},
			},
			Total: 2,
		})
	})
	defer server.Close()

	resp, err := sequencesClient.List(context.Background(), nil)

	require.NoError(t, err)
	assert.Len(t, resp.Sequences, 2)
}

func TestSequencesClient_Get(t *testing.T) {
	sequenceID := "seq-123"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SequenceWithNodes{
			Sequence: Sequence{
				ID:          sequenceID,
				Name:        "Welcome Series",
				TriggerType: SequenceTriggerContactAdded,
				Status:      SequenceStatusActive,
			},
			Nodes: []SequenceNode{
				{ID: "node-1", NodeType: SequenceNodeTrigger, Name: "Start"},
				{ID: "node-2", NodeType: SequenceNodeEmail, Name: "Welcome Email"},
			},
			Connections: []SequenceConnection{
				{ID: "conn-1", SourceNodeID: "node-1", TargetNodeID: "node-2"},
			},
		})
	})
	defer server.Close()

	resp, err := sequencesClient.Get(context.Background(), sequenceID)

	require.NoError(t, err)
	assert.Equal(t, sequenceID, resp.ID)
	assert.Len(t, resp.Nodes, 2)
	assert.Len(t, resp.Connections, 1)
}

func TestSequencesClient_Create(t *testing.T) {
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/sequences", r.URL.Path)

		var req CreateSequenceRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "New Sequence", req.Name)
		assert.Equal(t, SequenceTriggerManual, req.TriggerType)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Sequence{
			ID:          "seq-new",
			Name:        req.Name,
			TriggerType: req.TriggerType,
			Status:      SequenceStatusDraft,
		})
	})
	defer server.Close()

	resp, err := sequencesClient.Create(context.Background(), &CreateSequenceRequest{
		Name:        "New Sequence",
		TriggerType: SequenceTriggerManual,
	})

	require.NoError(t, err)
	assert.Equal(t, "seq-new", resp.ID)
	assert.Equal(t, SequenceStatusDraft, resp.Status)
}

func TestSequencesClient_Update(t *testing.T) {
	sequenceID := "seq-123"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Sequence{
			ID:   sequenceID,
			Name: "Updated Sequence",
		})
	})
	defer server.Close()

	name := "Updated Sequence"
	resp, err := sequencesClient.Update(context.Background(), &UpdateSequenceRequest{
		ID:   sequenceID,
		Name: &name,
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated Sequence", resp.Name)
}

func TestSequencesClient_Delete(t *testing.T) {
	sequenceID := "seq-123"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DeleteSequenceResponse{Success: true})
	})
	defer server.Close()

	resp, err := sequencesClient.Delete(context.Background(), sequenceID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestSequencesClient_Publish(t *testing.T) {
	sequenceID := "seq-123"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID+"/publish", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PublishSequenceResponse{Success: true})
	})
	defer server.Close()

	resp, err := sequencesClient.Publish(context.Background(), sequenceID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestSequencesClient_Pause(t *testing.T) {
	sequenceID := "seq-123"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID+"/pause", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PauseSequenceResponse{Success: true})
	})
	defer server.Close()

	resp, err := sequencesClient.Pause(context.Background(), sequenceID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestSequencesClient_Resume(t *testing.T) {
	sequenceID := "seq-123"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID+"/resume", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ResumeSequenceResponse{Success: true})
	})
	defer server.Close()

	resp, err := sequencesClient.Resume(context.Background(), sequenceID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestSequencesClient_Archive(t *testing.T) {
	sequenceID := "seq-123"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID+"/archive", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ArchiveSequenceResponse{Success: true})
	})
	defer server.Close()

	resp, err := sequencesClient.Archive(context.Background(), sequenceID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestSequencesClient_Duplicate(t *testing.T) {
	sequenceID := "seq-123"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID+"/duplicate", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Sequence{
			ID:     "seq-456",
			Name:   "Original Sequence (Copy)",
			Status: SequenceStatusDraft,
		})
	})
	defer server.Close()

	resp, err := sequencesClient.Duplicate(context.Background(), sequenceID, nil)

	require.NoError(t, err)
	assert.Equal(t, "seq-456", resp.ID)
}

func TestSequencesClient_CreateNode(t *testing.T) {
	sequenceID := "seq-123"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID+"/nodes", r.URL.Path)

		var req CreateNodeRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, SequenceNodeEmail, req.NodeType)
		assert.Equal(t, "Email Node", req.Name)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SequenceNode{
			ID:       "node-new",
			NodeType: req.NodeType,
			Name:     req.Name,
		})
	})
	defer server.Close()

	resp, err := sequencesClient.CreateNode(context.Background(), &CreateNodeRequest{
		ID:        sequenceID,
		NodeType:  SequenceNodeEmail,
		Name:      "Email Node",
		PositionX: 100,
		PositionY: 200,
	})

	require.NoError(t, err)
	assert.Equal(t, "node-new", resp.ID)
	assert.Equal(t, SequenceNodeEmail, resp.NodeType)
}

func TestSequencesClient_UpdateNode(t *testing.T) {
	sequenceID := "seq-123"
	nodeID := "node-456"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID+"/nodes/"+nodeID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SequenceNode{
			ID:   nodeID,
			Name: "Updated Node",
		})
	})
	defer server.Close()

	name := "Updated Node"
	resp, err := sequencesClient.UpdateNode(context.Background(), &UpdateNodeRequest{
		ID:     sequenceID,
		NodeID: nodeID,
		Name:   &name,
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated Node", resp.Name)
}

func TestSequencesClient_DeleteNode(t *testing.T) {
	sequenceID := "seq-123"
	nodeID := "node-456"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID+"/nodes/"+nodeID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DeleteNodeResponse{Success: true})
	})
	defer server.Close()

	resp, err := sequencesClient.DeleteNode(context.Background(), sequenceID, nodeID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestSequencesClient_CreateConnection(t *testing.T) {
	sequenceID := "seq-123"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID+"/connections", r.URL.Path)

		var req CreateConnectionRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "node-1", req.SourceNodeID)
		assert.Equal(t, "node-2", req.TargetNodeID)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SequenceConnection{
			ID:           "conn-new",
			SourceNodeID: req.SourceNodeID,
			TargetNodeID: req.TargetNodeID,
			CreatedAt:    time.Now(),
		})
	})
	defer server.Close()

	resp, err := sequencesClient.CreateConnection(context.Background(), &CreateConnectionRequest{
		ID:           sequenceID,
		SourceNodeID: "node-1",
		TargetNodeID: "node-2",
	})

	require.NoError(t, err)
	assert.Equal(t, "conn-new", resp.ID)
}

func TestSequencesClient_DeleteConnection(t *testing.T) {
	sequenceID := "seq-123"
	connectionID := "conn-456"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID+"/connections/"+connectionID, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DeleteConnectionResponse{Success: true})
	})
	defer server.Close()

	resp, err := sequencesClient.DeleteConnection(context.Background(), sequenceID, connectionID)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestSequencesClient_ListEntries(t *testing.T) {
	sequenceID := "seq-123"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/mail/sequences/"+sequenceID+"/entries")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListSequenceEntriesResponse{
			Entries: []SequenceEntry{
				{ID: "entry-1", ContactID: "contact-1", Status: SequenceEntryStatusActive},
			},
			Total: 1,
		})
	})
	defer server.Close()

	resp, err := sequencesClient.ListEntries(context.Background(), &ListSequenceEntriesRequest{
		ID: sequenceID,
	})

	require.NoError(t, err)
	assert.Len(t, resp.Entries, 1)
}

func TestSequencesClient_AddContact(t *testing.T) {
	sequenceID := "seq-123"
	contactID := "contact-456"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID+"/add-contact", r.URL.Path)

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, contactID, body["contactId"])

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SequenceEntry{
			ID:        "entry-new",
			ContactID: contactID,
			Status:    SequenceEntryStatusActive,
		})
	})
	defer server.Close()

	resp, err := sequencesClient.AddContact(context.Background(), &AddContactToSequenceRequest{
		ID:        sequenceID,
		ContactID: contactID,
	})

	require.NoError(t, err)
	assert.Equal(t, "entry-new", resp.ID)
	assert.Equal(t, contactID, resp.ContactID)
}

func TestSequencesClient_RemoveContact(t *testing.T) {
	sequenceID := "seq-123"
	entryID := "entry-456"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID+"/remove-contact", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(RemoveContactFromSequenceResponse{Success: true})
	})
	defer server.Close()

	resp, err := sequencesClient.RemoveContact(context.Background(), &RemoveContactFromSequenceRequest{
		ID:      sequenceID,
		EntryID: entryID,
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestSequencesClient_GetAnalytics(t *testing.T) {
	sequenceID := "seq-123"
	sequencesClient, server := setupSequencesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mail/sequences/"+sequenceID+"/analytics", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SequenceAnalyticsResponse{
			Sequence: struct {
				TotalEntered   int `json:"totalEntered"`
				TotalCompleted int `json:"totalCompleted"`
				TotalActive    int `json:"totalActive"`
			}{
				TotalEntered:   100,
				TotalCompleted: 80,
				TotalActive:    20,
			},
		})
	})
	defer server.Close()

	resp, err := sequencesClient.GetAnalytics(context.Background(), sequenceID)

	require.NoError(t, err)
	assert.Equal(t, 100, resp.Sequence.TotalEntered)
	assert.Equal(t, 80, resp.Sequence.TotalCompleted)
}

func TestSequenceStatus_Constants(t *testing.T) {
	assert.Equal(t, SequenceStatus("draft"), SequenceStatusDraft)
	assert.Equal(t, SequenceStatus("active"), SequenceStatusActive)
	assert.Equal(t, SequenceStatus("paused"), SequenceStatusPaused)
	assert.Equal(t, SequenceStatus("archived"), SequenceStatusArchived)
}

func TestSequenceTriggerType_Constants(t *testing.T) {
	assert.Equal(t, SequenceTriggerType("manual"), SequenceTriggerManual)
	assert.Equal(t, SequenceTriggerType("event_received"), SequenceTriggerEventReceived)
	assert.Equal(t, SequenceTriggerType("contact_added"), SequenceTriggerContactAdded)
	assert.Equal(t, SequenceTriggerType("api"), SequenceTriggerAPI)
	assert.Equal(t, SequenceTriggerType("scheduled"), SequenceTriggerScheduled)
}

func TestSequenceNodeType_Constants(t *testing.T) {
	assert.Equal(t, SequenceNodeType("trigger"), SequenceNodeTrigger)
	assert.Equal(t, SequenceNodeType("email"), SequenceNodeEmail)
	assert.Equal(t, SequenceNodeType("timer"), SequenceNodeTimer)
	assert.Equal(t, SequenceNodeType("filter"), SequenceNodeFilter)
	assert.Equal(t, SequenceNodeType("branch"), SequenceNodeBranch)
	assert.Equal(t, SequenceNodeType("experiment"), SequenceNodeExperiment)
	assert.Equal(t, SequenceNodeType("exit"), SequenceNodeExit)
}
