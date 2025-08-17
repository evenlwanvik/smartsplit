package auth

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/evenlwanvik/smartsplit/internal/auth"
	"github.com/evenlwanvik/smartsplit/internal/monolith"
)

const moduleName string = "auth"

type Module struct {
	logger   *slog.Logger
	name     string
	version  string
	db       *sql.DB
	mux      *http.ServeMux
	handlers auth.UserHandler
}

func (m *Module) Setup(ctx context.Context, mono monolith.Monolith) {
	m.initModuleLogger(mono.Logger())

	m.logger.Info("injecting database connection pool")
	m.db = mono.DB()

	m.handlers = auth.UserHandler{
		Service: auth.NewUserService(
			auth.NewUserRepository(m.db),
		),
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
