package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joselitophala/budget-planner-backend/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMain runs once before all tests
func TestMain(m *testing.M) {
	// Setup test database
	SetupTestDB(&testing.T{})

	// Run tests
	code := m.Run()

	// Cleanup
	CleanupTestDB(&testing.T{})
	os.Exit(code)
}

// TestAuthHandler_Login tests the Login endpoint
func TestAuthHandler_Login(t *testing.T) {
	ctx := context.Background()

	// Create test user
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	// Create handler
	h := NewAuthHandler(TestQueries, TestJWTClient)

	t.Run("Login with valid token returns user data", func(t *testing.T) {
		// Generate valid token
		token := GenerateTestToken(t, TestClerkID)

		// Create request
		reqBody := LoginRequest{Token: token}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call handler
		h.Login(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    struct {
				User UserResponse `json:"user"`
			} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, userID, response.Data.User.ID)
		assert.Equal(t, "Test User", response.Data.User.Name)
	})

	t.Run("Login with invalid token returns error", func(t *testing.T) {
		reqBody := LoginRequest{Token: "invalid-token"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Login(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Login for non-existent user returns onboarding prompt", func(t *testing.T) {
		// Generate token for unknown user
		token := GenerateTestToken(t, "unknown_clerk_user")

		reqBody := LoginRequest{Token: token}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Login(w, req)

		assert.Equal(t, http.StatusAccepted, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    struct {
				NeedsOnboarding bool   `json:"needsOnboarding"`
				ClerkUserID     string `json:"clerkUserID"`
			} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Data.NeedsOnboarding)
	})
}

// TestAuthHandler_CompleteOnboarding tests the CompleteOnboarding endpoint
func TestAuthHandler_CompleteOnboarding(t *testing.T) {
	ctx := context.Background()
	h := NewAuthHandler(TestQueries, TestJWTClient)

	t.Run("Complete onboarding creates new user", func(t *testing.T) {
		clerkID := "new_user_clerk_123"

		reqBody := OnboardingRequest{
			ClerkUserID: clerkID,
			Name:        "New User",
			Currency:    "EUR",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/auth/onboarding", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.CompleteOnboarding(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    UserResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "New User", response.Data.Name)
		assert.Equal(t, "EUR", response.Data.Currency)

		// Cleanup
		TestQueries.DeleteUser(ctx, response.Data.ID)
	})
}

// TestAuthHandler_Logout tests the Logout endpoint
func TestAuthHandler_Logout(t *testing.T) {
	h := NewAuthHandler(TestQueries, TestJWTClient)

	req := httptest.NewRequest("POST", "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	h.Logout(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			Message string `json:"message"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Logged out successfully", response.Data.Message)
}

// TestAuthHandler_RefreshToken tests the RefreshToken endpoint
func TestAuthHandler_RefreshToken(t *testing.T) {
	h := NewAuthHandler(TestQueries, TestJWTClient)

	req := httptest.NewRequest("POST", "/api/auth/refresh", nil)
	w := httptest.NewRecorder()

	h.RefreshToken(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			Message string `json:"message"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Success)
}

// TestAuthHandler_GetCurrentUser tests the GetCurrentUser endpoint
func TestAuthHandler_GetCurrentUser(t *testing.T) {
	ctx := context.Background()

	// Create test user
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewAuthHandler(TestQueries, TestJWTClient)

	t.Run("Get current user with valid auth context", func(t *testing.T) {
		// Create request with auth context
		req := httptest.NewRequest("GET", "/api/auth/me", nil)
		w := httptest.NewRecorder()

		// Add user context (simulating auth middleware)
		reqContext := auth.SetUserIDInContext(req, userID)
		req = req.WithContext(reqContext)

		h.GetCurrentUser(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    UserResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, userID, response.Data.ID)
	})

	t.Run("Get current user without auth returns unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/auth/me", nil)
		w := httptest.NewRecorder()

		h.GetCurrentUser(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
