package workout

import (
	"context"

	"github.com/evenlwanvik/smartsplit/internal/workout"
)

func (m *Module) ReadMuscles(ctx context.Context) ([]*workout.Muscle, error) {
	return m.svc.ReadMuscles(ctx)
}

func (m *Module) CreatePlanWithEntries(
	ctx context.Context,
	notes string,
	musclesIds []int,
) (*workout.Plan, error) {
	return m.svc.CreatePlanWithEntries(ctx, notes, musclesIds)
}

func (m *Module) UpdatePlanEntrySets(ctx context.Context, id int, sets int) (*workout.PlanEntry, error) {
	return m.svc.UpdatePlanEntrySets(ctx, id, sets)
}

func (m *Module) DeletePlan(ctx context.Context, id int) error {
	return m.svc.DeletePlan(ctx, id)
}
