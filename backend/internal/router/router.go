package router

import (
	"fmt"
	"net/http"

	"x509-pki/internal/handler"
	"x509-pki/internal/middleware"
)

func SetupRoutes() {
	// ===== AUTH =====
	http.HandleFunc("/api/auth/register", middleware.EnableCORS(handler.RegisterHandler))
	http.HandleFunc("/api/auth/login", middleware.EnableCORS(handler.LoginHandler))

	// ===== PRINT ROUTES =====
	fmt.Println("📌 API List:")
	for _, route := range Routes {
		fmt.Printf("%s %s → %s\n", route.Method, route.Path, route.Name)
	}
}