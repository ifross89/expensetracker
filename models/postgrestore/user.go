package postgrestore

import (
	"git.ianfross.com/ifross/expensetracker/auth"

	"github.com/juju/errors"

	"fmt"
)

const (
	insertUserStr = `
INSERT INTO users (name, email, pw_hash, admin, active, token)
    VALUES(:name, :email, :pw_hash, :admin, :active, :token) RETURNING *;`

	userByEmailStr = "SELECT * FROM users WHERE email=:email;"

	userByIDStr    = "SELECT * FROM users WHERE id=:id;"
	userByTokenStr = "SELECT * FROM users WHERE token=:token;"
	updateUserStr  = `
UPDATE users SET
		name=:name,
		email=:email,
		pw_hash=:pw_hash,
		admin=:admin,
		active=:active,
		token=:token
	WHERE id=:id;`
	deleteUserStr = "DELETE FROM users WHERE id=:id;"
	usersStr      = `SELECT * FROM users;`
)

func (s *postgresStore) Users() ([]*auth.User, error) {
	var us []*auth.User = make([]*auth.User, 0, 0)
	err := s.db.Select(&us, usersStr)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return us, nil
}

// Insert saves a new user to the database
func (s *postgresStore) Insert(u *auth.User) error {
	fmt.Println("USER:", u)
	if u.ID != 0 {
		return auth.ErrAlreadySaved
	}
	err := s.insertUserStmt.Get(u, u)
	if err != nil {
		return errors.Annotate(err, "Error inserting user")
	}
	return nil
}

// Update updated a user in the database
func (s *postgresStore) Update(u *auth.User) error {
	_, err := s.updateUserStmt.Exec(u)
	if err != nil {
		return errors.Annotate(err, "Error updating user")
	}

	return nil
}

// UserByToken retrieves a user by their unique token
func (s *postgresStore) UserByToken(tok string) (*auth.User, error) {
	var u = auth.User{Token: tok}
	err := s.userByTokenStmt.Get(&u, u)
	if err != nil {
		return nil, errors.Annotatef(err, "Could not find user with token %s", tok)
	}

	return &u, nil
}

// UserByID retrieves a user by their ID
func (s *postgresStore) UserByID(id int64) (*auth.User, error) {
	var u = auth.User{ID: id}
	err := s.userByIDStmt.Get(&u, u)
	if err != nil {
		return nil, errors.Annotatef(err, "Could not find user with id %d", id)
	}

	return &u, nil
}

// Delete removes a user from the database. If the user does not exist, then
// an error is returned
func (s *postgresStore) Delete(u *auth.User) error {
	result, err := s.deleteUserStmt.Exec(u)
	if err != nil {
		return errors.Annotatef(err, "Could not delete user with id %d", u.ID)
	}
	n, _ := result.RowsAffected()
	fmt.Println("Rows affected", n)
	if n != 1 {
		return errors.New("No user deleted")
	}
	u.ID = 0
	return nil
}

// UserByEmail obtains a user by their email address
func (s *postgresStore) UserByEmail(e string) (*auth.User, error) {
	var u = auth.User{Email: e}
	err := s.userByEmailStmt.Get(&u, u)
	if err != nil {
		return nil, errors.Annotatef(err, "Could not find user %s", e)
	}

	return &u, nil
}
