package main

import (
	"fmt"
	"net/http"

	"x509-pki/internal/repository"
	"x509-pki/internal/router"
)

func main() {
	// Initialize SQLite DB and create tables if they do not exist
	repository.InitDB()

	router.SetupRoutes()

	fmt.Println("\n🚀 Server running at http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}