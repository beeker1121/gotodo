package context

import (
	"log"

	"gotodo/api/config"
	"gotodo/services"
)

// Context defines the API context, which acts as a container for all assets
// used by the API.
type Context struct {
	Config   *config.Config
	Logger   *log.Logger
	Services *services.Services
}

// New returns a new API context.
func New(config *config.Config, logger *log.Logger, services *services.Services) *Context {
	return &Context{
		Config:   config,
		Logger:   logger,
		Services: services,
	}
}
