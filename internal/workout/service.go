package workout

import (
	"context"
)

type Client interface {
	GetMuscles(ctx context.Context) ([]*Muscle, error)
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetMuscles(ctx context.Context) ([]*Muscle, error) {
	return s.repo.SelectMuscles(ctx)
}

func (s *Service) CreateMuscle(ctx context.Context, input *MuscleInput) (*Muscle, error) {
	muscle, err := s.repo.InsertMuscle(ctx, input)
	if err != nil {
		return nil, err
	}
	return muscle, nil
}

func (s *Service) CreatePlanWithEntries(
	ctx context.Context,
	input *Plan,
) (*Plan, error) {
	plan, err := s.repo.InsertPlan(ctx, input)
	if err != nil {
		return nil, err
	}
	// TODO: Also accept entries and insert them?
	return plan, nil
}
