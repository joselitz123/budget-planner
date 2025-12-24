package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPaymentMethodsHandler_ListPaymentMethods tests listing payment methods
func TestPaymentMethodsHandler_ListPaymentMethods(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewPaymentMethodHandler(TestQueries)

	t.Run("List payment methods returns empty list initially", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/payment-methods", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.ListPaymentMethods(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    []PaymentMethodResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Empty(t, response.Data)
	})

	t.Run("List payment methods returns user methods", func(t *testing.T) {
		// Create a test payment method
		methodID := CreateTestPaymentMethod(t, ctx, userID, "Visa", "credit_card")

		req := httptest.NewRequest("GET", "/api/payment-methods", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.ListPaymentMethods(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    []PaymentMethodResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotEmpty(t, response.Data)
		assert.Equal(t, methodID, response.Data[0].ID)

		// Cleanup
		TestQueries.DeletePaymentMethod(ctx, methodID)
	})
}

// TestPaymentMethodsHandler_GetPaymentMethod tests getting a payment method by ID
func TestPaymentMethodsHandler_GetPaymentMethod(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewPaymentMethodHandler(TestQueries)

	t.Run("Get payment method by ID", func(t *testing.T) {
		// Create a test payment method
		methodID := CreateTestPaymentMethod(t, ctx, userID, "Visa", "credit_card")

		req := httptest.NewRequest("GET", "/api/payment-methods/"+methodID, nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.GetPaymentMethod(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    PaymentMethodResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, methodID, response.Data.ID)
		assert.Equal(t, "Visa", response.Data.Name)

		// Cleanup
		TestQueries.DeletePaymentMethod(ctx, methodID)
	})
}

// TestPaymentMethodsHandler_CreatePaymentMethod tests creating a payment method
func TestPaymentMethodsHandler_CreatePaymentMethod(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewPaymentMethodHandler(TestQueries)

	t.Run("Create payment method with valid data", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":     "Chase Sapphire",
			"type":     "credit_card",
			"lastFour": "4242",
			"brand":    "Visa",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/payment-methods", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.CreatePaymentMethod(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    PaymentMethodResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Chase Sapphire", response.Data.Name)
		assert.Equal(t, "4242", *response.Data.LastFour)

		// Cleanup
		TestQueries.DeletePaymentMethod(ctx, response.Data.ID)
	})

	t.Run("Create payment method with invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/payment-methods", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.CreatePaymentMethod(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestPaymentMethodsHandler_UpdatePaymentMethod tests updating a payment method
func TestPaymentMethodsHandler_UpdatePaymentMethod(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewPaymentMethodHandler(TestQueries)

	t.Run("Update payment method with valid data", func(t *testing.T) {
		// Create a test payment method
		methodID := CreateTestPaymentMethod(t, ctx, userID, "Visa", "credit_card")

		newName := "Chase Sapphire Preferred"
		reqBody := map[string]interface{}{
			"name": newName,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/api/payment-methods/"+methodID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.UpdatePaymentMethod(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    PaymentMethodResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, newName, response.Data.Name)

		// Cleanup
		TestQueries.DeletePaymentMethod(ctx, methodID)
	})
}

// TestPaymentMethodsHandler_DeletePaymentMethod tests deleting a payment method
func TestPaymentMethodsHandler_DeletePaymentMethod(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewPaymentMethodHandler(TestQueries)

	t.Run("Delete payment method successfully", func(t *testing.T) {
		// Create a test payment method
		methodID := CreateTestPaymentMethod(t, ctx, userID, "Visa", "credit_card")

		req := httptest.NewRequest("DELETE", "/api/payment-methods/"+methodID, nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.DeletePaymentMethod(w, req)

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
