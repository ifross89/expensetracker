package postgrestore

import (
	"github.com/jmoiron/sqlx"

	"fmt"
)

const (
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
	token               TEXT
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
}

func Create(d *sqlx.DB) *postgresStore {
	return &postgresStore{db: d}
}

func (p postgresStore) MustCreateTypes() {
	for _, s := range createTypesArr {
		if p.debug {
			fmt.Println("Executing: " + s)
		}
		p.db.MustExec(s)
	}
}

func (p postgresStore) MustDropTypes() {
	for _, s := range dropTypesArr {
		if p.debug {
			fmt.Println("Executing: " + s)
		}
		p.db.MustExec(s)
	}
}

func (p postgresStore) MustCreateTables() {
	for _, s := range createTablesArr {
		if p.debug {
			fmt.Println("Executing: " + s)
		}
		p.db.MustExec(s)
	}
}

func (p postgresStore) MustDropTables() {
	for _, s := range dropTablesArr {
		if p.debug {
			fmt.Println("Executing: " + s)
		}
		p.db.MustExec(s)
	}
}
