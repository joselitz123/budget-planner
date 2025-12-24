package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joselitophala/budget-planner-backend/internal/auth"
	"github.com/joselitophala/budget-planner-backend/internal/models"
	"github.com/joselitophala/budget-planner-backend/internal/utils"
)

// TransactionHandler handles transaction-related requests
type TransactionHandler struct {
	queries *models.Queries
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(queries *models.Queries) *TransactionHandler {
	return &TransactionHandler{queries: queries}
}

// TransactionResponse represents a transaction in API responses
type TransactionResponse struct {
	ID                  string      `json:"id"`
	BudgetID            *string     `json:"budgetId,omitempty"`
	CategoryID          *string     `json:"categoryId,omitempty"`
	PaymentMethodID     *string     `json:"paymentMethodId,omitempty"`
	Amount              float64     `json:"amount"`
	Type                string      `json:"type"`
	IsTransfer          bool        `json:"isTransfer"`
	TransferToAccountID *string     `json:"transferToAccountId,omitempty"`
	Description         *string     `json:"description,omitempty"`
	TransactionDate     string      `json:"transactionDate"`
	IsRecurring         bool        `json:"isRecurring"`
	RecurrencePattern   interface{} `json:"recurrencePattern,omitempty"`
	CreatedAt           string      `json:"createdAt"`
	UpdatedAt           string      `json:"updatedAt"`
}

// CreateTransactionRequest represents the create transaction request
type CreateTransactionRequest struct {
	BudgetID          *string     `json:"budgetId,omitempty"`
	CategoryID        *string     `json:"categoryId,omitempty"`
	PaymentMethodID   *string     `json:"paymentMethodId,omitempty"`
	Amount            float64     `json:"amount"`
	Type              string      `json:"type"`
	IsTransfer        bool        `json:"isTransfer"`
	TransferToAccountID *string   `json:"transferToAccountId,omitempty"`
	Description       *string     `json:"description,omitempty"`
	TransactionDate   string      `json:"transactionDate"`
	IsRecurring       bool        `json:"isRecurring"`
	RecurrencePattern interface{} `json:"recurrencePattern,omitempty"`
}

// UpdateTransactionRequest represents the update transaction request
type UpdateTransactionRequest struct {
	BudgetID          *string     `json:"budgetId,omitempty"`
	CategoryID        *string     `json:"categoryId,omitempty"`
	PaymentMethodID   *string     `json:"paymentMethodId,omitempty"`
	Amount            *float64    `json:"amount,omitempty"`
	Type              *string     `json:"type,omitempty"`
	IsTransfer        *bool       `json:"isTransfer,omitempty"`
	TransferToAccountID *string   `json:"transferToAccountId,omitempty"`
	Description       *string     `json:"description,omitempty"`
	TransactionDate   *string     `json:"transactionDate,omitempty"`
	IsRecurring       *bool       `json:"isRecurring,omitempty"`
	RecurrencePattern interface{} `json:"recurrencePattern,omitempty"`
}

// ListTransactions returns transactions with optional filters
func (h *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	// Parse query parameters
	var startDate, endDate *time.Time
	var categoryID, budgetID *string

	if startDateStr := r.URL.Query().Get("startDate"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr := r.URL.Query().Get("endDate"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &t
		}
	}
	if catID := r.URL.Query().Get("category"); catID != "" {
		categoryID = &catID
	}
	if budID := r.URL.Query().Get("budget"); budID != "" {
		budgetID = &budID
	}

	// Default pagination
	limit := 50
	offset := 0
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := parseInt(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := parseInt(offsetStr); err == nil {
			offset = o
		}
	}

	transactions, err := h.queries.ListTransactions(r.Context(), models.ListTransactionsParams{
		UserID:  utils.PgUUID(userID),
		Column2: utils.PgDatePtr(startDate),
		Column3: utils.PgDatePtr(endDate),
		Column4: stringPtrOrEmpty(categoryID),
		Column5: stringPtrOrEmpty(budgetID),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		utils.InternalError(w, "Failed to fetch transactions")
		return
	}

	response := make([]TransactionResponse, len(transactions))
	for i, t := range transactions {
		response[i] = transactionToResponse(t)
	}

	utils.SendSuccess(w, response)
}

// GetTransaction returns a single transaction by ID
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := chi.URLParam(r, "id")
	if transactionID == "" {
		utils.BadRequest(w, "Transaction ID is required")
		return
	}

	transaction, err := h.queries.GetTransactionByID(r.Context(), transactionID)
	if err != nil {
		utils.NotFound(w, "Transaction not found")
		return
	}

	utils.SendSuccess(w, transactionToResponse(transaction))
}

// CreateTransaction creates a new transaction
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	// Parse transaction date
	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		utils.BadRequest(w, "Invalid transaction date format. Use YYYY-MM-DD")
		return
	}

	// Handle recurrence pattern JSON
	var recurrencePattern []byte
	if req.RecurrencePattern != nil {
		if data, err := json.Marshal(req.RecurrencePattern); err == nil {
			recurrencePattern = data
		}
	}

	transaction, err := h.queries.CreateTransaction(r.Context(), models.CreateTransactionParams{
		UserID:              utils.PgUUID(userID),
		BudgetID:            utils.PgUUIDPtr(req.BudgetID),
		CategoryID:          utils.PgUUIDPtr(req.CategoryID),
		PaymentMethodID:     utils.PgUUIDPtr(req.PaymentMethodID),
		Amount:              utils.PgNumeric(req.Amount),
		Type:                utils.PgText(req.Type),
		IsTransfer:          pgBool(req.IsTransfer),
		TransferToAccountID: utils.PgUUIDPtr(req.TransferToAccountID),
		Description:         utils.PgTextPtr(req.Description),
		TransactionDate:     utils.PgDate(transactionDate),
		IsRecurring:         pgBool(req.IsRecurring),
		RecurrencePattern:   recurrencePattern,
	})
	if err != nil {
		utils.InternalError(w, "Failed to create transaction")
		return
	}

	utils.SendCreated(w, transactionToResponse(transaction))
}

