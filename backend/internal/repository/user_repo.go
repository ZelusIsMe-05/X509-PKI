package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"x509-pki/internal/model"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// InitDB khởi tạo kết nối SQLite và tạo bảng users nếu chưa có.
// File DB được lưu tại: data/users.db
func InitDB() {
	// Tạo thư mục data/ nếu chưa tồn tại
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("❌ Không thể tạo thư mục data: %v", err)
	}

	var err error
	db, err = sql.Open("sqlite", "data/users.db")
	if err != nil {
		log.Fatalf("❌ Không thể mở SQLite DB: %v", err)
	}

	// Tạo bảng users nếu chưa có
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password TEXT NOT NULL
	);`

	if _, err := db.Exec(createTable); err != nil {
		log.Fatalf("❌ Không thể tạo bảng users: %v", err)
	}

	fmt.Println("✅ SQLite DB đã sẵn sàng tại: data/users.db")
}

// Exists kiểm tra username đã tồn tại trong DB chưa.
func Exists(username string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		log.Printf("⚠️ Lỗi khi kiểm tra username: %v", err)
		return false
	}
	return count > 0
}

// Save lưu một user mới vào DB.
func Save(user model.User) error {
	_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	if err != nil {
		log.Printf("⚠️ Lỗi khi lưu user: %v", err)
		return err
	}
	return nil
}

// GetPassword lấy password theo username.
func GetPassword(username string) (string, bool) {
	var password string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&password)
	if err == sql.ErrNoRows {
		return "", false
	}
	if err != nil {
		log.Printf("⚠️ Lỗi khi lấy password: %v", err)
		return "", false
	}
	return password, true
}