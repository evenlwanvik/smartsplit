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

func (m *Module) ListPLans(
	ctx context.Context, filters workout.Filters,
) ([]*workout.Plan, *workout.Metadata, error) {
	return m.svc.ListPLans(ctx, filters)
}

func (m *Module) ReadPlan(ctx context.Context, id int) (*workout.Plan, error) {
	return m.svc.ReadPlan(ctx, id)
}
