package repository

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"x509-pki/internal/crypto"
	"x509-pki/internal/model"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// InitDB opens the SQLite connection and creates all required tables if they do not exist.
// The database file is stored at: data/users.db
func InitDB() {
	// Create the data/ directory if it does not exist
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("❌ Failed to create data directory: %v", err)
	}

	var err error
	db, err = sql.Open("sqlite", "data/users.db")
	if err != nil {
		log.Fatalf("❌ Failed to open SQLite DB: %v", err)
	}

	// users table: stores Argon2id hash and salt per user
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		username      TEXT PRIMARY KEY,
		password_hash TEXT NOT NULL,
		salt          TEXT NOT NULL,
		role          TEXT DEFAULT 'client'
	);`

	if _, err := db.Exec(createUsersTable); err != nil {
		log.Fatalf("❌ Failed to create users table: %v", err)
	}

	// Attempt to add role column to existing table (ignores error if column already exists)
	db.Exec("ALTER TABLE users ADD COLUMN role TEXT DEFAULT 'client'")

	// refresh_tokens table: stores SHA-256 hashed refresh tokens for rotation/revocation
	createRefreshTokensTable := `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		token_hash TEXT PRIMARY KEY,
		username   TEXT NOT NULL,
		expires_at DATETIME NOT NULL
	);`

	if _, err := db.Exec(createRefreshTokensTable); err != nil {
		log.Fatalf("❌ Failed to create refresh_tokens table: %v", err)
	}

	fmt.Println("✅ SQLite DB ready at: data/users.db")

	// Seed admin user if it does not exist
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	if count == 0 {
		salt, err := crypto.GenerateSalt()
		if err != nil {
			log.Fatalf("❌ Failed to generate salt for admin: %v", err)
		}
		hash, err := crypto.HashPassword("Admin@x509-pki", salt)
		if err != nil {
			log.Fatalf("❌ Failed to hash password for admin: %v", err)
		}
		_, err = db.Exec(
			"INSERT INTO users (username, password_hash, salt, role) VALUES (?, ?, ?, ?)",
			"admin", hash, salt, "admin",
		)
		if err != nil {
			log.Fatalf("❌ Failed to seed admin user: %v", err)
		}
		fmt.Println("✅ Seeded admin user (admin / Admin@x509-pki)")
	}

	// 🕒 Start background job to prune expired refresh tokens automatically
	go startTokenCleanupDaemon()
}

// startTokenCleanupDaemon runs an infinite loop in a background goroutine,
// triggering DeleteExpiredRefreshTokens once every hour to prevent DB bloat.
func startTokenCleanupDaemon() {
	// First run immediately
	DeleteExpiredRefreshTokens()
	
	// Then trigger every 1 hour
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		DeleteExpiredRefreshTokens()
	}
}


// ─────────────────────────────────────────────────────────────────
// USER REPOSITORY
// ─────────────────────────────────────────────────────────────────

// Exists returns true if the given username already exists in the DB.
func Exists(username string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		log.Printf("⚠️ Error checking username: %v", err)
		return false
	}
	return count > 0
}

// SaveHashed inserts a new user with a pre-computed password hash and salt.
func SaveHashed(username, passwordHash, salt string) error {
	role := "client"
	_, err := db.Exec(
		"INSERT INTO users (username, password_hash, salt, role) VALUES (?, ?, ?, ?)",
		username, passwordHash, salt, role,
	)
	if err != nil {
		log.Printf("⚠️ Error saving user: %v", err)
		return err
	}
	return nil
}

// GetUserByUsername retrieves a UserDB record by username.
func GetUserByUsername(username string) (model.UserDB, bool) {
	var u model.UserDB
	err := db.QueryRow(
		"SELECT username, password_hash, salt, role FROM users WHERE username = ?",
		username,
	).Scan(&u.Username, &u.PasswordHash, &u.Salt, &u.Role)

	if err == sql.ErrNoRows {
		return model.UserDB{}, false
	}
	if err != nil {
		log.Printf("⚠️ Error fetching user: %v", err)
		return model.UserDB{}, false
	}
	return u, true
}

// ─────────────────────────────────────────────────────────────────
// REFRESH TOKEN REPOSITORY
// ─────────────────────────────────────────────────────────────────

// hashToken returns the SHA-256 hex digest of a token string.
// Raw token values are never stored in the DB.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}

// SaveRefreshToken stores the SHA-256 hash of a refresh token along with its owner and expiry.
func SaveRefreshToken(token, username string, expiresAt time.Time) error {
	tokenHash := hashToken(token)
	_, err := db.Exec(
		"INSERT INTO refresh_tokens (token_hash, username, expires_at) VALUES (?, ?, ?)",
		tokenHash, username, expiresAt,
	)
	if err != nil {
		log.Printf("⚠️ Error saving refresh token: %v", err)
		return err
	}
	return nil
}

// GetRefreshToken looks up a refresh token in the DB by its raw value.
// Returns (username, expiresAt, found).
func GetRefreshToken(token string) (string, time.Time, bool) {
	tokenHash := hashToken(token)
	var username string
	var expiresAt time.Time

	err := db.QueryRow(
		"SELECT username, expires_at FROM refresh_tokens WHERE token_hash = ?",
		tokenHash,
	).Scan(&username, &expiresAt)

	if err == sql.ErrNoRows {
		return "", time.Time{}, false
	}
	if err != nil {
		log.Printf("⚠️ Error fetching refresh token: %v", err)
		return "", time.Time{}, false
	}
	return username, expiresAt, true
}

// DeleteRefreshToken removes a refresh token from the DB (used during rotation or logout).
func DeleteRefreshToken(token string) error {
	tokenHash := hashToken(token)
	_, err := db.Exec("DELETE FROM refresh_tokens WHERE token_hash = ?", tokenHash)
	if err != nil {
		log.Printf("⚠️ Error deleting refresh token: %v", err)
		return err
	}
	return nil
}

// DeleteExpiredRefreshTokens cleans up all refresh tokens that have passed their expiry.
func DeleteExpiredRefreshTokens() {
	_, err := db.Exec("DELETE FROM refresh_tokens WHERE expires_at < ?", time.Now())
	if err != nil {
		log.Printf("⚠️ Error pruning expired refresh tokens: %v", err)
	}
}