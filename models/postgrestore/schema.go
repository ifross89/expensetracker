package postgrestore

import (
	"github.com/jmoiron/sqlx"

	"fmt"
)

const (
	setTimeZoneStr      = "SET TIME ZONE 'UTC';"
	createCategoriesStr = `
CREATE TYPE category_t as ENUM(
	'groceries',
	'alcohol',
	'drugs',
	'household items',
	'bills',
	'presents',
	'tickets'
);`

	dropCategoriesStr = "DROP TYPE category_t;"

	createUsersTableStr = `
CREATE TABLE IF NOT EXISTS users (
	id                  SERIAL PRIMARY KEY,
	email               VARCHAR(64) NOT NULL CHECK (email <> '') UNIQUE,
	pw_hash             VARCHAR(128),
	admin               BOOLEAN NOT NULL DEFAULT false,
	active              BOOLEAN NOT NULL DEFAULT false,
	token               TEXT,
	created_at          TIMESTAMP DEFAULT LOCALTIMESTAMP NOT NULL
);`

	dropUsersTableStr = "DROP TABLE IF EXISTS users;"

	createGroupsTableStr = `
CREATE TABLE IF NOT EXISTS groups (
	id    SERIAL PRIMARY KEY,
	name  TEXT NOT NULL
);`

	dropGroupsTableStr = "DROP TABLE IF EXISTS groups;"

	createGroupsUsersTableStr = `
CREATE TABLE IF NOT EXISTS groups_users(
	id         SERIAL PRIMARY KEY,
	user_id    INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
	group_id   INTEGER REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE,
	admin      BOOLEAN NOT NULL DEFAULT FALSE
);`

	dropGroupUserTableStr = "DROP TABLE IF EXISTS groups_users;"

	createExpensesTableStr = `
CREATE TABLE IF NOT EXISTS expenses(
	id          SERIAL PRIMARY KEY,
	amount      INTEGER NOT NULL CHECK (amount >= 0),
	created_at  TIMESTAMP DEFAULT LOCALTIMESTAMP NOT NULL,
	group_id    INTEGER REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE,
	payer_id    INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
	category    category_t,
	description TEXT
);`

	dropExpensesTableStr = "DROP TABLE IF EXISTS expenses;"

	createExpenseAssignmentsTableStr = `
CREATE TABLE IF NOT EXISTS expense_assignments (
	id         SERIAL PRIMARY KEY,
	user_id    INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
	amount     INTEGER NOT NULL CHECK (amount >= 0),
	expense_id INTEGER REFERENCES expenses(id) ON UPDATE CASCADE ON DELETE CASCADE
);`

	dropExpenseAssingmentsTableStr = "DROP TABLE IF EXISTS expense_assignments;"

	createPaymentsTable = `
CREATE TABLE IF NOT EXISTS payments (
	id          SERIAL PRIMARY KEY,
	created_at  TIMESTAMP DEFAULT LOCALTIMESTAMP NOT NULL,
	amount      INTEGER NOT NULL CHECK (amount >= 0),
	giver_id    INTEGER REFERENCES users(id) NOT NULL,
	receiver_id INTEGER REFERENCES users(id) CHECK (giver_id <> receiver_id),
	group_id    INTEGER REFERENCES groups(id)
);`
	dropPaymentsTableStr = "DROP TABLE IF EXISTS payments;"
)

// user query format strings

var (
	createTablesArr = []string{
		createUsersTableStr,
		createGroupsTableStr,
		createGroupsUsersTableStr,
		createExpensesTableStr,
		createExpenseAssignmentsTableStr,
		createPaymentsTable,
	}

	// Ensure reverse order to above
	dropTablesArr = []string{
		dropPaymentsTableStr,
		dropExpenseAssingmentsTableStr,
		dropExpensesTableStr,
		dropGroupUserTableStr,
		dropGroupsTableStr,
		dropUsersTableStr,
	}

	createTypesArr = []string{
		createCategoriesStr,
	}

	dropTypesArr = []string{
		dropCategoriesStr,
	}
)

type postgresStore struct {
	db    *sqlx.DB
	debug bool

	// User statements
	insertUserStmt  *sqlx.NamedStmt
	userByEmailStmt *sqlx.NamedStmt
}

func MustCreate(d *sqlx.DB) *postgresStore {
	// Set the timezone for the store. This is needed for the logic here.
	// This obviously assumes that noone else will change this.
	// TODO: make this more robust.
	db.MustExec(setTimeZoneStr)
	return &postgresStore{db: d}
}

func (s *postgresStore) mustPrepareStmts() {
	s.insertUserStmt = s.mustPrepareStmt(insertUserStr)
	s.userByEmailStmt = s.mustPrepareStmt(userByEmailStr)
}

func (s *postgresStore) mustPrepareStmt(stmt string) *sqlx.NamedStmt {
	return s.mustPrepare(s.db.PrepareNamed(stmt))
}
func (s *postgresStore) mustPrepare(stmt *sqlx.NamedStmt, err error) *sqlx.NamedStmt {
	if err != nil {
		panic("Error during prepare: " + err.Error())
	}
	return stmt
}
func (p postgresStore) MustCreateTypes() {
	p.mustExecuteStatements(createTypesArr)
}

func (p postgresStore) MustDropTypes() {
	p.mustExecuteStatements(dropTypesArr)
}

func (p postgresStore) MustCreateTables() {
	p.mustExecuteStatements(createTablesArr)
}

func (p postgresStore) MustDropTables() {
	p.mustExecuteStatements(dropTablesArr)
}

func (p postgresStore) mustExecuteStatements(statements []string) {
	for _, s := range statements {
		if p.debug {
			fmt.Println("Executing: " + s)
		}
		p.db.MustExec(s)
	}
}
