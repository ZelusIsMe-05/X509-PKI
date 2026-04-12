package middleware

import (
	"context"
	"net/http"
	"strings"

	"x509-pki/internal/auth"
)

// contextKey is a package-private type to avoid collisions with other packages in context.
type contextKey string

const UsernameKey contextKey = "username"

// JWTAuth is a middleware that protects routes requiring authentication.
// It validates the Bearer token from the Authorization header and injects
// the username into the request context.
func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: invalid or expired token", http.StatusUnauthorized)
			return
		}

		if claims.TokenType != "access" {
			http.Error(w, "Unauthorized: not an access token", http.StatusUnauthorized)
			return
		}

		// Inject the authenticated username into the request context
		ctx := context.WithValue(r.Context(), UsernameKey, claims.Username)
		next(w, r.WithContext(ctx))
	}
}

// GetUsernameFromContext retrieves the username injected by JWTAuth from the request context.
func GetUsernameFromContext(r *http.Request) (string, bool) {
	username, ok := r.Context().Value(UsernameKey).(string)
	return username, ok
}
