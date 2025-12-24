package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joselitophala/budget-planner-backend/internal/auth"
	"github.com/joselitophala/budget-planner-backend/internal/models"
	"github.com/joselitophala/budget-planner-backend/internal/utils"
)

// UserHandler handles user-related requests
type UserHandler struct {
	queries *models.Queries
}

// NewUserHandler creates a new user handler
func NewUserHandler(queries *models.Queries) *UserHandler {
	return &UserHandler{queries: queries}
}

// UpdateUserRequest represents the update user request body
type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	Currency *string `json:"currency,omitempty"`
}

// GetProfile returns the current user's profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	user, err := h.queries.GetCurrentUser(r.Context(), userID)
	if err != nil {
		utils.NotFound(w, "User not found")
		return
	}

	utils.SendSuccess(w, UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Name:     utils.TextToString(user.Name),
		Currency: utils.TextToString(user.Currency),
	})
}

// UpdateProfile updates the current user's profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	// Convert *string to pgtype.Text
	var name, currency pgtype.Text
	if req.Name != nil {
		name = utils.PgText(*req.Name)
	}
	if req.Currency != nil {
		currency = utils.PgText(*req.Currency)
	}

	user, err := h.queries.UpdateUser(r.Context(), models.UpdateUserParams{
		ID:       userID,
		Name:     name,
		Currency: currency,
	})
	if err != nil {
		utils.InternalError(w, "Failed to update user")
		return
	}

	utils.SendSuccess(w, UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Name:     utils.TextToString(user.Name),
		Currency: utils.TextToString(user.Currency),
	})
}

// DeleteAccount deletes the current user's account
func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	// Soft delete the user
	err := h.queries.DeleteUser(r.Context(), userID)
	if err != nil {
		utils.InternalError(w, "Failed to delete account")
		return
	}

	utils.SendSuccess(w, map[string]string{
		"message": "Account deleted successfully",
	})
}
