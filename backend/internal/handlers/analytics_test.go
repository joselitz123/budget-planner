package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAnalyticsHandler_GetDashboard tests getting dashboard data
func TestAnalyticsHandler_GetDashboard(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewAnalyticsHandler(TestQueries)

	t.Run("Get dashboard for month", func(t *testing.T) {
		// Create a test budget
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)

		req := httptest.NewRequest("GET", "/api/analytics/dashboard/2025-01", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.GetDashboard(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)

		// Cleanup
		TestQueries.DeleteBudget(ctx, budgetID)
	})
}

// TestAnalyticsHandler_GetSpendingReport tests getting spending report
func TestAnalyticsHandler_GetSpendingReport(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewAnalyticsHandler(TestQueries)

	t.Run("Get spending report for month", func(t *testing.T) {
		// Create test budget and category
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)
		categoryID := CreateTestCategory(t, ctx, userID, "Food", "🍔", "#FF5733")

		req := httptest.NewRequest("GET", "/api/analytics/spending/2025-01", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.GetSpendingReport(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)

		// Cleanup
		TestQueries.DeleteBudget(ctx, budgetID)
		TestQueries.DeleteCategory(ctx, categoryID)
	})
}

// TestAnalyticsHandler_GetTrends tests getting spending trends
func TestAnalyticsHandler_GetTrends(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewAnalyticsHandler(TestQueries)

	t.Run("Get spending trends", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/analytics/trends?months=6", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.GetTrends(w, req)

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

// TestAnalyticsHandler_GetCategoryReport tests getting category-specific report
func TestAnalyticsHandler_GetCategoryReport(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewAnalyticsHandler(TestQueries)

	t.Run("Get category report", func(t *testing.T) {
		// Create a test category
		categoryID := CreateTestCategory(t, ctx, userID, "Food", "🍔", "#FF5733")

		req := httptest.NewRequest("GET", "/api/analytics/category/"+categoryID+"?months=3", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.GetCategoryReport(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)

		// Cleanup
		TestQueries.DeleteCategory(ctx, categoryID)
	})

	t.Run("Get category report with invalid months parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/analytics/category/test-id?months=invalid", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.GetCategoryReport(w, req)

		// Should return error
		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}
