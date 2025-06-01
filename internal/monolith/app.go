package monolith

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
)

type MonolithApp struct {
	Name   string
	db     *sql.DB
	mux    http.ServeMux
	logger slog.Logger
	done   <-chan os.Signal
}
