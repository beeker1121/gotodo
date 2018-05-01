package todos

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	apictx "gotodo/api/context"
	"gotodo/api/errors"
	"gotodo/api/middleware/auth"
	"gotodo/api/render"
	serverrors "gotodo/services/errors"
	servtodos "gotodo/services/todos"

	"github.com/beeker1121/httprouter"
)

// Todo defines the todo API type.
//
// This mirrors the service Todo type, which mirrors the database Todo type.
// However, we specify that the MemberID should not be included when encoding
// to JSON.
type Todo struct {
	ID        int       `json:"id"`
	MemberID  int       `json:"-"`
	Created   time.Time `json:"created"`
	Detail    string    `json:"detail"`
	Completed bool      `json:"completed"`
}

// Meta defines the response top level meta object.
type Meta struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

// Links defines the response top level links object.
type Links struct {
	Prev *string `json:"prev"`
	Next *string `json:"next"`
}

// ResultGet defines the response data for the HandleGet handler.
type ResultGet struct {
	Data  []*Todo `json:"data"`
	Meta  Meta    `json:"meta"`
	Links Links   `json:"links"`
}

// ResultGetTodo defines the response data for the HandleGetTodo handler.
type ResultGetTodo struct {
	Data *Todo `json:"data"`
}

// ResultPost defines the response data for the HandlePost handler.
type ResultPost struct {
	Data *Todo `json:"data"`
}

// ResultUpdate defines the response data for the HandleUpdate handler.
type ResultUpdate struct {
	Data *Todo `json:"data"`
}

// New creates the routes for the todo endpoints of the API.
func New(ac *apictx.Context, router *httprouter.Router) {
	// Handle the routes.
	router.GET("/api/v1/todos", auth.AuthenticateEndpoint(ac, HandleGet(ac)))
	router.GET("/api/v1/todos/:id", auth.AuthenticateEndpoint(ac, HandleGetTodo(ac)))
	router.POST("/api/v1/todos", auth.AuthenticateEndpoint(ac, HandlePost(ac)))
	router.POST("/api/v1/todos/:id", auth.AuthenticateEndpoint(ac, HandleUpdate(ac)))
}

