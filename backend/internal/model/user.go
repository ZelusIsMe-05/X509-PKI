package model

// User is the request body struct used for login and registration endpoints.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserDB represents a user record as stored in the database.
type UserDB struct {
	Username     string
	PasswordHash string
	Salt         string
}