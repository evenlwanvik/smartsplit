package identity

import (
	"time"
)

type User struct {
	ID           int
	Email        string
	FirstName    string
	LastName     string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateUser struct {
	Email        string
	FirstName    string
	LastName     string
	Username     string
	PasswordHash string
}

type UpdateUser struct {
	ID           int
	Email        *string
	FirstName    *string
	LastName     *string
	Username     *string
	PasswordHash *string
}
