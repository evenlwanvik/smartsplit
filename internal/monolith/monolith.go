package monolith

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/evenlwanvik/smartsplit/internal/identity"
)

// Monolith is the interface that represents the main application
type Monolith interface {
	DB() *sql.DB
	Logger() *slog.Logger
	Mux() *http.ServeMux
}

type Modules []Module

type Identity interface {
	identity.UserClient
}

type Module interface {
	// Setup sets up the module using the context and resources from the
	// monolith. For example, initializing database models, message queues
	// and so on.
	//
	// Be aware that in case of injecting resources from another module,
	// no method call should be made within the Setup process. This can
	// cause segmentation faults as there is no guarantee that injected modules
	// have completed their own startup process.
	//
	// If resources are required from other modules as part of the startup
	// process, add a PostSetup or Startup method, and set it to run after
	// Setup has completed.
	Setup(ctx context.Context, app Monolith) error
	// PostSetup performs any additional setup tasks, usually in cases where
	// external resources from other modules were required, and a guarantee
	// that the initial setup of any given module is completed to avoid
	// segmentation faults. Note that API calls through HTTP will not work
	// if the Monolith server is not running.
	PostSetup(ctx context.Context) error
	// Shutdown performs any necessary cleanup tasks before application
	// termination. Examples of such tasks could be closing channels,
	// connections created by the module and so on.
	//
	// Typically called with the defer keyword immediately after Setup and
	// PostSetup.
	Shutdown()
}
