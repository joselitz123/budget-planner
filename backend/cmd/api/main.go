package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/joselitophala/budget-planner-backend/internal/auth"
	"github.com/joselitophala/budget-planner-backend/internal/config"
	"github.com/joselitophala/budget-planner-backend/internal/database"
	"github.com/joselitophala/budget-planner-backend/internal/handlers"
	_ "github.com/joselitophala/budget-planner-backend/internal/models"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.NewConnection(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations (optional - comment out if using external migration tool)
	// err = database.Migrate(context.Background(), db, "sql/schema")
	// if err != nil {
	// 	log.Fatalf("Failed to run migrations: %v", err)
	// }

	// Initialize JWT client with Clerk domain for JWKS verification
	jwtClient, err := auth.NewJWTClient(cfg.ClerkDomain)
	if err != nil {
		log.Fatalf("Failed to initialize JWT client: %v", err)
	}

	// Create router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))
	r.Use(chiMiddleware.AllowContentType("application/json"))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Auth middleware
	authMiddleware := auth.NewMiddleware(jwtClient, db.Queries, true)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db.Queries, jwtClient)
	userHandler := handlers.NewUserHandler(db.Queries)
	categoryHandler := handlers.NewCategoryHandler(db.Queries)
	budgetHandler := handlers.NewBudgetHandler(db.Queries)
	transactionHandler := handlers.NewTransactionHandler(db.Queries)
	syncHandler := handlers.NewSyncHandler(db.Queries)
	paymentMethodHandler := handlers.NewPaymentMethodHandler(db.Queries)
	reflectionHandler := handlers.NewReflectionHandler(db.Queries)
	sharingHandler := handlers.NewSharingHandler(db.Queries)
	analyticsHandler := handlers.NewAnalyticsHandler(db.Queries)

	// Health check endpoint (no auth required)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Auth routes (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/logout", authHandler.Logout)
			r.Post("/refresh", authHandler.RefreshToken)
			r.Post("/onboarding", authHandler.CompleteOnboarding)
			r.Get("/me", authHandler.GetCurrentUser) // Requires auth in handler
		})

		// Protected routes (require authentication)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth())

			// Users routes (self-only)
			r.Route("/users", func(r chi.Router) {
				r.Get("/me", userHandler.GetProfile)
				r.Put("/me", userHandler.UpdateProfile)
				r.Delete("/me", userHandler.DeleteAccount)
			})

			// Categories routes
			r.Route("/categories", func(r chi.Router) {
				r.Get("/", categoryHandler.ListCategories)
				r.Get("/system", categoryHandler.GetSystemCategories)
				r.Post("/", categoryHandler.CreateCategory)
				r.Put("/{id}", categoryHandler.UpdateCategory)
				r.Delete("/{id}", categoryHandler.DeleteCategory)
			})

			// Budgets routes
			r.Route("/budgets", func(r chi.Router) {
				r.Get("/", budgetHandler.ListBudgets)
				r.Post("/", budgetHandler.CreateBudget)
				r.Get("/{month}", budgetHandler.GetBudgetByMonth)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", budgetHandler.GetBudget)
					r.Put("/", budgetHandler.UpdateBudget)
					r.Delete("/", budgetHandler.DeleteBudget)
					r.Get("/categories", budgetHandler.GetBudgetCategories)
					r.Post("/categories", budgetHandler.AddBudgetCategory)
				})
				r.Put("/categories/{categoryId}", budgetHandler.UpdateBudgetCategory)
				r.Delete("/categories/{categoryId}", budgetHandler.RemoveBudgetCategory)
			})

			// Transactions routes
			r.Route("/transactions", func(r chi.Router) {
				r.Get("/", transactionHandler.ListTransactions)
				r.Post("/", transactionHandler.CreateTransaction)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", transactionHandler.GetTransaction)
					r.Put("/", transactionHandler.UpdateTransaction)
					r.Delete("/", transactionHandler.DeleteTransaction)
				})
			})

			// Payment Methods routes
			r.Route("/payment-methods", func(r chi.Router) {
				r.Get("/", paymentMethodHandler.ListPaymentMethods)
				r.Post("/", paymentMethodHandler.CreatePaymentMethod)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", paymentMethodHandler.GetPaymentMethod)
					r.Put("/", paymentMethodHandler.UpdatePaymentMethod)
					r.Delete("/", paymentMethodHandler.DeletePaymentMethod)
				})
			})

			// Sync routes
			r.Route("/sync", func(r chi.Router) {
				r.Post("/push", syncHandler.Push)
				r.Post("/pull", syncHandler.Pull)
				r.Get("/status", syncHandler.GetStatus)
				r.Post("/resolve-conflict", syncHandler.ResolveConflict)
			})

			// Reflections routes
			r.Route("/reflections", func(r chi.Router) {
				r.Get("/month/{month}", reflectionHandler.GetReflectionByMonth)
				r.Post("/", reflectionHandler.CreateReflection)
				r.Route("/{id}", func(r chi.Router) {
					r.Put("/", reflectionHandler.UpdateReflection)
					r.Delete("/", reflectionHandler.DeleteReflection)
				})
				r.Get("/templates", reflectionHandler.ListReflectionTemplates)
			})

			// Sharing routes
			r.Route("/sharing", func(r chi.Router) {
				r.Post("/invite", sharingHandler.CreateShareInvitation)
				r.Get("/invitations", sharingHandler.GetMyInvitations)
				r.Route("/invitations/{id}", func(r chi.Router) {
					r.Put("/respond", sharingHandler.RespondToInvitation)
					r.Delete("/", sharingHandler.CancelInvitation)
				})
				r.Get("/budgets/{budgetId}", sharingHandler.GetBudgetSharing)
				r.Route("/access/{id}", func(r chi.Router) {
					r.Delete("/", sharingHandler.RemoveAccess)
				})
				r.Get("/shared-with-me", sharingHandler.GetSharedBudgets)
			})

			// Analytics routes
			r.Route("/analytics", func(r chi.Router) {
				r.Get("/dashboard/{month}", analyticsHandler.GetDashboard)
				r.Get("/spending/{month}", analyticsHandler.GetSpendingReport)
				r.Get("/trends", analyticsHandler.GetTrends)
				r.Get("/category/{categoryId}", analyticsHandler.GetCategoryReport)
			})
		})
	})

	// Start server
	addr := cfg.GetServerAddr()
	log.Printf("Server starting on %s", addr)

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	go func() {
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

// Helper functions for response writing
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// In a real implementation, you'd use json.NewEncoder
	fmt.Fprintf(w, `{"status":"%s"}`, payload)
}
