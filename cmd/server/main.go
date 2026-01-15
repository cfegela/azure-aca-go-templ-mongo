package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cfegela/azure-aca-go-templ-mongo/internal/auth"
	"github.com/cfegela/azure-aca-go-templ-mongo/internal/database"
	"github.com/cfegela/azure-aca-go-templ-mongo/internal/handlers"
)

func main() {
	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
	dbName := getEnv("MONGODB_DATABASE", "tasksdb")
	port := getEnv("PORT", "8080")
	jwtSecret := getEnv("JWT_SECRET", "default-secret-change-in-production")
	jwtExpiryStr := getEnv("JWT_EXPIRY", "24h")

	jwtExpiry, err := time.ParseDuration(jwtExpiryStr)
	if err != nil {
		log.Printf("Invalid JWT_EXPIRY, using default 24h: %v", err)
		jwtExpiry = 24 * time.Hour
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := database.Connect(ctx, mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	log.Println("Connected to MongoDB")

	// Initialize repositories
	taskRepo := database.NewTaskRepository(client, dbName)
	userRepo := database.NewUserRepository(client, dbName)
	inviteRepo := database.NewInviteRepository(client, dbName)

	// Create indexes
	if err := userRepo.CreateIndexes(context.Background()); err != nil {
		log.Printf("Warning: Failed to create user indexes: %v", err)
	}
	if err := inviteRepo.CreateIndexes(context.Background()); err != nil {
		log.Printf("Warning: Failed to create invite indexes: %v", err)
	}

	// Initialize auth config
	authConfig := &auth.Config{
		JWTSecret: jwtSecret,
	}

	// Initialize handlers
	taskHandler := handlers.NewTaskHandler(taskRepo)
	authHandler := handlers.NewAuthHandler(userRepo, inviteRepo, authConfig, jwtExpiry)
	pageHandler := handlers.NewPageHandler(taskRepo, userRepo, inviteRepo)

	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Public routes
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			pageHandler.ShowLogin(w, r)
		} else if r.Method == http.MethodPost {
			authHandler.HandleLogin(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/logout", authHandler.HandleLogout)
	mux.HandleFunc("/register/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			pageHandler.ShowRegister(w, r)
		} else if r.Method == http.MethodPost {
			authHandler.HandleRegister(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Auth form handler
	mux.Handle("/api/login", http.HandlerFunc(authHandler.HandleLogin))

	// Protected page routes
	mux.Handle("/", auth.RequireAuth(authConfig)(http.HandlerFunc(pageHandler.ShowDashboard)))
	mux.Handle("/tasks/new", auth.RequireAuth(authConfig)(http.HandlerFunc(pageHandler.ShowTaskForm)))
	mux.Handle("/tasks/", auth.RequireAuth(authConfig)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/tasks/" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if r.Method == http.MethodGet && hasEditSuffix(r.URL.Path) {
			pageHandler.ShowEditForm(w, r)
		} else if r.Method == http.MethodPost && hasDeleteSuffix(r.URL.Path) {
			pageHandler.DeleteTask(w, r)
		} else if r.Method == http.MethodPost {
			pageHandler.UpdateTask(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})))
	mux.Handle("/tasks", auth.RequireAuth(authConfig)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			pageHandler.CreateTask(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Admin routes
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/admin/invites", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			pageHandler.ShowInvites(w, r)
		} else if r.Method == http.MethodPost {
			pageHandler.CreateInvite(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.Handle("/admin/", auth.RequireAdmin(authConfig)(auth.RequireAuth(authConfig)(adminMux)))

	// API routes (protected)
	mux.Handle("/api/tasks", auth.RequireAuth(authConfig)(http.HandlerFunc(taskHandler.HandleTasks)))
	mux.Handle("/api/tasks/", auth.RequireAuth(authConfig)(http.HandlerFunc(taskHandler.HandleTasks)))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func hasEditSuffix(path string) bool {
	return len(path) > 5 && path[len(path)-5:] == "/edit"
}

func hasDeleteSuffix(path string) bool {
	return len(path) > 7 && path[len(path)-7:] == "/delete"
}
