package postgrestore

import (
	"git.ianfross.com/expensetracker/auth"
)

const (
	insertUserStr = `
INSERT INTO users (email, pw_hash, admin, active, token)
    VALUES(:email, :pw_hash, :admin, :active, :token) RETURNING *;`

	userByEmailStr = "SELECT * FROM users WHERE email=:email;"
)

func (s *postgresStore) Insert(u *auth.User) error {
	nstmt, err := s.db.PrepareNamed(insertUserStr)
	if err != nil {
		return err
	}

	err = nstmt.Get(u, u)
	if err != nil {
		return err
	}
	return nil
}

func (s *postgresStore) UserByEmail(e string) (*auth.User, error) {
	return nil, nil
}
