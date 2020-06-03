package resolver

import (
	"github.com/my/app/server/config"
)

// Resolver -
//
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
type Resolver struct {
	*config.Server
}
