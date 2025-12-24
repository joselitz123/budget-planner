package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joselitophala/budget-planner-backend/internal/auth"
	"github.com/joselitophala/budget-planner-backend/internal/models"
	"github.com/joselitophala/budget-planner-backend/internal/utils"
)

// AnalyticsHandler handles analytics and reporting requests
type AnalyticsHandler struct {
	queries *models.Queries
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(queries *models.Queries) *AnalyticsHandler {
	return &AnalyticsHandler{queries: queries}
}

// DashboardSummary represents the dashboard summary for a month
type DashboardSummary struct {
	Month              string  `json:"month"`
	TotalBudget        float64 `json:"totalBudget"`
	TotalSpent         float64 `json:"totalSpent"`
	Remaining          float64 `json:"remaining"`
	BudgetUsedPercent  float64 `json:"budgetUsedPercent"`
	TransactionCount   int32   `json:"transactionCount"`
	TopCategories      []CategorySpending `json:"topCategories"`
	RecentTransactions []TransactionSummary `json:"recentTransactions"`
}

// CategorySpending represents spending by category
type CategorySpending struct {
	CategoryID   string  `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	Amount       float64 `json:"amount"`
	Percent      float64 `json:"percent"`
}

// TransactionSummary represents a transaction summary
type TransactionSummary struct {
	ID          string  `json:"id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	Category    string  `json:"category,omitempty"`
}

// GetDashboard returns the dashboard summary for a specific month
func (h *AnalyticsHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	monthStr := chi.URLParam(r, "month")
	if monthStr == "" {
		monthStr = time.Now().Format("2006-01")
	}

	// Parse month
	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		utils.BadRequest(w, "Invalid month format. Use YYYY-MM")
		return
	}

	// Get budget for this month
	budget, err := h.queries.GetBudgetByMonth(r.Context(), models.GetBudgetByMonthParams{
		UserID: utils.PgUUID(userID),
		Month:  utils.PgDate(month),
	})
	if err != nil {
		utils.NotFound(w, "Budget not found for this month")
		return
	}

	// Get dashboard summary (using budget ID)
	spending, err := h.queries.GetDashboardSummary(r.Context(), budget.ID)
	if err != nil {
		utils.InternalError(w, "Failed to fetch dashboard data")
		return
	}

	// Get category breakdown (using budget ID)
	categorySpending, err := h.queries.GetSpendingByCategory(r.Context(), utils.PgUUID(budget.ID))
	if err != nil {
		utils.InternalError(w, "Failed to fetch category spending")
		return
	}

	// Get recent transactions
	recent, err := h.queries.GetRecentTransactions(r.Context(), models.GetRecentTransactionsParams{
		UserID: utils.PgUUID(userID),
		Limit:  5,
		Offset: 0,
	})
	if err != nil {
		utils.InternalError(w, "Failed to fetch recent transactions")
		return
	}

	// Build response
	totalLimit := utils.NumericToFloat64(budget.TotalLimit)
	totalSpent := toFloat64(spending.TotalSpent)
	transactionCount := int(spending.TransactionCount)

	summary := DashboardSummary{
		Month:             monthStr,
		TotalBudget:       totalLimit,
		TotalSpent:        totalSpent,
		Remaining:         totalLimit - totalSpent,
		BudgetUsedPercent: 0,
		TransactionCount:  int32(transactionCount),
	}
	if totalLimit > 0 {
		summary.BudgetUsedPercent = (totalSpent / totalLimit) * 100
	}

	// Build category spending
	summary.TopCategories = make([]CategorySpending, len(categorySpending))
	for i, cat := range categorySpending {
		amount := toFloat64(cat.TotalSpent)
		percent := 0.0
		if totalSpent > 0 {
			percent = (amount / totalSpent) * 100
		}
		summary.TopCategories[i] = CategorySpending{
			CategoryID:   cat.ID,
			CategoryName: cat.Name,
			Amount:       amount,
			Percent:      percent,
		}
	}

	// Build recent transactions
	summary.RecentTransactions = make([]TransactionSummary, len(recent))
	for i, t := range recent {
		summary.RecentTransactions[i] = TransactionSummary{
			ID:          t.ID,
			Amount:      utils.NumericToFloat64(t.Amount),
			Description: utils.TextToString(t.Description),
			Date:        utils.DateToTime(t.TransactionDate).Format("2006-01-02"),
			Category:    utils.TextToString(t.CategoryName),
		}
	}

	utils.SendSuccess(w, summary)
}

// GetSpendingReport returns a detailed spending report for a month
func (h *AnalyticsHandler) GetSpendingReport(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	monthStr := chi.URLParam(r, "month")
	if monthStr == "" {
		monthStr = time.Now().Format("2006-01")
	}

	// Parse month
	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		utils.BadRequest(w, "Invalid month format. Use YYYY-MM")
		return
	}

	// Get budget for this month
	budget, err := h.queries.GetBudgetByMonth(r.Context(), models.GetBudgetByMonthParams{
		UserID: utils.PgUUID(userID),
		Month:  utils.PgDate(month),
	})
	if err != nil {
		utils.NotFound(w, "Budget not found for this month")
		return
	}

	// Get spending by category
	categorySpending, err := h.queries.GetSpendingByCategory(r.Context(), utils.PgUUID(budget.ID))
	if err != nil {
		utils.InternalError(w, "Failed to fetch spending report")
		return
	}

	utils.SendSuccess(w, categorySpending)
}

// GetTrends returns historical spending trends
func (h *AnalyticsHandler) GetTrends(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	// Default to last 6 months
	months := 6
	if monthsStr := r.URL.Query().Get("months"); monthsStr != "" {
		if m, err := parseMonthInt(monthsStr); err == nil && m > 0 && m <= 24 {
			months = m
		}
	}

	// Calculate date range (last N months)
	endDate := time.Now().AddDate(0, 1, 0).Truncate(time.Hour * 24)
	startDate := endDate.AddDate(0, -months, 0)

	trends, err := h.queries.GetSpendingTrends(r.Context(), models.GetSpendingTrendsParams{
		UserID:          utils.PgUUID(userID),
		TransactionDate: utils.PgDate(startDate),
		TransactionDate_2: utils.PgDate(endDate),
	})
	if err != nil {
		utils.InternalError(w, "Failed to fetch trends")
		return
	}

	utils.SendSuccess(w, trends)
}

// GetCategoryReport returns a detailed report for a specific category
func (h *AnalyticsHandler) GetCategoryReport(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	categoryID := chi.URLParam(r, "categoryId")
	if categoryID == "" {
		utils.BadRequest(w, "Category ID is required")
		return
	}

	// Get date range from query params
	startDate := time.Now().AddDate(0, -1, 0) // Default: last month
	endDate := time.Now()

	if startStr := r.URL.Query().Get("startDate"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = t
		}
	}
	if endStr := r.URL.Query().Get("endDate"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = t
		}
	}

	report, err := h.queries.GetCategoryReport(r.Context(), models.GetCategoryReportParams{
		UserID:          utils.PgUUID(userID),
		CategoryID:      utils.PgUUID(categoryID),
		TransactionDate: utils.PgDate(startDate),
		TransactionDate_2: utils.PgDate(endDate),
	})
	if err != nil {
		utils.InternalError(w, "Failed to fetch category report")
		return
	}

	utils.SendSuccess(w, report)
}

// Helper functions
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case int64:
		return float64(val)
	case float64:
		return val
	case int32:
		return float64(val)
	default:
		return 0
	}
}

func parseMonthInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
