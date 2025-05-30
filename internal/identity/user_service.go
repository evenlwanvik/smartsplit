package identity

import (
	"context"
)

// UserService handles business logic for users.
type UserService struct {
	repo *UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser registers a new user.
func (s *UserService) CreateUser(ctx context.Context, user *CreateUser) (*User, error) {
	return s.repo.Create(ctx, user)
}

// ReadUser fetches by ID.
func (s *UserService) ReadUser(ctx context.Context, id int) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

// ListUsers returns all users.
func (s *UserService) ListUsers(ctx context.Context) ([]*User, error) {
	return s.repo.List(ctx)
}

// UpdateUser modifies user data.
func (s *UserService) UpdateUser(ctx context.Context, user *UpdateUser) error {
	return s.repo.Update(ctx, user)
}

// DeleteUser removes a user.
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
