package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/evenlwanvik/smartsplit/internal/db"
	"github.com/evenlwanvik/smartsplit/internal/identity"
	"github.com/evenlwanvik/smartsplit/internal/logging"
)

type Module struct {
	Name    string
	Version string
	DB      *sql.DB
	id      identity.UserHandler
}

func main() {
	ctx := context.Background()
	// Initialize the logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	ctx = logging.WithLogger(ctx, logger)

	dsn := db.PostgresDB{
		Host:     "localhost",
		Port:     5032,
		User:     "smartsplit",
		Password: "smartsplit",
		Database: "smartsplit",
		SSLMode:  "disable",
	}

	logger.Info("Connecting to the database", "dsn", dsn)
	db, err := db.NewDB(dsn)
	if err != nil {
		logger.Error("Failed to connect to the database", "error", err)
		panic(err)
	}
	// test Db connection
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping the database", "error", err)
		panic(err)
	}

	userHandlers := identity.UserHandler{Service: identity.NewUserService(identity.NewUserRepository(db))}

	mux := http.NewServeMux()
	userHandlers.RegisterRoutes(mux)

	logger.Info("Starting server on :5050")
	if err := http.ListenAndServe(":5050", mux); err != nil {
		logger.Error("Failed to start server", "error", err)
		panic(err)
	}
}
