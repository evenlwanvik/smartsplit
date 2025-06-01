package identity

import (
	"database/sql"

	"github.com/evenlwanvik/smartsplit/internal/identity"
)

type Module struct {
	Name    string
	Version string
	DB      *sql.DB
	id      identity.UserHandler
}

func (m *Module) Init() error {
	var err error
	m.id, err = identity.NewUserHandler(m.DB)
	if err != nil {
		return err
	}
	return nil
}
