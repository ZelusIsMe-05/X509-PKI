package router

import (
	"fmt"
	"net/http"

	"x509-pki/internal/handler"
	"x509-pki/internal/middleware"
)

func SetupRoutes() {
	// ===== AUTH =====
	// Register and login endpoints have rate limiting to prevent brute force
	http.HandleFunc("/api/auth/register", middleware.EnableCORS(
		middleware.RateLimit(handler.RegisterHandler),
	))
	http.HandleFunc("/api/auth/login", middleware.EnableCORS(
		middleware.RateLimit(handler.LoginHandler),
	))
	http.HandleFunc("/api/auth/refresh", middleware.EnableCORS(handler.RefreshHandler))

	// /api/auth/me is protected by the JWTAuth middleware
	http.HandleFunc("/api/auth/me", middleware.EnableCORS(
		middleware.JWTAuth(handler.MeHandler),
	))

	// /api/auth/logout is protected by the JWTAuth middleware
	http.HandleFunc("/api/auth/logout", middleware.EnableCORS(
		middleware.JWTAuth(handler.LogoutHandler),
	))

	// ===== PRINT ROUTES =====
	fmt.Println("📌 API List:")
	for _, route := range Routes {
		fmt.Printf("%s %s → %s\n", route.Method, route.Path, route.Name)
	}
}