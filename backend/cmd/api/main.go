package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", handleLogin)
			r.Post("/logout", handleLogout)
			r.Get("/me", handleGetCurrentUser)
			r.Post("/refresh", handleRefreshToken)
		})

		// Users routes (self-only)
		r.Route("/users", func(r chi.Router) {
			r.Get("/me", handleGetUserProfile)
			r.Put("/me", handleUpdateUserProfile)
			r.Delete("/me", handleDeleteUser)
		})

		// Budgets routes
		r.Route("/budgets", func(r chi.Router) {
			r.Get("/", handleListBudgets)
			r.Post("/", handleCreateBudget)
			r.Get("/{month}", handleGetBudgetByMonth)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handleGetBudget)
				r.Put("/", handleUpdateBudget)
				r.Delete("/", handleDeleteBudget)
				r.Post("/categories", handleAddBudgetCategory)
				r.Put("/categories/{categoryId}", handleUpdateBudgetCategory)
				r.Delete("/categories/{categoryId}", handleRemoveBudgetCategory)
				r.Put("/sharing-settings", handleUpdateSharingSettings)
			})
			r.Get("/shared", handleGetSharedBudgets)
		})

		// Continue with other routes...
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// Placeholder handlers (to be implemented)
func handleLogin(w http.ResponseWriter, r *http.Request)              {}
func handleLogout(w http.ResponseWriter, r *http.Request)             {}
func handleGetCurrentUser(w http.ResponseWriter, r *http.Request)     {}
func handleRefreshToken(w http.ResponseWriter, r *http.Request)       {}
func handleGetUserProfile(w http.ResponseWriter, r *http.Request)     {}
func handleUpdateUserProfile(w http.ResponseWriter, r *http.Request)  {}
func handleDeleteUser(w http.ResponseWriter, r *http.Request)         {}
func handleListBudgets(w http.ResponseWriter, r *http.Request)        {}
func handleCreateBudget(w http.ResponseWriter, r *http.Request)       {}
func handleGetBudgetByMonth(w http.ResponseWriter, r *http.Request)   {}
func handleGetBudget(w http.ResponseWriter, r *http.Request)          {}
func handleUpdateBudget(w http.ResponseWriter, r *http.Request)       {}
func handleDeleteBudget(w http.ResponseWriter, r *http.Request)       {}
func handleAddBudgetCategory(w http.ResponseWriter, r *http.Request)  {}
func handleUpdateBudgetCategory(w http.ResponseWriter, r *http.Request) {}
func handleRemoveBudgetCategory(w http.ResponseWriter, r *http.Request) {}
func handleUpdateSharingSettings(w http.ResponseWriter, r *http.Request) {}
func handleGetSharedBudgets(w http.ResponseWriter, r *http.Request)   {}
