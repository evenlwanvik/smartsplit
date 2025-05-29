package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/evenlwanvik/smartsplit/internal/identity/models"
	"github.com/evenlwanvik/smartsplit/internal/rest"

	"github.com/evenlwanvik/smartsplit/internal/identity/service"
)

// UserHandler defines HTTP handlers for users.
type UserHandler struct {
	Service *service.UserService
}

// RegisterRoutes hooks up endpoints.
func (h *UserHandler) RegisterRoutes(mux *chi.Mux) {

	routeDefinitions := rest.RouteDefinitionList{
		{
			"GET /api/v0/identity/users",
			h.createUser,
		},
		{
			"GET /api/v0/identity/users/{id}",
			h.getUser,
		},
		{
			"PUT /api/v0/identity/users/{id}",
			h.updateUser,
		},
		{
			"DELETE /api/v0/identity/users/{id}",
			h.deleteUser,
		},
		{
			"POST /api/v0/identity/users", // Rename to register?
			h.createUser,
		},
	}

	for _, d := range routeDefinitions {
		mux.Handle(d.Path, d.Handler)
	}
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var input models.CreateUser
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created, err := h.Service.CreateUser(r.Context(), &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *UserHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	u, err := h.Service.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var in models.User
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	in.ID = id
	if err := h.Service.UpdateUser(r.Context(), &in); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.Service.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
