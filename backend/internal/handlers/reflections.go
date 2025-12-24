package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joselitophala/budget-planner-backend/internal/auth"
	"github.com/joselitophala/budget-planner-backend/internal/models"
	"github.com/joselitophala/budget-planner-backend/internal/utils"
)

// ReflectionHandler handles reflection-related requests
type ReflectionHandler struct {
	queries *models.Queries
}

// NewReflectionHandler creates a new reflection handler
func NewReflectionHandler(queries *models.Queries) *ReflectionHandler {
	return &ReflectionHandler{queries: queries}
}

// ReflectionResponse represents a reflection in API responses
type ReflectionResponse struct {
	ID            string             `json:"id"`
	BudgetID      string             `json:"budgetId"`
	OverallRating *int32             `json:"overallRating,omitempty"`
	IsPrivate     bool               `json:"isPrivate"`
	CreatedAt     string             `json:"createdAt"`
	UpdatedAt     string             `json:"updatedAt"`
}

// CreateReflectionRequest represents the create reflection request
type CreateReflectionRequest struct {
	BudgetID      string             `json:"budgetId"`
	OverallRating *int32             `json:"overallRating,omitempty"`
	IsPrivate     bool               `json:"isPrivate"`
}

// GetReflectionByMonth returns a reflection for a specific month
func (h *ReflectionHandler) GetReflectionByMonth(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	monthStr := r.PathValue("month")
	if monthStr == "" {
		utils.BadRequest(w, "Month is required (format: YYYY-MM)")
		return
	}

	// Parse month
	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		utils.BadRequest(w, "Invalid month format. Use YYYY-MM")
		return
	}

	// Get budget for this month
	budget, err := h.queries.GetBudgetByMonth(r.Context(), models.GetBudgetByMonthParams{
		UserID: utils.PgUUID(userID),
		Month:  utils.PgDate(month),
	})
	if err != nil {
		utils.NotFound(w, "Budget not found for this month")
		return
	}

	// Get reflection for this budget
	reflection, err := h.queries.GetReflectionByBudget(r.Context(), utils.PgUUID(budget.ID))
	if err != nil {
		utils.SendSuccess(w, nil) // No reflection yet
		return
	}

	utils.SendSuccess(w, reflectionToResponse(reflection))
}

// CreateReflection creates a new reflection
func (h *ReflectionHandler) CreateReflection(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	var req CreateReflectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	// Verify the budget belongs to the user
	budget, err := h.queries.GetBudgetByID(r.Context(), req.BudgetID)
	if err != nil || !budget.UserID.Valid || utils.UUIDToString(budget.UserID) != userID {
		utils.Forbidden(w, "You can only create reflections for your own budgets")
		return
	}

	// Handle OverallRating conversion
	var rating pgtype.Int4
	if req.OverallRating != nil {
		rating = utils.PgInt4(*req.OverallRating)
	}

	reflection, err := h.queries.CreateReflection(r.Context(), models.CreateReflectionParams{
		UserID:        utils.PgUUID(userID),
		BudgetID:      utils.PgUUID(req.BudgetID),
		OverallRating: rating,
		IsPrivate:     utils.PgBool(req.IsPrivate),
	})
	if err != nil {
		utils.InternalError(w, "Failed to create reflection")
		return
	}

	utils.SendCreated(w, reflectionToResponse(reflection))
}

// UpdateReflection updates an existing reflection
func (h *ReflectionHandler) UpdateReflection(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	reflectionID := r.PathValue("id")
	if reflectionID == "" {
		utils.BadRequest(w, "Reflection ID is required")
		return
	}

	var req UpdateReflectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	// Verify ownership
	existing, err := h.queries.GetReflectionByID(r.Context(), reflectionID)
	if err != nil || !existing.UserID.Valid || utils.UUIDToString(existing.UserID) != userID {
		utils.Forbidden(w, "You can only update your own reflections")
		return
	}

	reflection, err := h.queries.UpdateReflection(r.Context(), models.UpdateReflectionParams{
		ID:            reflectionID,
		OverallRating: utils.PgInt4Ptr(req.OverallRating),
		IsPrivate:     utils.PgBoolPtr(req.IsPrivate),
	})
	if err != nil {
		utils.InternalError(w, "Failed to update reflection")
		return
	}

	utils.SendSuccess(w, reflectionToResponse(reflection))
}

// DeleteReflection deletes a reflection
func (h *ReflectionHandler) DeleteReflection(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	reflectionID := r.PathValue("id")
	if reflectionID == "" {
		utils.BadRequest(w, "Reflection ID is required")
		return
	}

	// Verify ownership
	existing, err := h.queries.GetReflectionByID(r.Context(), reflectionID)
	if err != nil || !existing.UserID.Valid || utils.UUIDToString(existing.UserID) != userID {
		utils.Forbidden(w, "You can only delete your own reflections")
		return
	}

	err = h.queries.DeleteReflection(r.Context(), reflectionID)
	if err != nil {
		utils.InternalError(w, "Failed to delete reflection")
		return
	}

	utils.SendSuccess(w, map[string]string{
		"message": "Reflection deleted successfully",
	})
}

// ListReflectionTemplates returns all reflection templates
func (h *ReflectionHandler) ListReflectionTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.queries.ListReflectionTemplates(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to fetch templates")
		return
	}

	utils.SendSuccess(w, templates)
}

type UpdateReflectionRequest struct {
	OverallRating *int32 `json:"overallRating,omitempty"`
	IsPrivate     *bool  `json:"isPrivate,omitempty"`
}

func reflectionToResponse(r models.Reflection) ReflectionResponse {
	return ReflectionResponse{
		ID:            r.ID,
		BudgetID:      utils.UUIDToString(r.BudgetID),
		OverallRating: utils.Int4ToInt32(r.OverallRating),
		IsPrivate:     r.IsPrivate.Bool,
		CreatedAt:     utils.TimestamptzToTime(r.CreatedAt).Format(time.RFC3339),
		UpdatedAt:     utils.TimestamptzToTime(r.UpdatedAt).Format(time.RFC3339),
	}
}
