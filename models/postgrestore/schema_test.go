package postgrestore

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

var (
	db *sqlx.DB
	s  *postgresStore
)

func init() {
	db = sqlx.MustOpen("postgres", "user=ian dbname=expense_test password=wedge89")
	s = MustCreate(db)
}

func TestSchemaCreate(t *testing.T) {
	s.debug = false
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

func TestMustPrepare(t *testing.T) {
	Convey("Attempt to prepare an invalid SQL string", t, func() {
		So(func() { s.mustPrepareStmt("INVALID SQL") }, ShouldPanic)
	})
}
