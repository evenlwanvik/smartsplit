package workout

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/evenlwanvik/smartsplit/internal/monolith"
	"github.com/evenlwanvik/smartsplit/internal/workout"
)

const moduleName string = "workout"

type Module struct {
	logger   *slog.Logger
	name     string
	version  string
	db       *sql.DB
	mux      *http.ServeMux
	handlers workout.Handlers
	svc      *workout.Service
}

func (m *Module) Setup(ctx context.Context, mono monolith.Monolith) {
	m.initModuleLogger(mono.Logger())

	m.logger.Info("injecting database connection pool")
	m.db = mono.DB()

	m.svc = workout.NewService(workout.NewRepository(m.db))

	m.handlers = workout.Handlers{
		Svc: m.svc,
	}

	m.logger.Info("injecting mux")
	m.mux = mono.Mux()

	m.logger.Info("registering routes")
	m.handlers.RegisterRoutes(ctx, m.mux)
}

func (m *Module) PostSetup() {
	m.logger.Info("performing post setup process")
}

func (m *Module) Shutdown() {}

func (m *Module) initModuleLogger(monoLogger *slog.Logger) {
	m.logger = monoLogger.With(slog.Group("module", slog.String("name", moduleName)))
}
