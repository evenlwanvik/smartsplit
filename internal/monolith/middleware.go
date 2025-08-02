package monolith

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/evenlwanvik/smartsplit/internal/logging"
	"github.com/evenlwanvik/smartsplit/internal/rest"
)

type requestUrlKey string

const RequestUrlKey requestUrlKey = "requestUrlKey"

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestLogger := app.logger.With(
			slog.Group(
				"request",
				slog.String("id", uuid.New().String()),
				slog.String("method", r.Method),
				slog.String("protocol", r.Proto),
				slog.String("url", r.URL.Path),
			),
		)
		ctx = logging.WithLogger(ctx, requestLogger)
		ctx = context.WithValue(ctx, RequestUrlKey, r.URL.Path)

		requestLogger.Info("received request")
		next.ServeHTTP(w, r.WithContext(ctx))
		requestLogger.Info("received completed")
	})
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				rest.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
