package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"x509-pki/internal/auth"
	"x509-pki/internal/model"
	"x509-pki/internal/service"
)

// ─────────────────────────────────────────────────────────────────
// REGISTER
// ─────────────────────────────────────────────────────────────────

// RegisterHandler handles new account creation.
// POST /api/auth/register
// Body: { "username": "...", "password": "..." }
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	err := service.Register(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}

// ─────────────────────────────────────────────────────────────────
// LOGIN
// ─────────────────────────────────────────────────────────────────

// LoginHandler authenticates a user and returns a JWT token pair.
// POST /api/auth/login
// Body:     { "username": "...", "password": "..." }
// Response: { "access_token": "...", "refresh_token": "...", "username": "..." }
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := service.Login(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"username":      result.Username,
	})
}

// ─────────────────────────────────────────────────────────────────
// REFRESH TOKEN
// ─────────────────────────────────────────────────────────────────

// RefreshHandler issues a new token pair when the access token has expired.
// POST /api/auth/refresh
// Body:     { "refresh_token": "..." }
// Response: { "access_token": "...", "refresh_token": "...", "username": "..." }
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		http.Error(w, "Invalid request body: refresh_token required", http.StatusBadRequest)
		return
	}

	result, err := service.RefreshToken(body.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"username":      result.Username,
	})
}

// ─────────────────────────────────────────────────────────────────
// ME — current session info
// ─────────────────────────────────────────────────────────────────

// MeHandler returns the username of the currently authenticated user.
// GET /api/auth/me
// Header:   Authorization: Bearer <access_token>
// Response: { "username": "..." }
func MeHandler(w http.ResponseWriter, r *http.Request) {
	// Extract Bearer token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	if claims.TokenType != "access" {
		http.Error(w, "Not an access token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"username": claims.Username,
	})
}