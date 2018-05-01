package todos

import (
	"gotodo/database"
	dbtodos "gotodo/database/todos"
	"gotodo/services/errors"
)

// Service defines the todos service.
type Service struct {
	db *database.Database
}

// New returns a new todos service.
func New(db *database.Database) *Service {
	return &Service{
		db: db,
	}
}

// Todo defines a todo.
type Todo dbtodos.Todo

// Todos defines a set of todos.
type Todos struct {
	Todos []*Todo `json:"todos"`
	Total int     `json:"total"`
}

// NewParams defines the parameters for the New method.
type NewParams dbtodos.NewParams

// New creates a new todo.
func (s *Service) New(mid int, params *NewParams) (*Todo, error) {
	// Create a new ParamErrors.
	pes := errors.NewParamErrors()

	// Check detail.
	if params.Detail == "" {
		pes.Add(errors.NewParamError("detail", ErrDetailEmpty))
	}

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return nil, pes
	}

	// Create this member in the database.
	dbt, err := s.db.Todos.New(mid, &dbtodos.NewParams{
		Detail: params.Detail,
	})
	if err != nil {
		return nil, err
	}

	// Create a new Todo.
	todo := &Todo{
		ID:        dbt.ID,
		MemberID:  dbt.MemberID,
		Created:   dbt.Created,
		Detail:    dbt.Detail,
		Completed: dbt.Completed,
	}

	return todo, nil
}

// GetParams defines the parameters for the Get method.
type GetParams dbtodos.GetParams

// Get gets a set of todos.
func (s *Service) Get(params *GetParams) (*Todos, error) {
	// Try to pull the todos from the database.
	dbts, err := s.db.Todos.Get(&dbtodos.GetParams{
		ID:        params.ID,
		MemberID:  params.MemberID,
		Created:   params.Created,
		Completed: params.Completed,
		Offset:    params.Offset,
		Limit:     params.Limit,
	})
	if err != nil {
		return nil, err
	}

	// Create a new Todos.
	todos := &Todos{
		Todos: []*Todo{},
		Total: dbts.Total,
	}

	// Loop through the set of todos.
	for _, t := range dbts.Todos {
		// Create a new Todo.
		todo := &Todo{
			ID:        t.ID,
			MemberID:  t.MemberID,
			Created:   t.Created,
			Detail:    t.Detail,
			Completed: t.Completed,
		}

		// Add to todos set.
		todos.Todos = append(todos.Todos, todo)
	}

	return todos, nil
}

// GetByIDAndMemberID retrieves a todo by its ID and member ID.
func (s *Service) GetByIDAndMemberID(id, mid int) (*Todo, error) {
	// Try to pull this todo from the database.
	dbt, err := s.db.Todos.GetByIDAndMemberID(id, mid)
	if err != nil {
		return nil, err
	}

	// Create a new Todo.
	todo := &Todo{
		ID:        dbt.ID,
		MemberID:  dbt.MemberID,
		Created:   dbt.Created,
		Detail:    dbt.Detail,
		Completed: dbt.Completed,
	}

	return todo, nil
}

// UpdateParams defines the parameters for the update methods.
type UpdateParams dbtodos.UpdateParams

// UpdateByIDAndMemberID updates a todo.
func (s *Service) UpdateByIDAndMemberID(id, mid int, params *UpdateParams) (*Todo, error) {
	// Try to pull this todo from the database.
	dbt, err := s.db.Todos.GetByIDAndMemberID(id, mid)
	if err == dbtodos.ErrTodoNotFound {
		return nil, ErrTodoNotFound
	} else if err != nil {
		return nil, err
	}

	// Update this todo in the database.
	dbt, err = s.db.Todos.Update(id, &dbtodos.UpdateParams{
		Created:   params.Created,
		Detail:    params.Detail,
		Completed: params.Completed,
	})
	if err != nil {
		return nil, err
	}

	// Create a new Todo.
	todo := &Todo{
		ID:        dbt.ID,
		MemberID:  dbt.MemberID,
		Created:   dbt.Created,
		Detail:    dbt.Detail,
		Completed: dbt.Completed,
	}

	return todo, nil
}
