package workout

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/evenlwanvik/smartsplit/internal/logging"
	"github.com/evenlwanvik/smartsplit/internal/rest"
)

// Handlers defines HTTP handlers for workout.
type Handlers struct {
	Svc *Service
}

// RegisterRoutes hooks up endpoints.
func (h *Handlers) RegisterRoutes(ctx context.Context, mux *http.ServeMux) {
	logger := logging.LoggerFromContext(ctx)

	routeDefinitions := rest.RouteDefinitionList{
		{
			"GET /api/v0/workout/muscles",
			h.listMuscles,
		},
		{
			"GET /api/v0/workout/muscles/create",
			h.createMuscle,
		},
	}

	for _, d := range routeDefinitions {
		logger.Info("adding route", "route", d.Path)
		mux.Handle(d.Path, d.Handler)
	}
}

func (h *Handlers) listMuscles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	ms, err := h.Svc.ReadMuscles(r.Context())
	if err != nil {
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
	err = rest.WriteJSONResponse(w, http.StatusOK, ms)
	if err != nil {
		logger.Error("failed to write response", "error", err)
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
}

func (h *Handlers) createMuscle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("decoding request body")
	var muscle MuscleInput
	if err := rest.DecodeJSONFromRequest(r, &muscle); err != nil {
		rest.BadRequestResponse(w, r, rest.UnableToDecodeRequestBody, err)
		return
	}
	logger = logger.With(slog.Group("input", slog.Any("muscle", muscle)))

	logger.Info("creating user")
	createdUser, err := h.Svc.CreateMuscle(ctx, &muscle)
	if err != nil {
		logger.Error("failed to create user", "error", err)
		rest.InternalServerErrorResponse(w, r, err)
		return
	}

	err = rest.WriteJSONResponse(w, http.StatusCreated, createdUser)
	if err != nil {
		logger.Error("failed to write response", "error", err)
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
}
