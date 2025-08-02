package monolith

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/justinas/alice"
)

type Application struct {
	db      *sql.DB
	mux     *http.ServeMux
	logger  *slog.Logger
	modules Modules
	done    <-chan os.Signal
}

func NewApplication(
	db *sql.DB, mux *http.ServeMux, logger *slog.Logger, modules Modules,
) *Application {
	return &Application{
		db:      db,
		mux:     mux,
		logger:  logger,
		modules: modules,
	}
}

func (app *Application) DB() *sql.DB          { return app.db }
func (app *Application) Logger() *slog.Logger { return app.logger }
func (app *Application) Mux() *http.ServeMux  { return app.mux }

func (app *Application) SetupModules(ctx context.Context) error {
	for _, module := range app.modules {
		if err := module.Setup(ctx, app); err != nil {
			return err
		}
	}
	return nil
}

func (app *Application) PostSetupModules(ctx context.Context) error {
	for _, module := range app.modules {
		if err := module.PostSetup(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (app *Application) ShutdownModules() {
	for i := len(app.modules) - 1; i >= 0; i-- {
		app.modules[i].Shutdown()
	}
}

func (app *Application) routes() http.Handler {
	app.logger.Info("creating standard middleware chain")
	standard := alice.New(
		app.recoverPanic,
		app.logRequest,
	)

	// healthcheck
	app.logger.Info("adding healthcheck route")
	app.mux.HandleFunc("GET /api/v1/healthcheck", app.healthcheckHandler)

	// profiling
	app.mux.HandleFunc("GET /debug/pprof/", http.DefaultServeMux.ServeHTTP)
	app.mux.HandleFunc("GET /debug/pprof/profile", http.DefaultServeMux.ServeHTTP)
	app.mux.HandleFunc("GET /debug/pprof/heap", http.DefaultServeMux.ServeHTTP)

	handler := standard.Then(app.mux)
	return handler
}

func (app *Application) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.App.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		slog.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.App.Env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", "addr", srv.Addr)

	return nil
}
