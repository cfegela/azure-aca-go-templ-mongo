package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cfegela/azure-aca-go-templ-mongo/internal/auth"
	"github.com/cfegela/azure-aca-go-templ-mongo/internal/database"
	"github.com/cfegela/azure-aca-go-templ-mongo/internal/models"
)

func main() {
	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
	dbName := getEnv("MONGODB_DATABASE", "tasksdb")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := database.Connect(ctx, mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	userRepo := database.NewUserRepository(client, dbName)

	// Check if admin already exists
	count, err := userRepo.Count(context.Background())
	if err != nil {
		log.Fatalf("Failed to count users: %v", err)
	}

	if count > 0 {
		log.Println("Users already exist. Skipping seed.")
		return
	}

	// Create admin user
	email := getEnv("ADMIN_EMAIL", "admin@example.com")
	password := getEnv("ADMIN_PASSWORD", "admin123")
	name := getEnv("ADMIN_NAME", "Admin User")

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	admin := &models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		Name:         name,
		Role:         models.RoleAdmin,
	}

	if err := admin.Validate(); err != nil {
		log.Fatalf("Invalid admin user: %v", err)
	}

	if err := userRepo.Create(context.Background(), admin); err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	fmt.Println("===========================================")
	fmt.Println("Admin user created successfully!")
	fmt.Println("===========================================")
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("Password: %s\n", password)
	fmt.Println("===========================================")
	fmt.Println("Please change the password after first login.")
	fmt.Println("===========================================")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
