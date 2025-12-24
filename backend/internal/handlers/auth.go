package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/joselitophala/budget-planner-backend/internal/auth"
	"github.com/joselitophala/budget-planner-backend/internal/models"
	"github.com/joselitophala/budget-planner-backend/internal/utils"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	queries *models.Queries
	jwt     *auth.JWTClient
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(queries *models.Queries, jwt *auth.JWTClient) *AuthHandler {
	return &AuthHandler{
		queries: queries,
		jwt:     jwt,
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Token string `json:"token"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	User *UserResponse `json:"user"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

// Login handles user login via JWT token
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	// Verify the JWT token
	clerkUserID, err := h.jwt.VerifyToken(req.Token)
	if err != nil {
		utils.Unauthorized(w, "Invalid token")
		return
	}

	// Look up user in our database
	user, err := h.queries.GetUserByClerkID(r.Context(), clerkUserID)
	if err != nil {
		// User exists in auth system but not in our database
		// They need to complete onboarding
		utils.SendJSON(w, http.StatusAccepted, map[string]interface{}{
			"needsOnboarding": true,
			"clerkUserID":     clerkUserID,
		})
		return
	}

	// Return user data
	utils.SendSuccess(w, LoginResponse{
		User: &UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Name:     utils.TextToString(user.Name),
			Currency: utils.TextToString(user.Currency),
		},
	})
}

// OnboardingRequest represents the onboarding request
type OnboardingRequest struct {
	ClerkUserID string `json:"clerkUserId"`
	Name        string `json:"name"`
	Currency    string `json:"currency"`
}

// CompleteOnboarding creates a user record after auth signup
func (h *AuthHandler) CompleteOnboarding(w http.ResponseWriter, r *http.Request) {
	var req OnboardingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	// Create user in our database
	// For now, we'll use a placeholder email since we don't have Clerk to fetch it
	user, err := h.queries.CreateUser(r.Context(), models.CreateUserParams{
		ClerkUserID: req.ClerkUserID,
		Email:       req.ClerkUserID + "@placeholder.local", // Placeholder
		Name:        utils.PgText(req.Name),
		Currency:    utils.PgText(req.Currency),
	})
	if err != nil {
		utils.InternalError(w, "Failed to create user")
		return
	}

	utils.SendCreated(w, UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Name:     utils.TextToString(user.Name),
		Currency: utils.TextToString(user.Currency),
	})
}

// GetCurrentUser returns the currently authenticated user
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
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

// Logout handles user logout
// Note: Auth tokens are managed on the frontend, so this is mainly for logging
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	utils.SendSuccess(w, map[string]string{
		"message": "Logged out successfully",
	})
}

// RefreshToken handles token refresh
// Note: Token refresh is handled on the frontend
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	utils.SendSuccess(w, map[string]string{
		"message": "Token refresh handled on frontend",
	})
}
