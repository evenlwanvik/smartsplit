package identity

import (
	"context"
)

type UserClient interface {
	ReadUser(ctx context.Context, id int) (*User, error)
}

// UserService handles business logic for usersvc.
type UserService struct {
	repo *UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser registers a new user.
func (svc *UserService) CreateUser(ctx context.Context, user *CreateUser) (*User, error) {
	return svc.repo.Create(ctx, user)
}

// ReadUser fetches by ID.
func (svc *UserService) ReadUser(ctx context.Context, id int) (*User, error) {
	return svc.repo.GetByID(ctx, id)
}

// ListUsers returns all usersvc.
func (svc *UserService) ListUsers(ctx context.Context) ([]*User, error) {
	return svc.repo.List(ctx)
}

// UpdateUser modifies user data.
func (svc *UserService) UpdateUser(ctx context.Context, id int, user *UpdateUser) (*User, error) {
	return svc.repo.Update(ctx, id, user)
}

// DeleteUser removes a user.
func (svc *UserService) DeleteUser(ctx context.Context, id int) (*User, error) {
	return svc.repo.Delete(ctx, id)
}
