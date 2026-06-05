package dto

import "time"

type TestResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TestRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsActive bool   `json:"is_active"`
}

type CreateTestRecord struct {
	Username     string
	PasswordHash string
	IsActive     bool
}
