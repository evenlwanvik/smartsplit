package workout

import (
	"context"

	"github.com/evenlwanvik/smartsplit/internal/workout"
)

func (m *Module) ReadMuscles(ctx context.Context) ([]*workout.Muscle, error) {
	return m.svc.ReadMuscles(ctx)
}
