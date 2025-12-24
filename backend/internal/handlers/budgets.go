package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joselitophala/budget-planner-backend/internal/auth"
	"github.com/joselitophala/budget-planner-backend/internal/models"
	"github.com/joselitophala/budget-planner-backend/internal/utils"
)

// BudgetHandler handles budget-related requests
type BudgetHandler struct {
	queries *models.Queries
}

// NewBudgetHandler creates a new budget handler
func NewBudgetHandler(queries *models.Queries) *BudgetHandler {
	return &BudgetHandler{queries: queries}
}

// BudgetResponse represents a budget in API responses
type BudgetResponse struct {
	ID         string  `json:"id"`
	UserID     string  `json:"userId"`
	Name       *string `json:"name,omitempty"`
	Month      string  `json:"month"`
	TotalLimit float64 `json:"totalLimit"`
	Spent      float64 `json:"spent"`
	Remaining  float64 `json:"remaining"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

// BudgetCategoryResponse represents a budget category in API responses
type BudgetCategoryResponse struct {
	ID          string  `json:"id"`
	BudgetID    string  `json:"budgetId"`
	CategoryID  string  `json:"categoryId"`
	Name        string  `json:"name"`
	Icon        *string `json:"icon,omitempty"`
	Color       string  `json:"color"`
	LimitAmount float64 `json:"limitAmount"`
	Spent       float64 `json:"spent"`
	Remaining   float64 `json:"remaining"`
}

// CreateBudgetRequest represents the create budget request
type CreateBudgetRequest struct {
	Name       string  `json:"name"`
	Month      string  `json:"month"` // Format: YYYY-MM-DD (first day of month)
	TotalLimit float64 `json:"totalLimit"`
}

// UpdateBudgetRequest represents the update budget request
type UpdateBudgetRequest struct {
	Name       *string  `json:"name,omitempty"`
	TotalLimit *float64 `json:"totalLimit,omitempty"`
}

// AddBudgetCategoryRequest represents the request to add a category to a budget
type AddBudgetCategoryRequest struct {
	CategoryID  string  `json:"categoryId"`
	LimitAmount float64 `json:"limitAmount"`
}

// UpdateBudgetCategoryRequest represents the request to update a budget category
type UpdateBudgetCategoryRequest struct {
	LimitAmount *float64 `json:"limitAmount,omitempty"`
}

// ListBudgets returns all budgets for the current user
func (h *BudgetHandler) ListBudgets(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	budgets, err := h.queries.ListUserBudgets(r.Context(), utils.PgUUID(userID))
	if err != nil {
		utils.InternalError(w, "Failed to fetch budgets")
		return
	}

	response := make([]BudgetResponse, len(budgets))
	for i, budget := range budgets {
		spent, _ := h.getBudgetSpent(r.Context(), budget.ID)
		totalLimit := utils.NumericToFloat64(budget.TotalLimit)
		name := utils.TextToStringPtr(budget.Name)
		response[i] = BudgetResponse{
			ID:         budget.ID,
			UserID:     userID,
			Name:       name,
			Month:      utils.DateToTime(budget.Month).Format("2006-01-02"),
			TotalLimit: totalLimit,
			Spent:      spent,
			Remaining:  totalLimit - spent,
			CreatedAt:  utils.TimestamptzToTime(budget.CreatedAt).Format(time.RFC3339),
			UpdatedAt:  utils.TimestamptzToTime(budget.UpdatedAt).Format(time.RFC3339),
		}
	}

	utils.SendSuccess(w, response)
}

// GetBudgetByMonth returns a budget for a specific month
func (h *BudgetHandler) GetBudgetByMonth(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	monthStr := chi.URLParam(r, "month")
	if monthStr == "" {
		utils.BadRequest(w, "Month is required (format: YYYY-MM)")
		return
	}

	// Parse month string to date (first day of month)
	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		utils.BadRequest(w, "Invalid month format. Use YYYY-MM")
		return
	}

	budget, err := h.queries.GetBudgetByMonth(r.Context(), models.GetBudgetByMonthParams{
		UserID: utils.PgUUID(userID),
		Month:  utils.PgDate(month),
	})
	if err != nil {
		utils.NotFound(w, "Budget not found for this month")
		return
	}

	spent, _ := h.getBudgetSpent(r.Context(), budget.ID)
	totalLimit := utils.NumericToFloat64(budget.TotalLimit)
	name := utils.TextToStringPtr(budget.Name)

	utils.SendSuccess(w, BudgetResponse{
		ID:         budget.ID,
		UserID:     userID,
		Name:       name,
		Month:      utils.DateToTime(budget.Month).Format("2006-01-02"),
		TotalLimit: totalLimit,
		Spent:      spent,
		Remaining:  totalLimit - spent,
		CreatedAt:  utils.TimestamptzToTime(budget.CreatedAt).Format(time.RFC3339),
		UpdatedAt:  utils.TimestamptzToTime(budget.UpdatedAt).Format(time.RFC3339),
	})
}

// GetBudget returns a budget by ID
func (h *BudgetHandler) GetBudget(w http.ResponseWriter, r *http.Request) {
	budgetID := chi.URLParam(r, "id")
	if budgetID == "" {
		utils.BadRequest(w, "Budget ID is required")
		return
	}

	budget, err := h.queries.GetBudgetByID(r.Context(), budgetID)
	if err != nil {
		utils.NotFound(w, "Budget not found")
		return
	}

	spent, _ := h.getBudgetSpent(r.Context(), budget.ID)
	totalLimit := utils.NumericToFloat64(budget.TotalLimit)
	name := utils.TextToStringPtr(budget.Name)
	userID := utils.UUIDToString(budget.UserID)

	utils.SendSuccess(w, BudgetResponse{
		ID:         budget.ID,
		UserID:     userID,
		Name:       name,
		Month:      utils.DateToTime(budget.Month).Format("2006-01-02"),
		TotalLimit: totalLimit,
		Spent:      spent,
		Remaining:  totalLimit - spent,
		CreatedAt:  utils.TimestamptzToTime(budget.CreatedAt).Format(time.RFC3339),
		UpdatedAt:  utils.TimestamptzToTime(budget.UpdatedAt).Format(time.RFC3339),
	})
}

// CreateBudget creates a new budget
func (h *BudgetHandler) CreateBudget(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	var req CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	// Parse month
	month, err := time.Parse("2006-01-02", req.Month)
	if err != nil {
		utils.BadRequest(w, "Invalid month format. Use YYYY-MM-DD")
		return
	}

	budget, err := h.queries.CreateBudget(r.Context(), models.CreateBudgetParams{
		UserID:     utils.PgUUID(userID),
		Name:       utils.PgText(req.Name),
		Month:      utils.PgDate(month),
		TotalLimit: utils.PgNumeric(req.TotalLimit),
	})
	if err != nil {
		utils.InternalError(w, "Failed to create budget")
		return
	}

	totalLimit := utils.NumericToFloat64(budget.TotalLimit)
	name := utils.TextToStringPtr(budget.Name)

	utils.SendCreated(w, BudgetResponse{
		ID:         budget.ID,
		UserID:     userID,
		Name:       name,
		Month:      utils.DateToTime(budget.Month).Format("2006-01-02"),
		TotalLimit: totalLimit,
		Spent:      0,
		Remaining:  totalLimit,
		CreatedAt:  utils.TimestamptzToTime(budget.CreatedAt).Format(time.RFC3339),
		UpdatedAt:  utils.TimestamptzToTime(budget.UpdatedAt).Format(time.RFC3339),
	})
}

// UpdateBudget updates an existing budget
func (h *BudgetHandler) UpdateBudget(w http.ResponseWriter, r *http.Request) {
	budgetID := chi.URLParam(r, "id")
	if budgetID == "" {
		utils.BadRequest(w, "Budget ID is required")
		return
	}

	var req UpdateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	budget, err := h.queries.UpdateBudget(r.Context(), models.UpdateBudgetParams{
		ID:         budgetID,
		Name:       utils.PgTextPtr(req.Name),
		TotalLimit: utils.PgNumericPtr(req.TotalLimit),
	})
	if err != nil {
		utils.InternalError(w, "Failed to update budget")
		return
	}

	spent, _ := h.getBudgetSpent(r.Context(), budget.ID)
	totalLimit := utils.NumericToFloat64(budget.TotalLimit)
	name := utils.TextToStringPtr(budget.Name)
	userID := utils.UUIDToString(budget.UserID)

	utils.SendSuccess(w, BudgetResponse{
		ID:         budget.ID,
		UserID:     userID,
		Name:       name,
		Month:      utils.DateToTime(budget.Month).Format("2006-01-02"),
		TotalLimit: totalLimit,
		Spent:      spent,
		Remaining:  totalLimit - spent,
		CreatedAt:  utils.TimestamptzToTime(budget.CreatedAt).Format(time.RFC3339),
		UpdatedAt:  utils.TimestamptzToTime(budget.UpdatedAt).Format(time.RFC3339),
	})
}

// DeleteBudget soft deletes a budget
func (h *BudgetHandler) DeleteBudget(w http.ResponseWriter, r *http.Request) {
	budgetID := chi.URLParam(r, "id")
	if budgetID == "" {
		utils.BadRequest(w, "Budget ID is required")
		return
	}

	err := h.queries.DeleteBudget(r.Context(), budgetID)
	if err != nil {
		utils.InternalError(w, "Failed to delete budget")
		return
	}

	utils.SendSuccess(w, map[string]string{
		"message": "Budget deleted successfully",
	})
}

// GetBudgetCategories returns all categories for a budget
func (h *BudgetHandler) GetBudgetCategories(w http.ResponseWriter, r *http.Request) {
	budgetID := chi.URLParam(r, "id")
	if budgetID == "" {
		utils.BadRequest(w, "Budget ID is required")
		return
	}

	categories, err := h.queries.GetBudgetCategories(r.Context(), utils.PgUUID(budgetID))
	if err != nil {
		utils.InternalError(w, "Failed to fetch budget categories")
		return
	}

	response := make([]BudgetCategoryResponse, len(categories))
	for i, cat := range categories {
		spent, _ := h.getCategorySpent(r.Context(), budgetID, utils.UUIDToString(cat.CategoryID))
		limitAmount := utils.NumericToFloat64(cat.LimitAmount)
		response[i] = BudgetCategoryResponse{
			ID:          cat.ID,
			BudgetID:    utils.UUIDToString(cat.BudgetID),
			CategoryID:  utils.UUIDToString(cat.CategoryID),
			Name:        cat.Name,
			Icon:        utils.TextToStringPtr(cat.Icon),
			Color:       utils.TextToString(cat.Color),
			LimitAmount: limitAmount,
			Spent:       spent,
			Remaining:   limitAmount - spent,
		}
	}

	utils.SendSuccess(w, response)
}

// AddBudgetCategory adds a category to a budget
func (h *BudgetHandler) AddBudgetCategory(w http.ResponseWriter, r *http.Request) {
	budgetID := chi.URLParam(r, "id")
	if budgetID == "" {
		utils.BadRequest(w, "Budget ID is required")
		return
	}

	var req AddBudgetCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	bc, err := h.queries.AddBudgetCategory(r.Context(), models.AddBudgetCategoryParams{
		BudgetID:    utils.PgUUID(budgetID),
		CategoryID:  utils.PgUUID(req.CategoryID),
		LimitAmount: utils.PgNumeric(req.LimitAmount),
	})
	if err != nil {
		utils.InternalError(w, "Failed to add category to budget")
		return
	}

	limitAmount := utils.NumericToFloat64(bc.LimitAmount)

	utils.SendCreated(w, BudgetCategoryResponse{
		ID:          bc.ID,
		BudgetID:    utils.UUIDToString(bc.BudgetID),
		CategoryID:  utils.UUIDToString(bc.CategoryID),
		LimitAmount: limitAmount,
		Spent:       0,
		Remaining:   limitAmount,
	})
}

// UpdateBudgetCategory updates a budget category limit
func (h *BudgetHandler) UpdateBudgetCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "categoryId")
	if categoryID == "" {
		utils.BadRequest(w, "Category ID is required")
		return
	}

	var req UpdateBudgetCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	bc, err := h.queries.UpdateBudgetCategory(r.Context(), models.UpdateBudgetCategoryParams{
		ID:          categoryID,
		LimitAmount: utils.PgNumericPtr(req.LimitAmount),
	})
	if err != nil {
		utils.InternalError(w, "Failed to update budget category")
		return
	}

	budgetID := utils.UUIDToString(bc.BudgetID)
	catID := utils.UUIDToString(bc.CategoryID)
	spent, _ := h.getCategorySpent(r.Context(), budgetID, catID)
	limitAmount := utils.NumericToFloat64(bc.LimitAmount)

	utils.SendSuccess(w, BudgetCategoryResponse{
		ID:          bc.ID,
		BudgetID:    budgetID,
		CategoryID:  catID,
		LimitAmount: limitAmount,
		Spent:       spent,
		Remaining:   limitAmount - spent,
	})
}

// RemoveBudgetCategory removes a category from a budget
func (h *BudgetHandler) RemoveBudgetCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "categoryId")
	if categoryID == "" {
		utils.BadRequest(w, "Category ID is required")
		return
	}

	err := h.queries.RemoveBudgetCategory(r.Context(), categoryID)
	if err != nil {
		utils.InternalError(w, "Failed to remove category from budget")
		return
	}

	utils.SendSuccess(w, map[string]string{
		"message": "Category removed from budget successfully",
	})
}

// Helper function to get budget spent amount
func (h *BudgetHandler) getBudgetSpent(ctx context.Context, budgetID string) (float64, error) {
	result, err := h.queries.GetBudgetSpent(ctx, utils.PgUUID(budgetID))
	if err != nil {
		return 0, err
	}

	// GetBudgetSpent returns interface{}, need to handle it
	if result == nil {
		return 0, nil
	}

	// Try to convert to float64
	switch v := result.(type) {
	case float64:
		return v, nil
	case []byte:
		// PostgreSQL may return numeric as []byte
		var f float64
		_, err := fmt.Sscanf(string(v), "%f", &f)
		return f, err
	default:
		return 0, nil
	}
}

// Helper function to get category spent amount for a budget
func (h *BudgetHandler) getCategorySpent(ctx context.Context, budgetID, categoryID string) (float64, error) {
	result, err := h.queries.GetCategorySpent(ctx, models.GetCategorySpentParams{
		BudgetID:   utils.PgUUID(budgetID),
		CategoryID: utils.PgUUID(categoryID),
	})
	if err != nil {
		return 0, err
	}

	// GetCategorySpent returns interface{}, need to handle it
	if result == nil {
		return 0, nil
	}

	// Try to convert to float64
	switch v := result.(type) {
	case float64:
		return v, nil
	case []byte:
		// PostgreSQL may return numeric as []byte
		var f float64
		_, err := fmt.Sscanf(string(v), "%f", &f)
		return f, err
	default:
		return 0, nil
	}
}
