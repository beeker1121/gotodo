package members

import "database/sql"

// Database defines the members database.
type Database struct {
	db *sql.DB
}

// New creates a new members database.
func New(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

// Member defines a member.
type Member struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

const (
	// stmtInsert defines the SQL statement to
	// insert a new member into the database.
	stmtInsert = `
INSERT INTO members (email, password)
VALUES (?, ?)
`

	// stmtSelectByID defines the SQL statement to
	// select a member by their ID.
	stmtSelectByID = `
SELECT id, email, password
FROM members
WHERE id=?
`

	// stmtSelectByEmail defines the SQL statement
	// to select a member by their email address.
	stmtSelectByEmail = `
SELECT id, email, password
FROM members
WHERE email=?
`
)

// NewParams defines the parameters for the New method.
type NewParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// New creates a new member.
func (db *Database) New(params *NewParams) (*Member, error) {
	// Create a new Member.
	member := &Member{
		Email:    params.Email,
		Password: params.Password,
	}

	// Create variable to hold the result.
	var res sql.Result
	var err error

	// Execute the query.
	if res, err = db.db.Exec(stmtInsert, member.Email, member.Password); err != nil {
		return nil, err
	}

	// Get last insert ID.
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	member.ID = int(id)

	return member, nil
}

// GetByID retrieves a member by their ID.
func (db *Database) GetByID(id int) (*Member, error) {
	// Create a new Member.
	member := &Member{}

	// Execute the query.
	err := db.db.QueryRow(stmtSelectByID, id).Scan(&member.ID, &member.Email, &member.Password)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrMemberNotFound
	case err != nil:
		return nil, err
	}

	return member, nil
}

// GetByEmail retrieves a member by their email.
func (db *Database) GetByEmail(email string) (*Member, error) {
	// Create a new Member.
	member := &Member{}

	// Execute the query.
	err := db.db.QueryRow(stmtSelectByEmail, email).Scan(&member.ID, &member.Email, &member.Password)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrMemberNotFound
	case err != nil:
		return nil, err
	}

	return member, nil
}
