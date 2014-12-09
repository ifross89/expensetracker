package postgrestore

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"testing"
)

func TestSchemaCreate(t *testing.T) {
	db := sqlx.MustOpen("postgres", "user=postgres dbname=expense_test password=wedge89")
	err := db.Ping()
	if err != nil {
		t.Fatalf("Error pinging DB: %v", err)
	}

	s := Create(db)
	s.debug = true
	s.MustCreateTypes()
	defer s.MustDropTypes()

	s.MustCreateTables()
	s.MustDropTables()
}
