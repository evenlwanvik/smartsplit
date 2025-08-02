package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"

	_ "github.com/mattn/go-sqlite3"

	"github.com/evenlwanvik/smartsplit/cmd/api/identity"
	"github.com/evenlwanvik/smartsplit/internal/config"
	"github.com/evenlwanvik/smartsplit/internal/monolith"
)

func main() {
	err := run()
	if err != nil {
		slog.Error("an error occured while starting the application", "error", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	ctx := context.Background()
	slog.Info("starting smartsplit application")

	slog.Info("initializing logger")
	handler := slog.NewJSONHandler(os.Stdout, nil)
	jsonLogger := slog.New(handler)
	logger := jsonLogger.With(
		slog.Group(
			"instance",
			slog.String("id", uuid.New().String()),
		),
	)
	slog.SetDefault(logger)

	logger.Info("initializing database connection")
	os.MkdirAll("db", os.ModePerm) // ensure db dir in root
	db, err := sql.Open("sqlite3", "db/smartsplit.db")
	if err != nil {
		log.Fatalf("failed to connect to SQLite: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(
				"unable to close database connection",
				slog.String("error", err.Error()),
			)
		}
	}()

	mux := http.NewServeMux()

	logger.Info("loading config")
	cfg, err := config.New()
	if err != nil {
		logger.Error("failed to load config")
		return err
	}

	app := monolith.NewApplication(
		db,
		mux,
		logger,
		cfg,
		monolith.Modules{
			Identity: &identity.Module{},
		},
	)

	logger.Info("running module startup procedures")
	app.SetupModules(ctx)
	app.PostSetupModules()

	err = app.Serve()
	if err != nil {
		logger.Error("unable to start server", "error", err)
		return err
	}

	app.Logger().Info("shutting down modules")
	app.ShutdownModules()

	app.Logger().Info("exiting...")

	return nil
}
