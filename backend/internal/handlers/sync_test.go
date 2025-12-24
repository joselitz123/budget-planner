package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSyncHandler_Push tests pushing sync data
func TestSyncHandler_Push(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewSyncHandler(TestQueries)

	t.Run("Push sync operations", func(t *testing.T) {
		syncData := []map[string]interface{}{
			{
				"tableName":  "categories",
				"recordId":   "test-id-123",
				"operation":  "create",
				"localData":  json.RawMessage(`{"name":"Test","color":"#FF5733"}`),
				"serverData": nil,
			},
		}

		reqBody := map[string]interface{}{
			"operations": syncData,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/sync/push", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.Push(w, req)

		// Check response
		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("Push with invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/sync/push", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.Push(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestSyncHandler_Pull tests pulling sync data
func TestSyncHandler_Pull(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewSyncHandler(TestQueries)

	t.Run("Pull sync data", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"lastSyncAt": "2024-01-01T00:00:00Z",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/sync/pull", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.Pull(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
	})
}

// TestSyncHandler_GetStatus tests getting sync status
func TestSyncHandler_GetStatus(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewSyncHandler(TestQueries)

	t.Run("Get sync status", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/sync/status", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.GetStatus(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
	})
}

// TestSyncHandler_ResolveConflict tests resolving sync conflicts
func TestSyncHandler_ResolveConflict(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewSyncHandler(TestQueries)

	t.Run("Resolve conflict with local version", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"syncId":     "test-sync-id",
			"resolution": "local",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/sync/resolve-conflict", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.ResolveConflict(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("Resolve conflict with invalid resolution", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"syncId":     "test-sync-id",
			"resolution": "invalid",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/sync/resolve-conflict", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.ResolveConflict(w, req)

		// Should return error
		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}
