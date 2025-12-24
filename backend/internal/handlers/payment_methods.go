package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joselitophala/budget-planner-backend/internal/auth"
	"github.com/joselitophala/budget-planner-backend/internal/models"
	"github.com/joselitophala/budget-planner-backend/internal/utils"
)

// PaymentMethodHandler handles payment method-related requests
type PaymentMethodHandler struct {
	queries *models.Queries
}

// NewPaymentMethodHandler creates a new payment method handler
func NewPaymentMethodHandler(queries *models.Queries) *PaymentMethodHandler {
	return &PaymentMethodHandler{queries: queries}
}

// PaymentMethodResponse represents a payment method in API responses
type PaymentMethodResponse struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	LastFour       *string `json:"lastFour,omitempty"`
	Brand          *string `json:"brand,omitempty"`
	IsDefault      bool    `json:"isDefault"`
	IsActive       bool    `json:"isActive"`
	CreditLimit    *float64 `json:"creditLimit,omitempty"`
	CurrentBalance *float64 `json:"currentBalance,omitempty"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

// CreatePaymentMethodRequest represents the create payment method request
type CreatePaymentMethodRequest struct {
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	LastFour       *string `json:"lastFour,omitempty"`
	Brand          *string `json:"brand,omitempty"`
	IsDefault      bool    `json:"isDefault"`
	CreditLimit    *float64 `json:"creditLimit,omitempty"`
	CurrentBalance *float64 `json:"currentBalance,omitempty"`
}

// UpdatePaymentMethodRequest represents the update payment method request
type UpdatePaymentMethodRequest struct {
	Name           *string `json:"name,omitempty"`
	Type           *string `json:"type,omitempty"`
	LastFour       *string `json:"lastFour,omitempty"`
	Brand          *string `json:"brand,omitempty"`
	IsDefault      *bool   `json:"isDefault,omitempty"`
	IsActive       *bool   `json:"isActive,omitempty"`
	CreditLimit    *float64 `json:"creditLimit,omitempty"`
	CurrentBalance *float64 `json:"currentBalance,omitempty"`
}

// ListPaymentMethods returns all payment methods for the current user
func (h *PaymentMethodHandler) ListPaymentMethods(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	methods, err := h.queries.ListPaymentMethods(r.Context(), utils.PgUUID(userID))
	if err != nil {
		utils.InternalError(w, "Failed to fetch payment methods")
		return
	}

	response := make([]PaymentMethodResponse, len(methods))
	for i, m := range methods {
		response[i] = paymentMethodToResponse(m)
	}

	utils.SendSuccess(w, response)
}

// GetPaymentMethod returns a single payment method by ID
func (h *PaymentMethodHandler) GetPaymentMethod(w http.ResponseWriter, r *http.Request) {
	_, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	methodID := chi.URLParam(r, "id")
	if methodID == "" {
		utils.BadRequest(w, "Payment method ID is required")
		return
	}

	method, err := h.queries.GetPaymentMethodByID(r.Context(), methodID)
	if err != nil {
		utils.NotFound(w, "Payment method not found")
		return
	}

	utils.SendSuccess(w, paymentMethodToResponse(method))
}

// CreatePaymentMethod creates a new payment method
func (h *PaymentMethodHandler) CreatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	var req CreatePaymentMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	method, err := h.queries.CreatePaymentMethod(r.Context(), models.CreatePaymentMethodParams{
		UserID:         utils.PgUUID(userID),
		Name:           req.Name,
		Type:           req.Type,
		LastFour:       utils.PgTextPtr(req.LastFour),
		Brand:          utils.PgTextPtr(req.Brand),
		IsDefault:      utils.PgBool(req.IsDefault),
		IsActive:       utils.PgBool(true),
		CreditLimit:    utils.PgNumericPtr(req.CreditLimit),
		CurrentBalance: utils.PgNumericPtr(req.CurrentBalance),
	})
	if err != nil {
		utils.InternalError(w, "Failed to create payment method")
		return
	}

	utils.SendCreated(w, paymentMethodToResponse(method))
}

// UpdatePaymentMethod updates an existing payment method
func (h *PaymentMethodHandler) UpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	_, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	methodID := chi.URLParam(r, "id")
	if methodID == "" {
		utils.BadRequest(w, "Payment method ID is required")
		return
	}

	var req UpdatePaymentMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body")
		return
	}

	method, err := h.queries.UpdatePaymentMethod(r.Context(), models.UpdatePaymentMethodParams{
		ID:             methodID,
		Name:           utils.PgTextPtr(req.Name),
		Type:           utils.PgTextPtr(req.Type),
		LastFour:       utils.PgTextPtr(req.LastFour),
		Brand:          utils.PgTextPtr(req.Brand),
		IsDefault:      utils.PgBoolPtr(req.IsDefault),
		IsActive:       utils.PgBoolPtr(req.IsActive),
		CreditLimit:    utils.PgNumericPtr(req.CreditLimit),
		CurrentBalance: utils.PgNumericPtr(req.CurrentBalance),
	})
	if err != nil {
		utils.InternalError(w, "Failed to update payment method")
		return
	}

	utils.SendSuccess(w, paymentMethodToResponse(method))
}

// DeletePaymentMethod soft deletes a payment method
func (h *PaymentMethodHandler) DeletePaymentMethod(w http.ResponseWriter, r *http.Request) {
	_, ok := auth.GetUserID(r)
	if !ok {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	methodID := chi.URLParam(r, "id")
	if methodID == "" {
		utils.BadRequest(w, "Payment method ID is required")
		return
	}

	err := h.queries.DeletePaymentMethod(r.Context(), methodID)
	if err != nil {
		utils.InternalError(w, "Failed to delete payment method")
		return
	}

	utils.SendSuccess(w, map[string]string{
		"message": "Payment method deleted successfully",
	})
}

func paymentMethodToResponse(m models.PaymentMethod) PaymentMethodResponse {
	return PaymentMethodResponse{
		ID:             m.ID,
		Name:           m.Name,
		Type:           m.Type,
		LastFour:       utils.TextToStringPtr(m.LastFour),
		Brand:          utils.TextToStringPtr(m.Brand),
		IsDefault:      m.IsDefault.Bool,
		IsActive:       m.IsActive.Bool,
		CreditLimit:    utils.NumericToFloat64Ptr(m.CreditLimit),
		CurrentBalance: utils.NumericToFloat64Ptr(m.CurrentBalance),
		CreatedAt:      utils.TimestamptzToTime(m.CreatedAt).Format(time.RFC3339),
		UpdatedAt:      utils.TimestamptzToTime(m.UpdatedAt).Format(time.RFC3339),
	}
}
