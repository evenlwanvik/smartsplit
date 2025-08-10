package workout

import (
	"context"

	"github.com/evenlwanvik/smartsplit/internal/workout"
)

func (m *Module) GetMuscles(ctx context.Context) ([]*workout.Muscle, error) {
	return m.svc.GetMuscles(ctx)
}
