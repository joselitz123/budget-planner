package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/joselitophala/budget-planner-backend/internal/utils"
)

// ErrorHandler is a middleware that recovers from panics and returns proper error responses
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				log.Printf("PANIC: %v\n%s", err, debug.Stack())

				// Return internal server error to client
				utils.InternalError(w, "An internal server error occurred")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// MethodOverride is a middleware that checks for X-HTTP-Method-Override header
// This is useful for clients that can't send certain HTTP methods
func MethodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if method := r.Header.Get("X-HTTP-Method-Override"); method != "" {
			r.Method = method
		}
		next.ServeHTTP(w, r)
	})
}

// ContentType is a middleware that ensures the request has JSON content type for POST/PUT/PATCH
func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			ct := r.Header.Get("Content-Type")
			if ct != "" && ct != "application/json" && !strings.HasPrefix(ct, "application/json") {
				utils.BadRequest(w, "Content-Type must be application/json")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
