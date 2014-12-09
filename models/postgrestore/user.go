package postgrestore

import (
	"git.ianfross.com/expensetracker/auth"
)

func (s *postgresStore) UserByEmail(e string) (*auth.User, error) {
	return nil, nil
}
