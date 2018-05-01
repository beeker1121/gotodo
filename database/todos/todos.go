package todos

import (
	"database/sql"
	"fmt"
	"time"
)

// Database defines the todos database.
type Database struct {
	db *sql.DB
}

// New creates a new todos database.
func New(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

// Todo defines a todo.
type Todo struct {
	ID        int       `json:"id"`
	MemberID  int       `json:"member_id"`
	Created   time.Time `json:"created"`
	Detail    string    `json:"detail"`
	Completed bool      `json:"completed"`
}

// Todos defines a set of todos.
type Todos struct {
	Todos []*Todo `json:"todos"`
	Total int     `json:"total"`
}

const (
	// stmtInsert defines the SQL statement to
	// insert a new todo into the database.
	stmtInsert = `
INSERT INTO todos (member_id, created, detail, completed)
VALUES (?, ?, ?, ?)
`

	// stmtSelect defines the SQL statement to
	// select a set of todos for a given member.
	stmtSelect = `
SELECT id, member_id, created, detail, completed
FROM todos
%s
LIMIT %v, %v
`

	// stmtSelectCount defines the SQL statement to
	// select the total number of todos found for a
	// a given member, according to the filters.
	stmtSelectCount = `
SELECT COUNT(*)
FROM todos
%s
`

	// stmtSelectByID defines the SQL statement to
	// select a todo by its ID.
	stmtSelectByID = `
SELECT id, member_id, created, detail, completed
FROM todos
WHERE id=?
`

	// stmtSelectByIDAndMemberID defines the SQL statement
	// to select a todo by its ID and member ID.
	stmtSelectByIDAndMemberID = `
SELECT id, member_id, created, detail, completed
FROM todos
WHERE id=? AND member_id=?
`

	// stmtUpdate defines the SQL statement to
	// update a todo.
	stmtUpdate = `
UPDATE todos
SET %s
WHERE id=?
`
)

// NewParams defines the parameters for the New method.
type NewParams struct {
	Detail string `json:"detail"`
}

// New creates a new todo.
func (db *Database) New(mid int, params *NewParams) (*Todo, error) {
	// Create a new Todo.
	todo := &Todo{
		MemberID: mid,
		Created:  time.Now(),
		Detail:   params.Detail,
	}

	// Create variable to hold the result.
	var res sql.Result
	var err error

	// Execute the query.
	if res, err = db.db.Exec(stmtInsert, todo.MemberID, todo.Created, todo.Detail, todo.Completed); err != nil {
		return nil, err
	}

	// Get last insert ID.
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	todo.ID = int(id)

	return todo, nil
}

// GetParams defines the parameters for the Get method.
type GetParams struct {
	ID        *int       `json:"id"`
	MemberID  *int       `json:"member_id"`
	Created   *time.Time `json:"created"`
	Completed *bool      `json:"completed"`
	Offset    int        `json:"offset"`
	Limit     int        `json:"limit"`
}

// Get gets a set of todos.
func (db *Database) Get(params *GetParams) (*Todos, error) {
	// Create variables to hold the query fields
	// being filtered on and their values.
	var queryFields string
	var queryValues []interface{}

	// Handle ID field.
	if params.ID != nil {
		if queryFields == "" {
			queryFields = "WHERE id=?"
		} else {
			queryFields += " AND id=?"
		}

		queryValues = append(queryValues, *params.ID)
	}

	// Handle member ID field.
	if params.MemberID != nil {
		if queryFields == "" {
			queryFields = "WHERE member_id=?"
		} else {
			queryFields += " AND member_id=?"
		}

		queryValues = append(queryValues, *params.MemberID)
	}

	// Handle created field.
	if params.Created != nil {
		if queryFields == "" {
			queryFields = "WHERE created=?"
		} else {
			queryFields += " AND created=?"
		}

		queryValues = append(queryValues, *params.Created)
	}

	// Handle completed field.
	if params.Completed != nil {
		if queryFields == "" {
			queryFields = "WHERE completed=?"
		} else {
			queryFields += " AND completed=?"
		}

		queryValues = append(queryValues, *params.Created)
	}

	// Build the full query.
	query := fmt.Sprintf(stmtSelect, queryFields, params.Offset, params.Limit)

	// Create a new Todos.
	todos := &Todos{
		Todos: []*Todo{},
	}

	// Execute the query.
	rows, err := db.db.Query(query, queryValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through the todo rows.
	for rows.Next() {
		// Create a new Todo.
		todo := &Todo{}

		// Scan row values into todo struct.
		if err := rows.Scan(&todo.ID, &todo.MemberID, &todo.Created, &todo.Detail, &todo.Completed); err != nil {
			return nil, err
		}

		// Add to todos set.
		todos.Todos = append(todos.Todos, todo)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Build the total count query.
	queryCount := fmt.Sprintf(stmtSelectCount, queryFields)

	// Get total count.
	var total int
	if err = db.db.QueryRow(queryCount, queryValues...).Scan(&total); err != nil {
		return nil, err
	}
	todos.Total = total

	return todos, nil
}

// GetByID retrieves a todo by its ID.
func (db *Database) GetByID(id int) (*Todo, error) {
	// Create a new Todo.
	todo := &Todo{}

	// Execute the query.
	err := db.db.QueryRow(stmtSelectByID, id).Scan(&todo.ID, &todo.MemberID, &todo.Created, &todo.Detail, &todo.Completed)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrTodoNotFound
	case err != nil:
		return nil, err
	}

	return todo, nil
}

// GetByIDAndMemberID retrieves a todo by its ID and member ID.
func (db *Database) GetByIDAndMemberID(id, mid int) (*Todo, error) {
	// Create a new Todo.
	todo := &Todo{}

	// Execute the query.
	err := db.db.QueryRow(stmtSelectByIDAndMemberID, id, mid).Scan(&todo.ID, &todo.MemberID, &todo.Created, &todo.Detail, &todo.Completed)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrTodoNotFound
	case err != nil:
		return nil, err
	}

	return todo, nil
}

// UpdateParams defines the parameters for the Update method.
type UpdateParams struct {
	Created   *time.Time `json:"created"`
	Detail    *string    `json:"detail"`
	Completed *bool      `json:"completed"`
}

// Update updates a todo.
func (db *Database) Update(id int, params *UpdateParams) (*Todo, error) {
	// Create variables to hold the query fields
	// being updated and their new values.
	var queryFields string
	var queryValues []interface{}

	// Handle created field.
	if params.Created != nil {
		if queryFields == "" {
			queryFields = "created=?"
		} else {
			queryFields += ", created=?"
		}

		queryValues = append(queryValues, *params.Created)
	}

	// Handle detail field.
	if params.Detail != nil {
		if queryFields == "" {
			queryFields = "detail=?"
		} else {
			queryFields += ", detail=?"
		}

		queryValues = append(queryValues, *params.Detail)
	}

	// Handle completed field.
	if params.Completed != nil {
		if queryFields == "" {
			queryFields = "completed=?"
		} else {
			queryFields += ", completed=?"
		}

		// Handle turning boolean into tinyint.
		if *params.Completed {
			queryValues = append(queryValues, 1)
		} else {
			queryValues = append(queryValues, 0)
		}
	}

	// Check if the query is empty.
	if queryFields == "" {
		return db.GetByID(id)
	}

	// Build the full query.
	query := fmt.Sprintf(stmtUpdate, queryFields)
	queryValues = append(queryValues, id)

	// Execute the query.
	_, err := db.db.Exec(query, queryValues...)
	if err != nil {
		return nil, err
	}

	// Since the GetByID method is straight forward,
	// we can use this method to retrieve the updated
	// todo. Anything more complicated should use the
	// original statement constants.
	return db.GetByID(id)
}
