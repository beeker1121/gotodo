package members

import (
	"gotodo/database"
	dbmembers "gotodo/database/members"
	"gotodo/services/errors"

	"golang.org/x/crypto/bcrypt"
)

// Service defines the members service.
type Service struct {
	db *database.Database
}

// New returns a new members service.
func New(db *database.Database) *Service {
	return &Service{
		db: db,
	}
}

// Members defines a member.
type Member dbmembers.Member

// NewParams defines the parameters for the New method.
type NewParams dbmembers.NewParams

// New creates a new member.
func (s *Service) New(params *NewParams) (*Member, error) {
	// Create a new ParamErrors.
	pes := errors.NewParamErrors()

	// Check email.
	if params.Email == "" {
		pes.Add(errors.NewParamError("email", ErrEmailEmpty))
	} else {
		_, err := s.db.Members.GetByEmail(params.Email)
		if err == nil {
			pes.Add(errors.NewParamError("email", ErrEmailExists))
		} else if err != nil && err != dbmembers.ErrMemberNotFound {
			return nil, err
		}
	}

	// Check password.
	if len(params.Password) < 8 {
		pes.Add(errors.NewParamError("password", ErrPassword))
	}

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return nil, pes
	}

	// Hash the password.
	pwhash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create this member in the database.
	dbm, err := s.db.Members.New(&dbmembers.NewParams{
		Email:    params.Email,
		Password: string(pwhash),
	})
	if err != nil {
		return nil, err
	}

	// Create a new Member.
	member := &Member{
		ID:       dbm.ID,
		Email:    dbm.Email,
		Password: dbm.Password,
	}

	return member, nil
}

// LoginParams defines the parameters for the Login method.
type LoginParams struct {
	Email    string
	Password string
}

// Login checks if a member exists in the database and can log in.
func (s *Service) Login(params *LoginParams) (*Member, error) {
	// Try to pull this member from the database.
	dbm, err := s.db.Members.GetByEmail(params.Email)
	if err == dbmembers.ErrMemberNotFound {
		return nil, ErrInvalidLogin
	} else if err != nil {
		return nil, err
	}

	// Validate the password.
	if err = bcrypt.CompareHashAndPassword([]byte(dbm.Password), []byte(params.Password)); err != nil {
		return nil, ErrInvalidLogin
	}

	// Create a new Member.
	member := &Member{
		ID:       dbm.ID,
		Email:    dbm.Email,
		Password: dbm.Password,
	}

	return member, nil
}

// GetByID retrieves a member by their ID.
func (s *Service) GetByID(id int) (*Member, error) {
	// Try to pull this member from the database.
	dbm, err := s.db.Members.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Create a new Member.
	member := &Member{
		ID:       dbm.ID,
		Email:    dbm.Email,
		Password: dbm.Password,
	}

	return member, nil
}
