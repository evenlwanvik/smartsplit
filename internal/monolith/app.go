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
	"reflect"
	"syscall"
	"time"

	"github.com/justinas/alice"

	"github.com/evenlwanvik/smartsplit/internal/config"
)

type Application struct {
	db      *sql.DB
	mux     *http.ServeMux
	config  *config.Config
	logger  *slog.Logger
	modules Modules
	done    <-chan os.Signal
}

func NewApplication(
	db *sql.DB,
	mux *http.ServeMux,
	logger *slog.Logger,
	config *config.Config,
	modules Modules,
) *Application {
	return &Application{
		db:      db,
		mux:     mux,
		logger:  logger,
		config:  config,
		modules: modules,
	}
}

func (app *Application) DB() *sql.DB            { return app.db }
func (app *Application) Logger() *slog.Logger   { return app.logger }
func (app *Application) Mux() *http.ServeMux    { return app.mux }
func (app *Application) Config() *config.Config { return app.config }
func (app *Application) Modules() *Modules {
	return &app.modules
}

func (app *Application) SetupModules(ctx context.Context) {
	app.logger.Info("running setupModules")
	val := reflect.ValueOf(app.modules)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if module, ok := field.Interface().(Module); ok {
			module.Setup(ctx, app)
		}
	}
}

func (app *Application) PostSetupModules() {
	app.logger.Info("running postSetupModules")
	val := reflect.ValueOf(app.modules)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if module, ok := field.Interface().(Module); ok {
			module.PostSetup()
		}
	}
}

func (app *Application) ShutdownModules() {
	app.logger.Info("running shutdownModules")
	val := reflect.ValueOf(app.modules)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if module, ok := field.Interface().(Module); ok {
			module.Shutdown()
		}
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
