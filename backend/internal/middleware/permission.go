package middleware

import (
	"context"
	"net/http"

	"github.com/joselitophala/budget-planner-backend/internal/auth"
	"github.com/joselitophala/budget-planner-backend/internal/models"
	"github.com/joselitophala/budget-planner-backend/internal/utils"
)

// PermissionLevel represents the access level for a shared resource
type PermissionLevel string

const (
	// PermissionOwner means the user owns the resource
	PermissionOwner PermissionLevel = "owner"
	// PermissionEdit means the user can edit the resource
	PermissionEdit PermissionLevel = "edit"
	// PermissionView means the user can only view the resource
	PermissionView PermissionLevel = "view"
)

// BudgetPermission checks if the user has access to a budget
type BudgetPermission struct {
	queries *models.Queries
}

// NewBudgetPermission creates a new budget permission checker
func NewBudgetPermission(queries *models.Queries) *BudgetPermission {
	return &BudgetPermission{queries: queries}
}

// RequireOwner checks if the user is the owner of the budget
func (bp *BudgetPermission) RequireOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		budgetID := r.PathValue("id")
		if budgetID == "" {
			budgetID = r.PathValue("budgetId")
		}
		if budgetID == "" {
			utils.BadRequest(w, "Budget ID is required")
			return
		}

		userID, ok := auth.GetUserID(r)
		if !ok {
			utils.Unauthorized(w, "User not authenticated")
			return
		}

		budget, err := bp.queries.GetBudgetByID(r.Context(), budgetID)
		if err != nil {
			utils.NotFound(w, "Budget not found")
			return
		}

		if budget.UserID != utils.PgUUID(userID) {
			utils.Forbidden(w, "You don't have permission to access this budget")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireAccess checks if the user has access to a budget (owner or shared)
func (bp *BudgetPermission) RequireAccess(minPermission PermissionLevel) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			budgetID := r.PathValue("id")
			if budgetID == "" {
				budgetID = r.PathValue("budgetId")
			}
			if budgetID == "" {
				// Try to get budget ID from query param
				budgetID = r.URL.Query().Get("budgetId")
			}
			if budgetID == "" {
				utils.BadRequest(w, "Budget ID is required")
				return
			}

			userID, ok := auth.GetUserID(r)
			if !ok {
				utils.Unauthorized(w, "User not authenticated")
				return
			}

			// Check permission
			permission, err := bp.checkBudgetAccess(r.Context(), budgetID, userID)
			if err != nil || permission == "" {
				utils.NotFound(w, "Budget not found or no access")
				return
			}

			// Check if user has required permission level
			if !hasSufficientPermission(permission, minPermission) {
				utils.Forbidden(w, "Insufficient permissions")
				return
			}

			// Add permission to context
			ctx := context.WithValue(r.Context(), "budgetPermission", permission)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// checkBudgetAccess checks what level of access a user has to a budget
func (bp *BudgetPermission) checkBudgetAccess(ctx context.Context, budgetID, userID string) (PermissionLevel, error) {
	// First check if user is the owner
	budget, err := bp.queries.GetBudgetByID(ctx, budgetID)
	if err != nil {
		return "", err
	}

	if budget.UserID == utils.PgUUID(userID) {
		return PermissionOwner, nil
	}

	// Check share_access table
	access, err := bp.queries.GetShareAccessForBudgetAndUser(ctx, models.GetShareAccessForBudgetAndUserParams{
		BudgetID:    utils.PgUUID(budgetID),
		SharedWithID: utils.PgUUID(userID),
	})
	if err != nil {
		return "", err
	}

	return PermissionLevel(access.Permission), nil
}

// hasSufficientPermission checks if the user's permission meets the minimum requirement
func hasSufficientPermission(userPerm, minPerm PermissionLevel) bool {
	// Define permission hierarchy
	permissionOrder := map[PermissionLevel]int{
		PermissionOwner: 3,
		PermissionEdit:  2,
		PermissionView: 1,
	}

	userLevel := permissionOrder[userPerm]
	minLevel := permissionOrder[minPerm]

	return userLevel >= minLevel
}

// GetBudgetPermission retrieves the budget permission from the request context
func GetBudgetPermission(r *http.Request) PermissionLevel {
	if perm, ok := r.Context().Value("budgetPermission").(PermissionLevel); ok {
		return perm
	}
	return ""
}

// CanEditBudget checks if the current request can edit the budget
func CanEditBudget(r *http.Request) bool {
	perm := GetBudgetPermission(r)
	return perm == PermissionOwner || perm == PermissionEdit
}

// IsBudgetOwner checks if the current user is the budget owner
func IsBudgetOwner(r *http.Request) bool {
	return GetBudgetPermission(r) == PermissionOwner
}
