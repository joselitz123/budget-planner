package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/joselitophala/budget-planner-backend/internal/auth"
	"github.com/joselitophala/budget-planner-backend/internal/models"
	"github.com/joselitophala/budget-planner-backend/internal/utils"
)

// SharingHandler handles budget sharing-related requests
type SharingHandler struct {
	queries *models.Queries
}

// NewSharingHandler creates a new sharing handler
func NewSharingHandler(queries *models.Queries) *SharingHandler {
	return &SharingHandler{queries: queries}
}

// ShareInvitationRequest represents the create share invitation request
type ShareInvitationRequest struct {
	BudgetID    string `json:"budgetId"`
	RecipientEmail string `json:"recipientEmail"`
	Permission  string `json:"permission"` // "view" or "edit"
}

// UpdateInvitationRequest represents the update invitation request
type UpdateInvitationRequest struct {
	Status string `json:"status"` // "accepted" or "declined"
}

// CreateShareInvitation creates a new share invitation
func (h *SharingHandler) CreateShareInvitation(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	var req ShareInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	// Verify ownership of the budget
	budget, err := h.queries.GetBudgetByID(r.Context(), req.BudgetID)
	if err != nil || !budget.UserID.Valid || utils.UUIDToString(budget.UserID) != userID {
		utils.Forbidden(w, "You can only share your own budgets")
		return
	}

	// Validate permission
	if req.Permission != "view" && req.Permission != "edit" {
		utils.BadRequest(w, "Permission must be 'view' or 'edit'")
		return
	}

	invitation, err := h.queries.CreateShareInvitation(r.Context(), models.CreateShareInvitationParams{
		BudgetID:       utils.PgUUID(req.BudgetID),
		OwnerID:        utils.PgUUID(userID),
		RecipientEmail: req.RecipientEmail,
		Permission:     req.Permission,
		ExpiresAt:      utils.PgTimestamptz(time.Now().Add(7 * 24 * time.Hour)), // 7 days
	})
	if err != nil {
		utils.InternalError(w, "Failed to create invitation")
		return
	}

	utils.SendCreated(w, invitation)
}

// GetMyInvitations returns pending invitations for the current user
func (h *SharingHandler) GetMyInvitations(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	// Get user email to check for invitations
	user, err := h.queries.GetCurrentUser(r.Context(), userID)
	if err != nil {
		utils.InternalError(w, "Failed to fetch user")
		return
	}

	invitations, err := h.queries.GetPendingInvitationsByRecipient(r.Context(), user.Email)
	if err != nil {
		utils.InternalError(w, "Failed to fetch invitations")
		return
	}

	utils.SendSuccess(w, invitations)
}

// RespondToInvitation accepts or declines a share invitation
func (h *SharingHandler) RespondToInvitation(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	invitationID := r.PathValue("id")
	if invitationID == "" {
		utils.BadRequest(w, "Invitation ID is required")
		return
	}

	var req UpdateInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	// Get the invitation
	invitation, err := h.queries.GetInvitationByID(r.Context(), invitationID)
	if err != nil {
		utils.NotFound(w, "Invitation not found")
		return
	}

	// Get user email to verify
	user, err := h.queries.GetCurrentUser(r.Context(), userID)
	if err != nil || user.Email != invitation.RecipientEmail {
		utils.Forbidden(w, "This invitation is not for you")
		return
	}

	// Update invitation status
	updated, err := h.queries.UpdateInvitationStatus(r.Context(), models.UpdateInvitationStatusParams{
		ID:     invitationID,
		Status: utils.PgText(req.Status),
	})
	if err != nil {
		utils.InternalError(w, "Failed to update invitation")
		return
	}

	// If accepted, create share access record
	if req.Status == "accepted" {
		_, err = h.queries.CreateShareAccess(r.Context(), models.CreateShareAccessParams{
			BudgetID:     updated.BudgetID,
			OwnerID:      updated.OwnerID,
			SharedWithID: utils.PgUUID(userID),
			Permission:   updated.Permission,
		})
		if err != nil {
			utils.InternalError(w, "Failed to create share access")
			return
		}
	}

	utils.SendSuccess(w, updated)
}

// CancelInvitation cancels a pending invitation (owner only)
func (h *SharingHandler) CancelInvitation(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	invitationID := r.PathValue("id")
	if invitationID == "" {
		utils.BadRequest(w, "Invitation ID is required")
		return
	}

	// Get the invitation to verify ownership
	invitation, err := h.queries.GetInvitationByID(r.Context(), invitationID)
	if err != nil {
		utils.NotFound(w, "Invitation not found")
		return
	}

	if !invitation.OwnerID.Valid || utils.UUIDToString(invitation.OwnerID) != userID {
		utils.Forbidden(w, "You can only cancel your own invitations")
		return
	}

	err = h.queries.DeleteInvitation(r.Context(), invitationID)
	if err != nil {
		utils.InternalError(w, "Failed to cancel invitation")
		return
	}

	utils.SendSuccess(w, map[string]string{
		"message": "Invitation cancelled successfully",
	})
}

// GetBudgetSharing returns who has access to a budget
func (h *SharingHandler) GetBudgetSharing(w http.ResponseWriter, r *http.Request) {
	budgetID := r.PathValue("budgetId")
	if budgetID == "" {
		utils.BadRequest(w, "Budget ID is required")
		return
	}

	accessList, err := h.queries.GetShareAccessByBudget(r.Context(), utils.PgUUID(budgetID))
	if err != nil {
		utils.InternalError(w, "Failed to fetch sharing info")
		return
	}

	utils.SendSuccess(w, accessList)
}

// RemoveAccess removes someone's access to a budget (owner only)
func (h *SharingHandler) RemoveAccess(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	accessID := r.PathValue("id")
	if accessID == "" {
		utils.BadRequest(w, "Access ID is required")
		return
	}

	// Get the access record to verify ownership
	access, err := h.queries.GetShareAccessByID(r.Context(), accessID)
	if err != nil {
		utils.NotFound(w, "Access record not found")
		return
	}

	if !access.OwnerID.Valid || utils.UUIDToString(access.OwnerID) != userID {
		utils.Forbidden(w, "You can only remove access for your own budgets")
		return
	}

	err = h.queries.DeleteShareAccess(r.Context(), accessID)
	if err != nil {
		utils.InternalError(w, "Failed to remove access")
		return
	}

	utils.SendSuccess(w, map[string]string{
		"message": "Access removed successfully",
	})
}

// GetSharedBudgets returns budgets shared with the current user
func (h *SharingHandler) GetSharedBudgets(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	sharedBudgets, err := h.queries.GetShareAccessForUser(r.Context(), utils.PgUUID(userID))
	if err != nil {
		utils.InternalError(w, "Failed to fetch shared budgets")
		return
	}

	utils.SendSuccess(w, sharedBudgets)
}
