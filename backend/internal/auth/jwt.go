package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// In production, load the secret from an environment variable.
const (
	jwtSecret            = "x509-pki-super-secret-key-2024"
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
)

// Claims is the JWT payload carrying the username and token type.
type Claims struct {
	Username  string `json:"username"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// GenerateTokenPair issues a new access token (15 min) and refresh token (7 days).
func GenerateTokenPair(username string) (accessToken, refreshToken string, err error) {
	accessToken, err = generateToken(username, "access", AccessTokenDuration)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = generateToken(username, "refresh", RefreshTokenDuration)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// generateToken creates a signed JWT of the given type and TTL.
func generateToken(username, tokenType string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		Username:  username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "x509-pki",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// ValidateToken parses and verifies a JWT token string.
// Returns the Claims on success, or an error if the token is invalid or expired.
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Enforce HMAC signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
