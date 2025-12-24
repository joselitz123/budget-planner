package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// Logger returns a request logger middleware
// It wraps Chi's built-in logger middleware with additional timing info
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Call the next handler
		next.ServeHTTP(ww, r)

		// Log the request details
		duration := time.Since(start)
		log.Printf(
			"%s %s %s %d %s %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			ww.Status(),
			duration,
			r.UserAgent(),
		)
	})
}

// RequestID adds a unique request ID to each request for tracing
func RequestID(next http.Handler) http.Handler {
	return middleware.RequestID(next)
}

// RealIP gets the real IP from behind a proxy
func RealIP(next http.Handler) http.Handler {
	return middleware.RealIP(next)
}
