package web

import (
	"context"
	"net/http"

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
			"GET /",
			h.indexPage,
		},
		{
			"GET /dashboard",
			h.dashboardPage,
		},
		{
			"GET /counters",
			h.counterPage,
		},
		{
			"GET /muscles",
			h.musclesPage,
		},
	}

	for _, d := range routeDefinitions {
		logger.Info("adding route", "route", d.Path)
		mux.Handle(d.Path, d.Handler)
	}
}

func (h *WebHandlers) indexPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := h.Service.Index(w)
	if err != nil {
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
}

func (h *WebHandlers) dashboardPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := h.Service.Dashboard(w); err != nil {
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
}

func (h *WebHandlers) counterPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := h.Service.Dashboard(w); err != nil {
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
}

func (h *WebHandlers) musclesPage(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := h.Service.Muscles(ctx, w); err != nil {
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
}
