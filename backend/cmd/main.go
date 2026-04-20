package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"x509-pki/internal/auth"
	"x509-pki/internal/repository"
	"x509-pki/internal/router"
)

// init loads environment variables from .env file before main() runs
func init() {
	// Try multiple paths to find .env
	paths := []string{
		".env",                                          // Current directory
		filepath.Join("..", ".env"),                     // Parent directory (from backend/)
		filepath.Join(os.Getenv("HOME"), "X509-PKI", ".env"), // User home
	}

	var err error
	for _, path := range paths {
		if err = godotenv.Load(path); err == nil {
			log.Printf("✅ Loaded .env from: %s\n", path)
			return
		}
	}

	log.Println("⚠️ No .env file found, using environment variables from system")
}

func main() {
	// Initialize SQLite DB and create tables if they do not exist
	repository.InitDB()

	// Initialize JWT secret after .env is loaded
	auth.InitJWTSecret()

	router.SetupRoutes()

	fmt.Println("\n🚀 Server running at http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}