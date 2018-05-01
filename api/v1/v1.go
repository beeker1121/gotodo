package v1

import (
	apictx "gotodo/api/context"
	"gotodo/api/v1/handlers/login"
	"gotodo/api/v1/handlers/signup"
	"gotodo/api/v1/handlers/todos"

	"github.com/beeker1121/httprouter"
)

// New creates a new API v1 application. All of the necessary routes for
// v1 of the API will be created on the given router, which should then be
// used to create the web server. The root domain should be api.maildb.io
// or something similar.
func New(ac *apictx.Context, router *httprouter.Router) {
	// Create all of the API v1 routes.
	signup.New(ac, router)
	login.New(ac, router)
	todos.New(ac, router)
}
