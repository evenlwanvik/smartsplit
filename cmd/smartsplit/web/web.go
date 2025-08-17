package web

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/evenlwanvik/smartsplit/internal/monolith"
	"github.com/evenlwanvik/smartsplit/internal/web"
	"github.com/evenlwanvik/smartsplit/internal/workout"
)

const moduleName string = "web"

type Module struct {
	logger  *slog.Logger
	name    string
	version string
	mux     *http.ServeMux
	web     web.Service
	workout workout.Client
}

func (m *Module) Setup(ctx context.Context, mono monolith.Monolith) {
	m.initModuleLogger(mono.Logger())

	m.web = web.NewService(mono.Modules().Workout)

	// TODO: We have to wait for the monolith to be fully initialized before we can inject modules
	m.logger.Info("injecting mux")
	m.mux = mono.Mux()

	m.logger.Info("registering routes")
	m.web.RegisterRoutes(ctx, m.mux)
}

func (m *Module) PostSetup() {
	m.logger.Info("performing post setup process")
}

func (m *Module) Shutdown() {}

func (m *Module) initModuleLogger(monoLogger *slog.Logger) {
	m.logger = monoLogger.With(slog.Group("module", slog.String("name", moduleName)))
}
