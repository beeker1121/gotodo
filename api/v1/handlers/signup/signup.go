package signup

import (
	"encoding/json"
	"net/http"

	apictx "gotodo/api/context"
	"gotodo/api/errors"
	"gotodo/api/middleware/auth"
	"gotodo/api/render"
	serverrors "gotodo/services/errors"
	"gotodo/services/members"

	"github.com/beeker1121/httprouter"
)

// ResultPost defines the response data for the HandlePost handler.
type ResultPost struct {
	Data string `json:"data"`
}

// New creates the routes for the signup endpoints of the API.
func New(ac *apictx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.POST("/api/v1/signup", HandlePost(ac))
}

// HandlePost handles the /api/v1/signup POST route of the API.
func HandlePost(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the parameters from the request body.
		var params members.NewParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			errors.Default(ac.Logger, w, errors.ErrBadRequest)
			return
		}

		// Create the member.
		member, err := ac.Services.Members.New(&params)
		if pes, ok := err.(*serverrors.ParamErrors); ok && err != nil {
			errors.Params(ac.Logger, w, http.StatusBadRequest, pes)
			return
		} else if err != nil {
			ac.Logger.Printf("members.New() service error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Issue a new JWT for this member.
		token, err := auth.NewJWT(ac, member.Password, member.ID)
		if err != nil {
			ac.Logger.Printf("auth.NewJWT() error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultPost{
			Data: token,
		}

		// Render output.
		if err := render.JSON(w, true, result); err != nil {
			ac.Logger.Printf("render.JSON() error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}
	}
}
