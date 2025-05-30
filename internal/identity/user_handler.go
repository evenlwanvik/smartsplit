package identity

import (
	"encoding/json"
	"net/http"

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

	// Decode the request body into a RegisterUser struct
	var user CreateUser
	if err := rest.ReadJSONFromRequest(r, &user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdUser, err := h.Service.CreateUser(ctx, &user)
	if err != nil {
		rest.InternalServerError(w, err)
		return
	}

	err = rest.WriteJSONResponse(w, http.StatusCreated, createdUser)
	if err != nil {
		rest.InternalServerError(w, err)
		return
	}
}

func (h *UserHandler) listUsersHandler(w http.ResponseWriter, r *http.Request) {
}

func (h *UserHandler) getUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract user ID from the request URL
	id, err := rest.GetPathParamInt(r, "id")

	// Call the service to get the user details
	user, err := h.Service.ReadUser(ctx, id)
	if err != nil {
		rest.InternalServerError(w, err)
		return
	}

	// Respond with the user details
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) updateUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) deleteUserHandler(w http.ResponseWriter, r *http.Request) {

}
