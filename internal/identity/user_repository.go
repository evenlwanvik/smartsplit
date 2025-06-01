package identity

import (
	"context"
	"database/sql"
	"errors"
)

var (
	// ErrNotFound is returned when a user cannot be found in the database.
	ErrNotFound = errors.New("user not found")
)

// UserRepository provides access to the users store.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user into the identity.user table.
func (r *UserRepository) Create(ctx context.Context, user *CreateUser) (*User, error) {
	query := `
	INSERT INTO identity.user (
		email, first_name, last_name, username, password_hash
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, email, first_name, last_name, username, password_hash, created_at, updated_at
	`

	var u User

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Username,
		user.PasswordHash,
	).Scan(
		&u.ID,
		&u.Email,
		&u.FirstName,
		&u.LastName,
		&u.Username,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, err
}

// GetByID fetches a user by ID.
func (r *UserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	query := `
	SELECT id, email, first_name, last_name, username, password_hash, created_at, updated_at
	FROM identity.user
	WHERE id = $1
	`
	var u User
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(
			&u.ID,
			&u.Email,
			&u.FirstName,
			&u.LastName,
			&u.Username,
			&u.PasswordHash,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, err
}

// List retrieves all users from the identity.user table.
func (r *UserRepository) List(ctx context.Context) ([]*User, error) {
	query := `
	SELECT id, email, first_name, last_name, username, password_hash, created_at, updated_at
	FROM identity.user
	ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.FirstName,
			&u.LastName,
			&u.Username,
			&u.PasswordHash,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// Update modifies an existing user's details.
func (r *UserRepository) Update(ctx context.Context, id int, user *UpdateUser) (*User, error) {
	query := `
	UPDATE identity.user
	SET
		email = $2,
		first_name = $3,
		last_name = $4,
		username = $5,
		password_hash = $6,
		updated_at = NOW()
	WHERE id = $1
	RETURNING id, email, first_name, last_name, username, password_hash, created_at, updated_at
	`

	var u User

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Username,
		user.PasswordHash,
	).Scan(
		&u.ID,
		&u.Email,
		&u.FirstName,
		&u.LastName,
		&u.Username,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}

// Delete removes a user and all associated workout data.
func (r *UserRepository) Delete(ctx context.Context, id int) (*User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// TODO: also remove associated workout records in workout schema

	deleteQuery := `
	DELETE FROM identity.user WHERE id = $1
	RETURNING id, email, first_name, last_name, username, password_hash, created_at, updated_at
	`

	var u User

	err = tx.QueryRowContext(ctx, deleteQuery, id).Scan(
		&u.ID,
		&u.Email,
		&u.FirstName,
		&u.LastName,
		&u.Username,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return &u, err
	}
	return &u, nil
}
