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
	http.HandleFunc("/api/auth/refresh", middleware.EnableCORS(handler.RefreshHandler))

	// /api/auth/me is protected by the JWTAuth middleware
	http.HandleFunc("/api/auth/me", middleware.EnableCORS(
		middleware.JWTAuth(handler.MeHandler),
	))

	// ===== PRINT ROUTES =====
	fmt.Println("📌 API List:")
	for _, route := range Routes {
		fmt.Printf("%s %s → %s\n", route.Method, route.Path, route.Name)
	}
}