package monolith

import (
	"net/http"

	"github.com/evenlwanvik/smartsplit/internal/config"
	"github.com/evenlwanvik/smartsplit/internal/rest"
)

type HealthCheckMessage struct {
	Status      string             `json:"status"`
	Environment config.Environment `json:"environment"`
}

// @Summary Healthcheck
// @Description Endpoint to check if the API is running
// @Tags    Healthcheck
// @Produce json
// @Success 200 {object} HealthCheckMessage "OK"
// @Failure 500 {object} rest.ErrorMessage "Internal Server Error"
// @Router /api/v1/healthcheck [get]
func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	healthCheckMessage := HealthCheckMessage{
		Status:      "available",
		Environment: app.config.App.Env,
	}

	rest.RespondWithJSON(w, r, http.StatusOK, healthCheckMessage, nil)
}
