package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/joselitophala/budget-planner-backend/internal/models"
	"github.com/joselitophala/budget-planner-backend/internal/utils"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// UserIDKey is the context key for the user ID
	UserIDKey contextKey = "userID"
	// ClerkUserIDKey is the context key for the Clerk user ID (kept for compatibility)
	ClerkUserIDKey contextKey = "clerkUserID"
)

// User represents the authenticated user in the database
type User struct {
	ID          string
	ClerkUserID string
	Email       string
	Name        string
	Currency    string
}

// Middleware creates an authentication middleware using JWT
type Middleware struct {
	jwt       *JWTClient
	queries   *models.Queries
	requireAuth bool
}

// NewMiddleware creates a new auth middleware
func NewMiddleware(jwt *JWTClient, queries *models.Queries, requireAuth bool) *Middleware {
	return &Middleware{
		jwt:       jwt,
		queries:   queries,
		requireAuth: requireAuth,
	}
}

// Authenticate verifies the JWT and adds user context to the request
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log all headers for debugging
		fmt.Printf("[Auth] All headers received:\n")
		for key, values := range r.Header {
			for _, value := range values {
				fmt.Printf("[Auth]   %s: %s\n", key, value)
			}
		}
		
		// Extract token from request
		token, err := ExtractTokenFromRequest(r)
		if err != nil {
			if m.requireAuth {
				fmt.Printf("[Auth] Token extraction failed: %v\n", err)
				fmt.Printf("[Auth] Request path: %s\n", r.URL.Path)
				fmt.Printf("[Auth] Authorization header: %s\n", r.Header.Get("Authorization"))
				utils.Unauthorized(w, "Invalid or missing authorization token")
				return
			}
			// If auth is optional, continue without user context
			next.ServeHTTP(w, r)
			return
		}

		fmt.Printf("[Auth] Token extracted successfully (length: %d)\n", len(token))

		// Verify token
		userID, err := m.jwt.VerifyToken(token)
		if err != nil {
			if m.requireAuth {
				fmt.Printf("[Auth] Token verification failed: %v\n", err)
				fmt.Printf("[Auth] Request path: %s\n", r.URL.Path)
				utils.Unauthorized(w, "Invalid token")
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		fmt.Printf("[Auth] Token verified successfully for user: %s\n", userID)

		// Look up user in our database (assuming userID matches clerk_user_id)
		user, err := m.queries.GetUserByClerkID(r.Context(), userID)
		if err != nil {
			// User exists in auth system but not in our database
			// This could mean they haven't completed onboarding
			if m.requireAuth {
				fmt.Printf("[Auth] User lookup failed for clerk_user_id %s: %v\n", userID, err)
				utils.NotFound(w, "User not found in database")
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		fmt.Printf("[Auth] User found in database: %s (internal ID: %s)\n", user.ClerkUserID, user.ID)

		// Add user context to request
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIDKey, user.ID)
		ctx = context.WithValue(ctx, ClerkUserIDKey, user.ClerkUserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth creates an auth middleware that doesn't require authentication
// but will add user context if a valid token is provided
func (m *Middleware) OptionalAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return NewMiddleware(m.jwt, m.queries, false).Authenticate(next)
	}
}

// RequireAuth creates an auth middleware that requires valid authentication
func (m *Middleware) RequireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return NewMiddleware(m.jwt, m.queries, true).Authenticate(next)
	}
}

// Helper functions to get user context from request

// GetUserID retrieves the user ID from the request context
func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	return userID, ok
}

// GetClerkUserID retrieves the Clerk user ID from the request context
func GetClerkUserID(r *http.Request) (string, bool) {
	clerkUserID, ok := r.Context().Value(ClerkUserIDKey).(string)
	return clerkUserID, ok
}

// MustGetUserID retrieves the user ID or panics (for use in handlers where auth is required)
func MustGetUserID(r *http.Request) string {
	userID, ok := GetUserID(r)
	if !ok {
		panic("user ID not found in context - auth middleware may not be configured")
	}
	return userID
}

// SetUserIDInContext sets the user ID in the request context (for testing purposes)
func SetUserIDInContext(r *http.Request, userID string) context.Context {
	ctx := r.Context()
	return context.WithValue(ctx, UserIDKey, userID)
}
