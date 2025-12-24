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

// TestCategoriesHandler_ListCategories tests listing categories
func TestCategoriesHandler_ListCategories(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewCategoryHandler(TestQueries)

	t.Run("List categories returns empty list initially", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/categories", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.ListCategories(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    []CategoryResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Empty(t, response.Data)
	})

	t.Run("List categories returns user categories", func(t *testing.T) {
		// Create a test category
		categoryID := CreateTestCategory(t, ctx, userID, "Food", "🍔", "#FF5733")

		req := httptest.NewRequest("GET", "/api/categories", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.ListCategories(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    []CategoryResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotEmpty(t, response.Data)
		assert.Equal(t, categoryID, response.Data[0].ID)
		assert.Equal(t, "Food", response.Data[0].Name)

		// Cleanup
		TestQueries.DeleteCategory(ctx, categoryID)
	})
}

// TestCategoriesHandler_GetSystemCategories tests getting system categories
func TestCategoriesHandler_GetSystemCategories(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewCategoryHandler(TestQueries)

	t.Run("Get system categories", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/categories/system", nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.GetSystemCategories(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    []CategoryResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		// System categories should be empty initially unless seeded
		assert.NotNil(t, response.Data)
	})
}

// TestCategoriesHandler_CreateCategory tests creating a category
func TestCategoriesHandler_CreateCategory(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewCategoryHandler(TestQueries)

	t.Run("Create category with valid data", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":  "Transport",
			"icon":  "🚗",
			"color": "#3498db",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.CreateCategory(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    CategoryResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Transport", response.Data.Name)
		assert.Equal(t, "🚗", *response.Data.Icon)
		assert.Equal(t, "#3498db", response.Data.Color)

		// Cleanup
		TestQueries.DeleteCategory(ctx, response.Data.ID)
	})

	t.Run("Create category with invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/categories", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.CreateCategory(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestCategoriesHandler_UpdateCategory tests updating a category
func TestCategoriesHandler_UpdateCategory(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewCategoryHandler(TestQueries)

	t.Run("Update category with valid data", func(t *testing.T) {
		// Create a test category
		categoryID := CreateTestCategory(t, ctx, userID, "Food", "🍔", "#FF5733")

		reqBody := map[string]interface{}{
			"name":  "Food & Drinks",
			"color": "#FF5734",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/api/categories/"+categoryID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.UpdateCategory(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool `json:"success"`
			Data    CategoryResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Food & Drinks", response.Data.Name)

		// Cleanup
		TestQueries.DeleteCategory(ctx, categoryID)
	})
}

// TestCategoriesHandler_DeleteCategory tests deleting a category
func TestCategoriesHandler_DeleteCategory(t *testing.T) {
	ctx := context.Background()
	userID := CreateTestUser(t, ctx)
	defer CleanupTestUser(t, ctx, userID)

	h := NewCategoryHandler(TestQueries)

	t.Run("Delete category successfully", func(t *testing.T) {
		// Create a test category
		categoryID := CreateTestCategory(t, ctx, userID, "Food", "🍔", "#FF5733")

		req := httptest.NewRequest("DELETE", "/api/categories/"+categoryID, nil)
		setAuthContext(req, userID)
		w := httptest.NewRecorder()

		h.DeleteCategory(w, req)

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