// UpdateTransaction updates an existing transaction
func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := chi.URLParam(r, "id")
	if transactionID == "" {
		utils.BadRequest(w, "Transaction ID is required")
		return
	}

	var req UpdateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	// Parse transaction date if provided
	var transactionDate *time.Time
	if req.TransactionDate != nil {
		if t, err := time.Parse("2006-01-02", *req.TransactionDate); err == nil {
			transactionDate = &t
		}
	}

	// Handle recurrence pattern JSON
	var recurrencePattern []byte
	if req.RecurrencePattern != nil {
		if data, err := json.Marshal(req.RecurrencePattern); err == nil {
			recurrencePattern = data
		}
	}

	transaction, err := h.queries.UpdateTransaction(r.Context(), models.UpdateTransactionParams{
		ID:                  transactionID,
		BudgetID:            utils.PgUUIDPtr(req.BudgetID),
		CategoryID:          utils.PgUUIDPtr(req.CategoryID),
		PaymentMethodID:     utils.PgUUIDPtr(req.PaymentMethodID),
		Amount:              utils.PgNumericPtr(req.Amount),
		Type:                utils.PgTextPtr(req.Type),
		IsTransfer:          pgBoolPtr(req.IsTransfer),
		TransferToAccountID: utils.PgUUIDPtr(req.TransferToAccountID),
		Description:         utils.PgTextPtr(req.Description),
		TransactionDate:     utils.PgDatePtr(transactionDate),
		IsRecurring:         pgBoolPtr(req.IsRecurring),
		RecurrencePattern:   recurrencePattern,
	})
	if err != nil {
		utils.InternalError(w, "Failed to update transaction")
		return
	}

	utils.SendSuccess(w, transactionToResponse(transaction))
}

// DeleteTransaction soft deletes a transaction
func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := chi.URLParam(r, "id")
	if transactionID == "" {
		utils.BadRequest(w, "Transaction ID is required")
		return
	}

	err := h.queries.DeleteTransaction(r.Context(), transactionID)
	if err != nil {
		utils.InternalError(w, "Failed to delete transaction")
		return
	}

	utils.SendSuccess(w, map[string]string{
		"message": "Transaction deleted successfully",
	})
}

// GetBudgetTransactions returns all transactions for a specific budget
func (h *TransactionHandler) GetBudgetTransactions(w http.ResponseWriter, r *http.Request) {
	budgetID := chi.URLParam(r, "budgetId")
	if budgetID == "" {
		utils.BadRequest(w, "Budget ID is required")
		return
	}

	transactions, err := h.queries.GetTransactionsByBudget(r.Context(), utils.PgUUID(budgetID))
	if err != nil {
		utils.InternalError(w, "Failed to fetch transactions")
		return
	}

	response := make([]TransactionResponse, len(transactions))
	for i, t := range transactions {
		response[i] = transactionToResponse(t)
	}

	utils.SendSuccess(w, response)
}

// Helper function to convert transaction model to response
func transactionToResponse(t models.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:                  t.ID,
		BudgetID:            uuidPtrToString(t.BudgetID),
		CategoryID:          uuidPtrToString(t.CategoryID),
		PaymentMethodID:     uuidPtrToString(t.PaymentMethodID),
		Amount:              utils.NumericToFloat64(t.Amount),
		Type:                utils.TextToString(t.Type),
		IsTransfer:          t.IsTransfer.Bool,
		TransferToAccountID: uuidPtrToString(t.TransferToAccountID),
		Description:         utils.TextToStringPtr(t.Description),
		TransactionDate:     utils.DateToTime(t.TransactionDate).Format("2006-01-02"),
		IsRecurring:         t.IsRecurring.Bool,
		RecurrencePattern:   t.RecurrencePattern,
		CreatedAt:           utils.TimestamptzToTime(t.CreatedAt).Format(time.RFC3339),
		UpdatedAt:           utils.TimestamptzToTime(t.UpdatedAt).Format(time.RFC3339),
	}
}

// Helper function to parse int
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

// Helper function to convert pgtype.UUID pointer to string pointer
func uuidPtrToString(ptr pgtype.UUID) *string {
	if ptr.Valid {
		s := utils.UUIDToString(ptr)
		return &s
	}
	return nil
}

// Helper function to convert *bool to pgtype.Bool (for update operations)
func pgBoolPtr(b *bool) pgtype.Bool {
	if b != nil {
		return pgtype.Bool{Valid: true, Bool: *b}
	}
	return pgtype.Bool{}
}

// Helper function to convert bool to pgtype.Bool
func pgBool(b bool) pgtype.Bool {
	return pgtype.Bool{Valid: true, Bool: b}
}

// Helper function to convert *string to string (empty if nil)
func stringPtrOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
