package auth

import (
	"time"
)

type User struct {
	ID           int       `json:"id,omitempty"`
	Email        string    `json:"email,omitempty"`
	FirstName    string    `json:"first_name,omitempty"`
	LastName     string    `json:"last_name,omitempty"`
	Username     string    `json:"username,omitempty"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateUser struct {
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"user_name"`
	PasswordHash string `json:"password_hash"`
}

type UpdateUser struct {
	Email        *string `json:"email"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	Username     *string `json:"user_name"`
	PasswordHash *string `json:"password_hash"`
}

type RegisterUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
