package workout

import (
	"context"
	"time"
)

type Client interface {
	ReadMuscles(ctx context.Context) ([]*Muscle, error)
	CreatePlanWithEntries(context.Context, string, []int) (*Plan, error)
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ReadMuscles(ctx context.Context) ([]*Muscle, error) {
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
	notes string,
	musclesIds []int,
) (*Plan, error) {
	planInput := PlanInput{
		Notes: notes,
		// Properly set user ID.
		UserID: 1,
		// TODO: Let user choose date
		Date: time.Now(),
	}
	// TODO: Setup transactions
	plan, err := s.repo.InsertPlan(ctx, planInput)
	if err != nil {
		return nil, err
	}

	var entries []*PlanEntry
	for _, muscleID := range musclesIds {
		entryInput := PlanEntry{
			MuscleID: muscleID,
			Sets:     1,
			PlanID:   plan.ID,
		}
		entry, err := s.repo.InsertPlanEntry(ctx, entryInput)
		if err != nil {
			return nil, err
		}
		entry.Muscle, err = s.repo.SelectMuscle(ctx, muscleID)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	plan.Entries = entries
	return plan, nil
}
