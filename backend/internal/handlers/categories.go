package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joselitophala/budget-planner-backend/internal/auth"
	"github.com/joselitophala/budget-planner-backend/internal/models"
	"github.com/joselitophala/budget-planner-backend/internal/utils"
)

// CategoryHandler handles category-related requests
type CategoryHandler struct {
	queries *models.Queries
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(queries *models.Queries) *CategoryHandler {
	return &CategoryHandler{queries: queries}
}

// CategoryResponse represents a category in API responses
type CategoryResponse struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Icon         *string  `json:"icon,omitempty"`
	Color        string   `json:"color"`
	IsSystem     bool     `json:"isSystem"`
	DefaultLimit *float64 `json:"defaultLimit,omitempty"`
}

// CreateCategoryRequest represents the create category request
type CreateCategoryRequest struct {
	Name         string  `json:"name"`
	Icon         *string `json:"icon,omitempty"`
	Color        string  `json:"color"`
	DefaultLimit *float64 `json:"defaultLimit,omitempty"`
}

// UpdateCategoryRequest represents the update category request
type UpdateCategoryRequest struct {
	Name         *string `json:"name,omitempty"`
	Icon         *string `json:"icon,omitempty"`
	Color        *string `json:"color,omitempty"`
	DefaultLimit *float64 `json:"defaultLimit,omitempty"`
}

// ListCategories returns all categories for the current user
func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	categories, err := h.queries.GetUserCategories(r.Context(), utils.PgUUID(userID))
	if err != nil {
		utils.InternalError(w, "Failed to fetch categories")
		return
	}

	response := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		response[i] = CategoryResponse{
			ID:           cat.ID,
			Name:         cat.Name,
			Icon:         utils.TextToStringPtr(cat.Icon),
			Color:        utils.TextToString(cat.Color),
			IsSystem:     cat.IsSystem.Valid,
			DefaultLimit: utils.NumericToFloat64Ptr(cat.DefaultLimit),
		}
	}

	utils.SendSuccess(w, response)
}

// GetSystemCategories returns all system default categories
func (h *CategoryHandler) GetSystemCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.queries.GetSystemCategories(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to fetch system categories")
		return
	}

	response := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		response[i] = CategoryResponse{
			ID:           cat.ID,
			Name:         cat.Name,
			Icon:         utils.TextToStringPtr(cat.Icon),
			Color:        utils.TextToString(cat.Color),
			IsSystem:     cat.IsSystem.Valid,
			DefaultLimit: utils.NumericToFloat64Ptr(cat.DefaultLimit),
		}
	}

	utils.SendSuccess(w, response)
}

// CreateCategory creates a new custom category
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	var icon pgtype.Text
	if req.Icon != nil {
		icon = utils.PgText(*req.Icon)
	}

	var defaultLimit pgtype.Numeric
	if req.DefaultLimit != nil {
		defaultLimit = utils.PgNumeric(*req.DefaultLimit)
	}

	category, err := h.queries.CreateCategory(r.Context(), models.CreateCategoryParams{
		UserID:       utils.PgUUID(userID),
		Name:         req.Name,
		Icon:         icon,
		Color:        utils.PgText(req.Color),
		IsSystem:     pgtype.Bool{Valid: true, Bool: false},
		DefaultLimit: defaultLimit,
	})
	if err != nil {
		utils.InternalError(w, "Failed to create category")
		return
	}

	utils.SendCreated(w, CategoryResponse{
		ID:           category.ID,
		Name:         category.Name,
		Icon:         utils.TextToStringPtr(category.Icon),
		Color:        utils.TextToString(category.Color),
		IsSystem:     category.IsSystem.Valid,
		DefaultLimit: utils.NumericToFloat64Ptr(category.DefaultLimit),
	})
}

// UpdateCategory updates an existing category
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	_, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	categoryID := r.PathValue("id")
	if categoryID == "" {
		utils.BadRequest(w, "Category ID is required")
		return
	}

	var req UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	var name pgtype.Text
	if req.Name != nil {
		name = utils.PgText(*req.Name)
	}

	var icon pgtype.Text
	if req.Icon != nil {
		icon = utils.PgText(*req.Icon)
	}

	var color pgtype.Text
	if req.Color != nil {
		color = utils.PgText(*req.Color)
	}

	var defaultLimit pgtype.Numeric
	if req.DefaultLimit != nil {
		defaultLimit = utils.PgNumeric(*req.DefaultLimit)
	}

	category, err := h.queries.UpdateCategory(r.Context(), models.UpdateCategoryParams{
		ID:           categoryID,
		Name:         name,
		Icon:         icon,
		Color:        color,
		DefaultLimit: defaultLimit,
	})
	if err != nil {
		utils.InternalError(w, "Failed to update category")
		return
	}

	utils.SendSuccess(w, CategoryResponse{
		ID:           category.ID,
		Name:         category.Name,
		Icon:         utils.TextToStringPtr(category.Icon),
		Color:        utils.TextToString(category.Color),
		IsSystem:     category.IsSystem.Valid,
		DefaultLimit: utils.NumericToFloat64Ptr(category.DefaultLimit),
	})
}

// DeleteCategory soft deletes a category
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	_, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	categoryID := r.PathValue("id")
	if categoryID == "" {
		utils.BadRequest(w, "Category ID is required")
		return
	}

	err := h.queries.DeleteCategory(r.Context(), categoryID)
	if err != nil {
		utils.InternalError(w, "Failed to delete category")
		return
	}

	utils.SendSuccess(w, map[string]string{
		"message": "Category deleted successfully",
	})
}
