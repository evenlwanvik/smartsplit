package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/evenlwanvik/smartsplit/internal/common"
	"github.com/evenlwanvik/smartsplit/internal/identity"
)

func main() {
	ctx := context.Background()
	// Initialize the logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	ctx = common.WithLogger(ctx, logger)

	// postgresql dsn
	dsn := "host=localhost port=5432 user=postgres password=yourpassword dbname=yourdb sslmode=disable"

	logger.Info("Connecting to the database", "dsn", dsn)
	db, err := common.NewDB(dsn)
	if err != nil {
		logger.Error("Failed to connect to the database", "error", err)
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
