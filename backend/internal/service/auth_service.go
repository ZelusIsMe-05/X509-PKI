package service

import (
	"errors"
	"regexp"
	"time"

	pkicrypto "x509-pki/internal/crypto"
	"x509-pki/internal/auth"
	"x509-pki/internal/model"
	"x509-pki/internal/repository"
)

// ─────────────────────────────────────────────────────────────────
// INPUT VALIDATION
// ─────────────────────────────────────────────────────────────────

// ValidateUsername checks if username meets security requirements.
func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(username) > 50 {
		return errors.New("username must be at most 50 characters long")
	}
	// Only allow alphanumeric, underscore, and hyphen
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(username) {
		return errors.New("username can only contain letters, numbers, underscores, and hyphens")
	}
	return nil
}

// ValidatePassword checks if password meets security requirements.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > 256 {
		return errors.New("password must be at most 256 characters long")
	}
	// Check for at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	// Check for at least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	// Check for at least one digit
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("password must contain at least one number")
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────
// REGISTER
// ─────────────────────────────────────────────────────────────────

// Register creates a new user whose password is hashed with Argon2id before
// being stored in the database.
func Register(user model.User) error {
	// Validate username
	if err := ValidateUsername(user.Username); err != nil {
		return err
	}

	// Validate password
	if err := ValidatePassword(user.Password); err != nil {
		return err
	}

	if repository.Exists(user.Username) {
		return errors.New("username already exists")
	}

	// Generate a cryptographically random salt
	salt, err := pkicrypto.GenerateSalt()
	if err != nil {
		return errors.New("failed to generate salt")
	}

	// Derive Argon2id hash — all parameters are defined in the crypto package
	passwordHash, err := pkicrypto.HashPassword(user.Password, salt)
	if err != nil {
		return errors.New("failed to hash password")
	}

	if err := repository.SaveHashed(user.Username, passwordHash, salt); err != nil {
		return errors.New("failed to save user")
	}

	return nil
}

// ─────────────────────────────────────────────────────────────────
// LOGIN
// ─────────────────────────────────────────────────────────────────

// LoginResult holds the JWT token pair and username returned after login.
type LoginResult struct {
	AccessToken  string
	RefreshToken string
	Username     string
}

// Login authenticates a user and returns a JWT token pair on success.
func Login(user model.User) (*LoginResult, error) {
	// Fetch stored user record from DB
	userDB, exists := repository.GetUserByUsername(user.Username)
	if !exists {
		return nil, errors.New("invalid username or password")
	}

	// Re-derive Argon2id hash with stored salt and compare via constant-time check
	if !pkicrypto.VerifyPassword(user.Password, userDB.Salt, userDB.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	// Issue JWT access + refresh token pair
	accessToken, refreshToken, err := auth.GenerateTokenPair(user.Username)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Persist the refresh token hash in DB
	expiresAt := time.Now().Add(auth.RefreshTokenDuration)
	if err := repository.SaveRefreshToken(refreshToken, user.Username, expiresAt); err != nil {
		return nil, errors.New("failed to save session")
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Username:     user.Username,
	}, nil
}

// ─────────────────────────────────────────────────────────────────
// REFRESH TOKEN
// ─────────────────────────────────────────────────────────────────

// RefreshToken validates the old refresh token and issues a new token pair
// using single-use rotation (old token is revoked immediately).
func RefreshToken(oldRefreshToken string) (*LoginResult, error) {
	// 1. Validate JWT signature and expiry
	claims, err := auth.ValidateToken(oldRefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// 2. Confirm this is a refresh token (not an access token)
	if claims.TokenType != "refresh" {
		return nil, errors.New("not a refresh token")
	}

	// 3. Verify the token exists in DB and has not been revoked
	username, _, found := repository.GetRefreshToken(oldRefreshToken)
	if !found || username != claims.Username {
		return nil, errors.New("refresh token not found or revoked")
	}

	// 4. Revoke the old refresh token (single-use rotation)
	if err := repository.DeleteRefreshToken(oldRefreshToken); err != nil {
		return nil, errors.New("failed to revoke old token")
	}

	// 5. Generate a new token pair
	accessToken, newRefreshToken, err := auth.GenerateTokenPair(claims.Username)
	if err != nil {
		return nil, errors.New("failed to generate new tokens")
	}

	// 6. Persist the new refresh token
	expiresAt := time.Now().Add(auth.RefreshTokenDuration)
	if err := repository.SaveRefreshToken(newRefreshToken, claims.Username, expiresAt); err != nil {
		return nil, errors.New("failed to save new session")
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		Username:     claims.Username,
	}, nil
}
