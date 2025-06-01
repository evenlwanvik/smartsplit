package monolith

import (
	"github.com/evenlwanvik/smartsplit/internal/identity"
)

type Module struct {
	identity Identity
}

type Identity interface {
	identity.UserClient
}
