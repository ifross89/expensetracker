package postgrestore

import (
	"git.ianfross.com/expensetracker/auth"

	"github.com/juju/errors"
)

const (
	insertUserStr = `
INSERT INTO users (email, pw_hash, admin, active, token)
    VALUES(:email, :pw_hash, :admin, :active, :token) RETURNING *;`

	userByEmailStr = "SELECT * FROM users WHERE email=:email;"

	userByIdStr    = "SELECT * FROM users WHERE id=:id;"
	userByTokenStr = "SELECT * FROM USERS WHERE token=:token;"
)

func (s *postgresStore) Insert(u *auth.User) error {
	err := s.insertUserStmt.Get(u, u)
	if err != nil {
		return errors.Annotate(err, "Error inserting user")
	}
	return nil
}

func (s *postgresStore) UserByEmail(e string) (*auth.User, error) {
	var u = auth.User{Email: e}
	err := s.userByEmailStmt.Get(&u, u)
	if err != nil {
		return nil, errors.Annotatef(err, "Could not find user %s", e)
	}

	return &u, nil
}
