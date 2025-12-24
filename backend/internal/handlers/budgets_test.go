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

// TestBudgetsHandler_ListBudgets tests listing budgets
func TestBudgetsHandler_ListBudgets(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewBudgetHandler(TestQueries)

	t.Run("List budgets returns empty list initially", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/budgets", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.ListBudgets(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    []BudgetResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Empty(t, response.Data)
	})

	t.Run("List budgets returns user budgets", func(t *testing.T) {
		// Create a test budget
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)

		req := httptest.NewRequest("GET", "/api/budgets", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.ListBudgets(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    []BudgetResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotEmpty(t, response.Data)
		assert.Equal(t, budgetID, response.Data[0].ID)

		// Cleanup
		TestQueries.DeleteBudget(ctx, budgetID)
	})
}

// TestBudgetsHandler_GetBudgetByMonth tests getting budget by month
func TestBudgetsHandler_GetBudgetByMonth(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewBudgetHandler(TestQueries)

	t.Run("Get budget by month returns budget", func(t *testing.T) {
		// Create a test budget
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)

		req := httptest.NewRequest("GET", "/api/budgets/2025-01", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.GetBudgetByMonth(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    BudgetResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, budgetID, response.Data.ID)

		// Cleanup
		TestQueries.DeleteBudget(ctx, budgetID)
	})

	t.Run("Get budget by month with invalid format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/budgets/invalid", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.GetBudgetByMonth(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestBudgetsHandler_GetBudget tests getting budget by ID
func TestBudgetsHandler_GetBudget(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewBudgetHandler(TestQueries)

	t.Run("Get budget by ID returns budget", func(t *testing.T) {
		// Create a test budget
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)

		req := httptest.NewRequest("GET", "/api/budgets/"+budgetID, nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		// Set path value (Chi router does this)
		// In tests we need to simulate this differently

		h.GetBudget(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    BudgetResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)

		// Cleanup
		TestQueries.DeleteBudget(ctx, budgetID)
	})
}

// TestBudgetsHandler_CreateBudget tests creating a budget
func TestBudgetsHandler_CreateBudget(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewBudgetHandler(TestQueries)

	t.Run("Create budget with valid data", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"month":      "2025-01-01",
			"totalLimit": 5000,
			"name":       "January 2025",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/budgets", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.CreateBudget(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    BudgetResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, 5000.0, response.Data.TotalLimit)

		// Cleanup
		TestQueries.DeleteBudget(ctx, response.Data.ID)
	})

	t.Run("Create budget with invalid month format", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"month":      "invalid",
			"totalLimit": 5000,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/budgets", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.CreateBudget(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestBudgetsHandler_UpdateBudget tests updating a budget
func TestBudgetsHandler_UpdateBudget(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewBudgetHandler(TestQueries)

	t.Run("Update budget with valid data", func(t *testing.T) {
		// Create a test budget
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)

		newLimit := 6000.0
		reqBody := map[string]interface{}{
			"totalLimit": newLimit,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/api/budgets/"+budgetID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.UpdateBudget(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    BudgetResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, newLimit, response.Data.TotalLimit)

		// Cleanup
		TestQueries.DeleteBudget(ctx, budgetID)
	})
}

// TestBudgetsHandler_DeleteBudget tests deleting a budget
func TestBudgetsHandler_DeleteBudget(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewBudgetHandler(TestQueries)

	t.Run("Delete budget successfully", func(t *testing.T) {
		// Create a test budget
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)

		req := httptest.NewRequest("DELETE", "/api/budgets/"+budgetID, nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.DeleteBudget(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    map[string]string `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
	})
}

// TestBudgetsHandler_AddBudgetCategory tests adding category to budget
func TestBudgetsHandler_AddBudgetCategory(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewBudgetHandler(TestQueries)

	t.Run("Add category to budget", func(t *testing.T) {
		// Create test budget and category
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)
		categoryID := CreateTestCategory(t, ctx, userID, "Food", "🍔", "#FF5733")

		reqBody := map[string]interface{}{
			"categoryId":  categoryID,
			"limitAmount": 1000,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/budgets/"+budgetID+"/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.AddBudgetCategory(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    BudgetCategoryResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, 1000.0, response.Data.LimitAmount)

		// Cleanup
		TestQueries.DeleteBudget(ctx, budgetID)
		TestQueries.DeleteCategory(ctx, categoryID)
	})
}

// TestBudgetsHandler_GetBudgetCategories tests getting budget categories
func TestBudgetsHandler_GetBudgetCategories(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewBudgetHandler(TestQueries)

	t.Run("Get budget categories", func(t *testing.T) {
		// Create test budget and category
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)
		categoryID := CreateTestCategory(t, ctx, userID, "Food", "🍔", "#FF5733")

		// Add category to budget
		TestQueries.AddBudgetCategory(ctx, models.AddBudgetCategoryParams{
			BudgetID:    utils.PgUUID(budgetID),
			CategoryID:  utils.PgUUID(categoryID),
			LimitAmount: utils.PgNumeric(1000),
		})

		req := httptest.NewRequest("GET", "/api/budgets/"+budgetID+"/categories", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.GetBudgetCategories(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    []BudgetCategoryResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotEmpty(t, response.Data)

		// Cleanup
		TestQueries.DeleteBudget(ctx, budgetID)
		TestQueries.DeleteCategory(ctx, categoryID)
	})
}