// HandleGet handles the /api/v1/todos GET route of the API.
func HandleGet(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get this member from the request context.
		member, err := auth.GetMemberFromRequest(r)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new GetParams.
		params := &servtodos.GetParams{
			MemberID: &member.ID,
		}

		// Create a new API Errors.
		errs := &errors.Errors{}

		// Handle created.
		if createdqs, ok := r.URL.Query()["created"]; ok && len(createdqs) == 1 {
			t, err := time.Parse(time.RFC3339, createdqs[0])
			if err != nil {
				errs.Add(errors.New(http.StatusBadRequest, "created", ErrCreatedInvalid.Error()))
			} else {
				params.Created = &t
			}
		}

		// Handle offset.
		if offsetqs, ok := r.URL.Query()["offset"]; ok && len(offsetqs) == 1 {
			offset64, err := strconv.ParseInt(offsetqs[0], 10, 32)
			if err != nil {
				errs.Add(errors.New(http.StatusBadRequest, "offset", ErrLimitInvalid.Error()))
			} else {
				params.Offset = int(offset64)
			}
		} else {
			params.Offset = 0
		}

		// Handle limit.
		if limitqs, ok := r.URL.Query()["limit"]; ok && len(limitqs) == 1 {
			limit64, err := strconv.ParseInt(limitqs[0], 10, 32)
			if err != nil {
				errs.Add(errors.New(http.StatusBadRequest, "limit", ErrLimitInvalid.Error()))
			} else {
				if int(limit64) > ac.Config.LimitMax {
					errs.Add(errors.New(http.StatusBadRequest, "limit", ErrLimitMax.Error()+" of "+strconv.FormatUint(uint64(ac.Config.LimitMax), 10)))
				} else {
					params.Limit = int(limit64)
				}
			}
		} else {
			params.Limit = ac.Config.LimitDefault
		}

		// Return if there were errors.
		if errs.Length() > 0 {
			errors.Multiple(ac.Logger, w, http.StatusBadRequest, errs)
			return
		}

		// Try to get the todos.
		todos, err := ac.Services.Todos.Get(params)
		if err != nil {
			ac.Logger.Printf("todos.Get() service error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultGet{
			Data: []*Todo{},
			Meta: Meta{
				Offset: params.Offset,
				Limit:  params.Limit,
				Total:  todos.Total,
			},
			Links: Links{},
		}

		// Loop through the todos.
		for _, t := range todos.Todos {
			// Copy the Todo type over.
			todo := &Todo{
				ID:        t.ID,
				MemberID:  t.MemberID,
				Created:   t.Created,
				Detail:    t.Detail,
				Completed: t.Completed,
			}

			result.Data = append(result.Data, todo)
		}

		// Handle previous link.
		if params.Offset > 0 {
			limitstr := "&limit=" + strconv.FormatInt(int64(params.Limit), 10)

			offsetstr := "?offset="
			if params.Offset-params.Limit < 0 {
				offsetstr += "0"
			} else {
				offsetstr += strconv.FormatInt(int64(params.Offset-params.Limit), 10)
			}

			prev := "https://" + ac.Config.APIHost + "/api/v1/todos" + offsetstr + limitstr
			result.Links.Prev = &prev
		}

		// Handle next link.
		if params.Offset+params.Limit < todos.Total {
			offsetstr := "?offset=" + strconv.FormatInt(int64(params.Offset+params.Limit), 10)
			limitstr := "&limit=" + strconv.FormatInt(int64(params.Limit), 10)

			next := "https://" + ac.Config.APIHost + "/api/v1/todos" + offsetstr + limitstr
			result.Links.Next = &next
		}

		// Render output.
		if err := render.JSON(w, true, result); err != nil {
			ac.Logger.Printf("render.JSON() error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}
	}
}

// HandleGetTodo handles the /api/v1/todos/:id GET route of the API.
func HandleGetTodo(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to get the todo ID.
		var id int
		id64, err := strconv.ParseInt(httprouter.GetParam(r, "id"), 10, 32)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrBadRequest)
			return
		}
		id = int(id64)

		// Get this member from the request context.
		member, err := auth.GetMemberFromRequest(r)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Try to get this todo.
		todo, err := ac.Services.Todos.GetByIDAndMemberID(id, member.ID)
		if err == servtodos.ErrTodoNotFound {
			errors.Default(ac.Logger, w, errors.New(http.StatusNotFound, "", err.Error()))
			return
		} else if err != nil {
			ac.Logger.Printf("todos.GetByIDAndMemberID() service error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultGetTodo{
			Data: &Todo{
				ID:        todo.ID,
				MemberID:  todo.MemberID,
				Created:   todo.Created,
				Detail:    todo.Detail,
				Completed: todo.Completed,
			},
		}

		// Render output.
		if err := render.JSON(w, true, result); err != nil {
			ac.Logger.Printf("render.JSON() error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}
	}
}

// HandlePost handles the /api/v1/todos POST route of the API.
func HandlePost(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the parameters from the request body.
		var params servtodos.NewParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			errors.Default(ac.Logger, w, errors.ErrBadRequest)
			return
		}

		// Get this member from the request context.
		member, err := auth.GetMemberFromRequest(r)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Try to create a new todo.
		todo, err := ac.Services.Todos.New(member.ID, &params)
		if pes, ok := err.(*serverrors.ParamErrors); ok && err != nil {
			errors.Params(ac.Logger, w, http.StatusBadRequest, pes)
			return
		} else if err != nil {
			ac.Logger.Printf("todos.New() service error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultPost{
			Data: &Todo{
				ID:        todo.ID,
				MemberID:  todo.MemberID,
				Created:   todo.Created,
				Detail:    todo.Detail,
				Completed: todo.Completed,
			},
		}

		// Render output.
		if err := render.JSON(w, true, result); err != nil {
			ac.Logger.Printf("render.JSON() error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}
	}
}

// HandleUpdate handles the /api/v1/todos/:id POST route of the API.
func HandleUpdate(ac *apictx.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the parameters from the request body.
		var params servtodos.UpdateParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			errors.Default(ac.Logger, w, errors.ErrBadRequest)
			return
		}

		// Try to get the todo ID.
		var id int
		id64, err := strconv.ParseInt(httprouter.GetParam(r, "id"), 10, 32)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrBadRequest)
			return
		}
		id = int(id64)

		// Get this member from the request context.
		member, err := auth.GetMemberFromRequest(r)
		if err != nil {
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Try to update this todo.
		todo, err := ac.Services.Todos.UpdateByIDAndMemberID(id, member.ID, &params)
		if pes, ok := err.(*serverrors.ParamErrors); ok && err != nil {
			errors.Params(ac.Logger, w, http.StatusBadRequest, pes)
			return
		} else if err == servtodos.ErrTodoNotFound {
			errors.Default(ac.Logger, w, errors.New(http.StatusNotFound, "", err.Error()))
			return
		} else if err != nil {
			ac.Logger.Printf("todos.New() service error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}

		// Create a new Result.
		result := ResultUpdate{
			Data: &Todo{
				ID:        todo.ID,
				MemberID:  todo.MemberID,
				Created:   todo.Created,
				Detail:    todo.Detail,
				Completed: todo.Completed,
			},
		}

		// Render output.
		if err := render.JSON(w, true, result); err != nil {
			ac.Logger.Printf("render.JSON() error: %s\n", err)
			errors.Default(ac.Logger, w, errors.ErrInternalServerError)
			return
		}
	}
}
