package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joselitophala/budget-planner-backend/internal/models"
	"github.com/joselitophala/budget-planner-backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTransactionsHandler_ListTransactions tests listing transactions
func TestTransactionsHandler_ListTransactions(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewTransactionHandler(TestQueries)

	t.Run("List transactions returns empty list initially", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/transactions", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.ListTransactions(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    []TransactionResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Empty(t, response.Data)
	})

	t.Run("List transactions with filters", func(t *testing.T) {
		// Create test budget and category
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)
		categoryID := CreateTestCategory(t, ctx, userID, "Food", "🍔", "#FF5733")

		// Create a test transaction
		transactionID := CreateTestTransaction(t, ctx, userID, budgetID, categoryID, 50.00)

		req := httptest.NewRequest("GET", "/api/transactions?limit=10", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.ListTransactions(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    []TransactionResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotEmpty(t, response.Data)

		// Cleanup
		TestQueries.DeleteTransaction(ctx, transactionID)
		TestQueries.DeleteBudget(ctx, budgetID)
		TestQueries.DeleteCategory(ctx, categoryID)
	})
}

// TestTransactionsHandler_GetTransaction tests getting a transaction by ID
func TestTransactionsHandler_GetTransaction(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewTransactionHandler(TestQueries)

	t.Run("Get transaction by ID", func(t *testing.T) {
		// Create test budget and category
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)
		categoryID := CreateTestCategory(t, ctx, userID, "Food", "🍔", "#FF5733")
		transactionID := CreateTestTransaction(t, ctx, userID, budgetID, categoryID, 50.00)

		req := httptest.NewRequest("GET", "/api/transactions/"+transactionID, nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.GetTransaction(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    TransactionResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, transactionID, response.Data.ID)

		// Cleanup
		TestQueries.DeleteTransaction(ctx, transactionID)
		TestQueries.DeleteBudget(ctx, budgetID)
		TestQueries.DeleteCategory(ctx, categoryID)
	})
}

// TestTransactionsHandler_CreateTransaction tests creating a transaction
func TestTransactionsHandler_CreateTransaction(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewTransactionHandler(TestQueries)

	t.Run("Create transaction with valid data", func(t *testing.T) {
		// Create test budget and category
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)
		categoryID := CreateTestCategory(t, ctx, userID, "Food", "🍔", "#FF5733")

		reqBody := map[string]interface{}{
			"budgetId":     budgetID,
			"categoryId":   categoryID,
			"amount":       50.00,
			"type":         "expense",
			"description":  "Lunch",
			"transactionDate": time.Now().Format("2006-01-02"),
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/transactions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.CreateTransaction(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    TransactionResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, 50.00, response.Data.Amount)

		// Cleanup
		TestQueries.DeleteTransaction(ctx, response.Data.ID)
		TestQueries.DeleteBudget(ctx, budgetID)
		TestQueries.DeleteCategory(ctx, categoryID)
	})

	t.Run("Create transaction with invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/transactions", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.CreateTransaction(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestTransactionsHandler_UpdateTransaction tests updating a transaction
func TestTransactionsHandler_UpdateTransaction(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewTransactionHandler(TestQueries)

	t.Run("Update transaction with valid data", func(t *testing.T) {
		// Create test budget, category, and transaction
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)
		categoryID := CreateTestCategory(t, ctx, userID, "Food", "🍔", "#FF5733")
		transactionID := CreateTestTransaction(t, ctx, userID, budgetID, categoryID, 50.00)

		newAmount := 75.00
		reqBody := map[string]interface{}{
			"amount": newAmount,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/api/transactions/"+transactionID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.UpdateTransaction(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    TransactionResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, newAmount, response.Data.Amount)

		// Cleanup
		TestQueries.DeleteTransaction(ctx, transactionID)
		TestQueries.DeleteBudget(ctx, budgetID)
		TestQueries.DeleteCategory(ctx, categoryID)
	})
}

// TestTransactionsHandler_DeleteTransaction tests deleting a transaction
func TestTransactionsHandler_DeleteTransaction(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewTransactionHandler(TestQueries)

	t.Run("Delete transaction successfully", func(t *testing.T) {
		// Create test budget, category, and transaction
		budgetID := CreateTestBudget(t, ctx, userID, "2025-01-01", 5000)
		categoryID := CreateTestCategory(t, ctx, userID, "Food", "🍔", "#FF5733")
		transactionID := CreateTestTransaction(t, ctx, userID, budgetID, categoryID, 50.00)

		req := httptest.NewRequest("DELETE", "/api/transactions/"+transactionID, nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.DeleteTransaction(w, req)

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
		TestQueries.DeleteCategory(ctx, categoryID)
	})
}

// Helper function to create a test transaction
func CreateTestTransaction(t *testing.T, ctx context.Context, userID, budgetID, categoryID string, amount float64) string {
	t.Helper()

	transaction, err := TestQueries.CreateTransaction(ctx, models.CreateTransactionParams{
		UserID:         utils.PgUUID(userID),
		BudgetID:       utils.PgUUIDPtr(&budgetID),
		CategoryID:     utils.PgUUIDPtr(&categoryID),
		Amount:         utils.PgNumeric(amount),
		Type:           utils.PgText("expense"),
		Description:    utils.PgTextPtr(strPtr("Test transaction")),
		TransactionDate: utils.PgDate(time.Now()),
	})
	if err != nil {
		t.Fatalf("Failed to create test transaction: %v", err)
	}
	return transaction.ID
}

func strPtr(s string) *string {
	return &s
}
