package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/evenlwanvik/smartsplit/internal/logging"
	"github.com/evenlwanvik/smartsplit/internal/rest"
)

// WebHandlers defines HTTP handlers for web pages.
type WebHandlers struct {
	Service *WebService
}

// RegisterRoutes hooks up endpoints.
func (h *WebHandlers) RegisterRoutes(ctx context.Context, mux *http.ServeMux) {
	logger := logging.LoggerFromContext(ctx)

	routeDefinitions := rest.RouteDefinitionList{
		{
			"GET /dashboard",
			h.dashboardPage,
		},
		{
			"POST /plans/new",
			h.newPlanPage,
		},
	}

	for _, d := range routeDefinitions {
		logger.Info("adding route", "route", d.Path)
		mux.Handle(d.Path, d.Handler)
	}
}

func (h *WebHandlers) dashboardPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := h.Service.Dashboard(ctx, w); err != nil {
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
}

func parseIds(idInput []string) ([]int, error) {
	var parsedIds []int
	for _, id := range idInput {
		parsedId, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		parsedIds = append(parsedIds, parsedId)
	}
	return parsedIds, nil
}

func (h *WebHandlers) newPlanPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	muscles := r.Form["muscles"]
	if len(muscles) == 0 {
		http.Error(w, "pick at least one muscle", http.StatusBadRequest)
		return
	}
	muscleIds, err := parseIds(muscles)
	if err != nil {
		rest.InternalServerErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = h.Service.NewPlan(ctx, w, muscleIds)
	if err != nil {
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
}
