package postgrestore

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"testing"
)

var (
	db *sqlx.DB
	s  *postgresStore
)

func TestSetup(t *testing.T) {
	db = sqlx.MustOpen("postgres", "user=ian dbname=expense_test password=wedge89")
	s = MustCreate(db)
}

func TestSchemaCreate(t *testing.T) {
	s.debug = true
	defer func() { s.debug = false }()
	err := db.Ping()
	if err != nil {
		t.Fatalf("Error pinging DB: %v", err)
		return
	}

	s.MustCreateTypes()
	defer s.MustDropTypes()

	s.MustCreateTables()
	s.MustPrepareStmts()
	s.MustDropTables()

}
