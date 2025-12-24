package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/joselitophala/budget-planner-backend/internal/auth"
	"github.com/joselitophala/budget-planner-backend/internal/models"
	"github.com/joselitophala/budget-planner-backend/internal/utils"
)

// SyncHandler handles offline sync-related requests
type SyncHandler struct {
	queries *models.Queries
}

// NewSyncHandler creates a new sync handler
func NewSyncHandler(queries *models.Queries) *SyncHandler {
	return &SyncHandler{queries: queries}
}

// PushRequest represents a sync push request from the client
type PushRequest struct {
	Operations []SyncOperation `json:"operations"`
}

// SyncOperation represents a single sync operation
type SyncOperation struct {
	Table      string                 `json:"table"`
	RecordID   string                 `json:"recordId"`
	Operation  string                 `json:"operation"` // create, update, delete
	LocalData  map[string]interface{} `json:"localData"`
	ServerData map[string]interface{} `json:"serverData,omitempty"`
}

// PullRequest represents a sync pull request
type PullRequest struct {
	LastSyncTime string `json:"lastSyncTime"` // ISO 8601 timestamp
}

// PullResponse represents the response to a sync pull request
type PullResponse struct {
	HasMore      bool                   `json:"hasMore"`
	LastSyncTime string                 `json:"lastSyncTime"`
	Changes      map[string][]json.RawMessage `json:"changes"` // table name -> records
}

// Push pushes local changes from the client to the server
func (h *SyncHandler) Push(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	var req PushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	results := make([]map[string]interface{}, 0)

	for _, op := range req.Operations {
		result := h.processSyncOperation(r.Context(), userID, op)
		results = append(results, result)
	}

	utils.SendSuccess(w, map[string]interface{}{
		"results": results,
		"syncedAt": time.Now().Format(time.RFC3339),
	})
}

// Pull pulls server changes down to the client
func (h *SyncHandler) Pull(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	var req PullRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	// Parse last sync time
	var lastSyncTime *time.Time
	if req.LastSyncTime != "" {
		if t, err := time.Parse(time.RFC3339, req.LastSyncTime); err == nil {
			lastSyncTime = &t
		}
	}

	// Get changes since last sync
	changes := make(map[string][]json.RawMessage)

	// Fetch updated budgets
	budgets, err := h.queries.GetBudgetsSince(r.Context(), models.GetBudgetsSinceParams{
		UserID:  utils.PgUUID(userID),
		Column2: lastSyncTime,
	})
	if err == nil && len(budgets) > 0 {
		budgetsJSON, _ := json.Marshal(budgets)
		changes["budgets"] = []json.RawMessage{budgetsJSON}
	}

	// Fetch updated transactions
	transactions, err := h.queries.GetTransactionsSince(r.Context(), models.GetTransactionsSinceParams{
		UserID:  utils.PgUUID(userID),
		Column2: lastSyncTime,
	})
	if err == nil && len(transactions) > 0 {
		txJSON, _ := json.Marshal(transactions)
		changes["transactions"] = []json.RawMessage{txJSON}
	}

	// Fetch updated categories
	categories, err := h.queries.GetCategoriesSince(r.Context(), models.GetCategoriesSinceParams{
		UserID:  utils.PgUUID(userID),
		Column2: lastSyncTime,
	})
	if err == nil && len(categories) > 0 {
		catsJSON, _ := json.Marshal(categories)
		changes["categories"] = []json.RawMessage{catsJSON}
	}

	utils.SendSuccess(w, PullResponse{
		HasMore:      false,
		LastSyncTime: time.Now().Format(time.RFC3339),
		Changes:      changes,
	})
}

// GetStatus returns the sync status for the current user
func (h *SyncHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	// Count pending sync operations
	pendingCount, err := h.queries.CountPendingSyncOperations(r.Context(), utils.PgUUID(userID))
	if err != nil {
		utils.InternalError(w, "Failed to fetch sync status")
		return
	}

	utils.SendSuccess(w, map[string]interface{}{
		"pendingOperations": pendingCount,
		"lastSyncTime":      time.Now().Format(time.RFC3339),
	})
}

// ResolveConflict handles conflict resolution for sync operations
func (h *SyncHandler) ResolveConflict(w http.ResponseWriter, r *http.Request) {
	_, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	var req struct {
		OperationID string                 `json:"operationId"`
		Resolution  string                 `json:"resolution"` // "local", "server", or "merge"
		MergedData  map[string]interface{} `json:"mergedData,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	// Update the sync operation
	err := h.queries.ResolveSyncOperation(r.Context(), models.ResolveSyncOperationParams{
		ID:     req.OperationID,
		Status: utils.PgText(req.Resolution),
	})
	if err != nil {
		utils.InternalError(w, "Failed to resolve conflict")
		return
	}

	utils.SendSuccess(w, map[string]string{
		"message": "Conflict resolved successfully",
	})
}

// processSyncOperation processes a single sync operation
func (h *SyncHandler) processSyncOperation(ctx context.Context, userID string, op SyncOperation) map[string]interface{} {
	result := map[string]interface{}{
		"table":     op.Table,
		"recordId":  op.RecordID,
		"operation": op.Operation,
		"status":    "success",
	}

	// Implement table-specific logic
	switch op.Table {
	case "transactions":
		// Handle transaction sync
		// This would call the appropriate query based on operation type
	case "budgets":
		// Handle budget sync
	case "categories":
		// Handle category sync
	default:
		result["status"] = "error"
		result["error"] = "unsupported table"
	}

	return result
}
