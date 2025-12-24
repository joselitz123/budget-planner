package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joselitophala/budget-planner-backend/internal/models"
	"github.com/joselitophala/budget-planner-backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReflectionsHandler_GetReflectionByMonth tests getting reflection by month
func TestReflectionsHandler_GetReflectionByMonth(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewReflectionHandler(TestQueries)

	t.Run("Get reflection returns not found initially", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/reflections/month/2025-01", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.GetReflectionByMonth(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Get reflection by month", func(t *testing.T) {
		// Create a test budget and reflection
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)
		reflectionID := CreateTestReflection(t, ctx, userID, budgetID, 8)

		req := httptest.NewRequest("GET", "/api/reflections/month/2025-01", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.GetReflectionByMonth(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    ReflectionResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, reflectionID, response.Data.ID)

		// Cleanup
		TestQueries.DeleteReflection(ctx, reflectionID)
		TestQueries.DeleteBudget(ctx, budgetID)
	})
}

// TestReflectionsHandler_CreateReflection tests creating a reflection
func TestReflectionsHandler_CreateReflection(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewReflectionHandler(TestQueries)

	t.Run("Create reflection with valid data", func(t *testing.T) {
		// Create a test budget
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)

		reqBody := map[string]interface{}{
			"budgetId":       budgetID,
			"overallRating":  8,
			"isPrivate":      true,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/reflections", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.CreateReflection(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    ReflectionResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)

		// Cleanup
		TestQueries.DeleteReflection(ctx, response.Data.ID)
		TestQueries.DeleteBudget(ctx, budgetID)
	})

	t.Run("Create reflection with invalid rating", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"budgetId":       "test-budget-id",
			"overallRating":  15, // Invalid: must be 1-10
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/reflections", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.CreateReflection(w, req)

		// Should return error
		assert.NotEqual(t, http.StatusCreated, w.Code)
	})
}

// TestReflectionsHandler_UpdateReflection tests updating a reflection
func TestReflectionsHandler_UpdateReflection(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewReflectionHandler(TestQueries)

	t.Run("Update reflection with valid data", func(t *testing.T) {
		// Create a test budget and reflection
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)
		reflectionID := CreateTestReflection(t, ctx, userID, budgetID, 7)

		newRating := int32(9)
		reqBody := map[string]interface{}{
			"overallRating": newRating,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/api/reflections/"+reflectionID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.UpdateReflection(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    ReflectionResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)

		// Cleanup
		TestQueries.DeleteReflection(ctx, reflectionID)
		TestQueries.DeleteBudget(ctx, budgetID)
	})
}

// TestReflectionsHandler_DeleteReflection tests deleting a reflection
func TestReflectionsHandler_DeleteReflection(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewReflectionHandler(TestQueries)

	t.Run("Delete reflection successfully", func(t *testing.T) {
		// Create a test budget and reflection
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)
		reflectionID := CreateTestReflection(t, ctx, userID, budgetID, 8)

		req := httptest.NewRequest("DELETE", "/api/reflections/"+reflectionID, nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.DeleteReflection(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    map[string]string `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)

		// Cleanup
		TestQueries.DeleteBudget(ctx, budgetID)
	})
}

// TestReflectionsHandler_ListReflectionTemplates tests listing reflection templates
func TestReflectionsHandler_ListReflectionTemplates(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewReflectionHandler(TestQueries)

	t.Run("List reflection templates", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/reflections/templates", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.ListReflectionTemplates(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    []models.ReflectionTemplate `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)
	})
}

// Helper function to create a test reflection
func CreateTestReflection(t *testing.T, ctx context.Context, userID, budgetID string, rating int32) string {
	t.Helper()

	reflection, err := TestQueries.CreateReflection(ctx, models.CreateReflectionParams{
		UserID:        utils.PgUUID(userID),
		BudgetID:      utils.PgUUID(budgetID),
		OverallRating: utils.PgInt4(rating),
		IsPrivate:     utils.PgBool(true),
	})
	if err != nil {
		t.Fatalf("Failed to create test reflection: %v", err)
	}
	return reflection.ID
}
