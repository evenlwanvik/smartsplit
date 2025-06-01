package identity

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/evenlwanvik/smartsplit/internal/logging"
	"github.com/evenlwanvik/smartsplit/internal/rest"
)

// UserHandler defines HTTP handlers for users.
type UserHandler struct {
	Service *UserService
}

// RegisterRoutes hooks up endpoints.
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {

	routeDefinitions := rest.RouteDefinitionList{
		{
			"GET /api/v0/identity/users",
			h.listUsersHandler,
		},
		{
			"GET /api/v0/identity/users/{id}",
			h.getUserHandler,
		},
		{
			"PUT /api/v0/identity/users/{id}",
			h.updateUserHandler,
		},
		{
			"DELETE /api/v0/identity/users/{id}",
			h.deleteUserHandler,
		},
		{
			"POST /api/v0/identity/users/register",
			h.RegisterUserHandler,
		},
	}

	for _, d := range routeDefinitions {
		mux.Handle(d.Path, d.Handler)
	}
}

func (h *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("decoding request body")
	var user CreateUser
	if err := rest.DecodeJSONFromRequest(r, &user); err != nil {
		rest.BadRequestResponse(w, r, rest.UnableToDecodeRequestBody, err)
		return
	}
	logger = logger.With(slog.Group("input", slog.Any("user", user)))

	logger.Info("creating user")
	createdUser, err := h.Service.CreateUser(ctx, &user)
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

func (h *UserHandler) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("reading users")
	users, err := h.Service.ListUsers(ctx)
	if err != nil {
		logger.Error("failed to read users", "error", err)
		rest.InternalServerErrorResponse(w, r, err)
		return
	}

	err = rest.WriteJSONResponse(w, http.StatusOK, users)
	if err != nil {
		logger.Error("failed to write response", "error", err)
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
}

func (h *UserHandler) getUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	id, err := rest.GetPathParamInt(r, "id")
	if err != nil {
		rest.UnableToGetPathParamFromRequest(w, r, "id", err)
		return
	}
	logger = logger.With(slog.Group("input", slog.Int("id", id)))

	logger.Info("reading user")
	users, err := h.Service.ReadUser(ctx, id)
	if err != nil {
		logger.Error("failed to read user", "error", err)
		switch {
		case errors.Is(err, ErrNotFound):
			rest.NotFoundResponse(w, r, err)
		default:
			rest.InternalServerErrorResponse(w, r, err)
		}
		return
	}

	err = rest.WriteJSONResponse(w, http.StatusOK, users)
	if err != nil {
		logger.Error("failed to write response", "error", err)
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
}

func (h *UserHandler) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	id, err := rest.GetPathParamInt(r, "id")
	if err != nil {
		rest.UnableToGetPathParamFromRequest(w, r, "id", err)
		return
	}

	logger.Info("decoding request body")
	var user UpdateUser
	if err := rest.DecodeJSONFromRequest(r, &user); err != nil {
		rest.BadRequestResponse(w, r, rest.UnableToDecodeRequestBody, err)
		return
	}
	logger = logger.With(slog.Group(
		"input",
		slog.Int("id", id),
		slog.Any("user", user),
	))

	logger.Info("updating user")
	updatedUser, err := h.Service.UpdateUser(ctx, id, &user)
	if err != nil {
		logger.Error("failed to update user", "error", err)
		switch {
		case errors.Is(err, ErrNotFound):
			rest.NotFoundResponse(w, r, err)
		default:
			rest.InternalServerErrorResponse(w, r, err)
		}
		return
	}

	err = rest.WriteJSONResponse(w, http.StatusOK, updatedUser)
	if err != nil {
		logger.Error("failed to write response", "error", err)
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
}

func (h *UserHandler) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	id, err := rest.GetPathParamInt(r, "id")
	if err != nil {
		rest.UnableToGetPathParamFromRequest(w, r, "id", err)
		return
	}
	logger = logger.With(slog.Group("input", slog.Int("id", id)))

	logger.Info("deleting user")
	updatedUser, err := h.Service.DeleteUser(ctx, id)
	if err != nil {
		logger.Error("failed to delete user", "error", err)
		switch {
		case errors.Is(err, ErrNotFound):
			rest.NotFoundResponse(w, r, err)
		default:
			rest.InternalServerErrorResponse(w, r, err)
		}
		return
	}

	err = rest.WriteJSONResponse(w, http.StatusOK, updatedUser)
	if err != nil {
		logger.Error("failed to write response", "error", err)
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
}
