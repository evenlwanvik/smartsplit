package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/evenlwanvik/smartsplit/internal/identity/models"
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
func (r *UserRepository) Create(ctx context.Context, user *models.CreateUser) (*models.User, error) {
	query := `
	INSERT INTO identity.user (
		email, first_name, last_name, username, password_hash
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, email, first_name, last_name, username, password_hash, created_at, updated_at
	`

	var u models.User

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
func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
	SELECT id, email, first_name, last_name, username, password_hash, created_at, updated_at
	FROM identity.user
	WHERE id = $1
	`
	var u models.User
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
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &u, err
}

// List retrieves all users from the identity.user table.
func (r *UserRepository) List(ctx context.Context) ([]*models.User, error) {
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

	var users []*models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.FirstName,
			&u.LastName,
			&u.Username,
			&u.PasswordHash,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
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
func (r *UserRepository) Update(ctx context.Context, u *models.UpdateUser) error {
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
	`
	res, err := r.db.ExecContext(
		ctx,
		query,
		u.ID,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Username,
		u.PasswordHash,
	)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete removes a user and all associated workout data.
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// TODO: also remove associated workout records in workout schema

	deleteQuery := `
	DELETE FROM identity.user WHERE id = $1
	RETURNING id, email, first_name, last_name, username, password_hash, created_at, updated_at
	`

	var u models.User

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
		return err
	}

	if u.ID == 0 {
		return ErrNotFound
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
