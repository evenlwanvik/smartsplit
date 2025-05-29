package service

import (
	"context"

	"github.com/evenlwanvik/smartsplit/internal/identity/models"
	"github.com/evenlwanvik/smartsplit/internal/identity/repository"
)

// UserService handles business logic for users.
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser registers a new user.
func (s *UserService) CreateUser(ctx context.Context, user *models.CreateUser) (*models.User, error) {
	return s.repo.Create(ctx, user)
}

// GetUser fetches by ID.
func (s *UserService) GetUser(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

// ListUsers returns all users.
func (s *UserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	return s.repo.List(ctx)
}

// UpdateUser modifies user data.
func (s *UserService) UpdateUser(ctx context.Context, user *models.UpdateUser) error {
	return s.repo.Update(ctx, user)
}

// DeleteUser removes a user.
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
