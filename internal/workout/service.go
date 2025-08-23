package workout

import (
	"context"
	"log/slog"
	"time"

	"github.com/evenlwanvik/smartsplit/internal/logging"
)

type Client interface {
	ReadMuscles(ctx context.Context) ([]*Muscle, error)
	CreatePlanWithEntries(ctx context.Context, notes string, muscleIDs []int) (*Plan, error)
	UpdatePlanEntrySets(ctx context.Context, id int, sets int) (*PlanEntry, error)
	DeletePlan(ctx context.Context, id int) error
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

func (s *Service) UpdatePlanEntrySets(
	ctx context.Context,
	id int,
	sets int,
) (*PlanEntry, error) {
	planEntryPatch := PlanEntryPatch{
		ID:   id,
		Sets: sets,
	}
	entry, err := s.repo.PatchPlanEntry(ctx, planEntryPatch)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *Service) DeletePlan(ctx context.Context, id int) error {
	logger := logging.LoggerFromContext(ctx)
	logger = logger.With(slog.Group("DeletePlan", slog.Int("plan_id", id)))

	// TODO: Introduce transactions
	nDeleted, err := s.repo.DeleteManyPlanEntries(ctx, Filters{PlanID: &id})
	if err != nil {
		logger.Error("failed to delete plan entries", slog.Any("error", err))
		return err
	}
	logger.Info("deleted plan entries", slog.Int64("n_deleted", nDeleted))

	_, err = s.repo.DeletePlan(ctx, id)
	if err != nil {
		logger.Error("failed to delete plan", slog.Any("error", err))
		return err
	}
	logger.Info("deleted plan")
	return nil
}
