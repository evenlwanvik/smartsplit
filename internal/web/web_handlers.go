package web

import (
	"context"
	"net/http"

	"github.com/evenlwanvik/smartsplit/internal/logging"
	"github.com/evenlwanvik/smartsplit/internal/rest"
)

// WebHandler defines HTTP handlers for web pages.
type WebHandler struct {
	Service *WebService
}

// RegisterRoutes hooks up endpoints.
func (h *WebHandler) RegisterRoutes(ctx context.Context, mux *http.ServeMux) {
	logger := logging.LoggerFromContext(ctx)

	routeDefinitions := rest.RouteDefinitionList{
		{
			"GET /index",
			h.getIndexHandler,
		},
	}

	for _, d := range routeDefinitions {
		logger.Info("adding route", "route", d.Path)
		mux.HandleFunc(d.Path, d.Handler)
	}
}

func (h *WebHandler) getIndexHandler(w http.ResponseWriter, r *http.Request) {
	err := h.Service.ShowIndex(w)
	if err != nil {
		rest.InternalServerErrorResponse(w, r, err)
		return
	}
}
