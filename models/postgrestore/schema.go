package postgrestore

import (
	"github.com/jmoiron/sqlx"

	"database/sql"
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
CREATE TABLE users (
	id                  SERIAL PRIMARY KEY,
	email               VARCHAR(64) NOT NULL CHECK (email <> '') UNIQUE,
	pw_hash             VARCHAR(128),
	admin               BOOLEAN NOT NULL DEFAULT false,
	active              BOOLEAN NOT NULL DEFAULT false,
	token               TEXT
);`

	dropUsersTableStr = "DROP TABLE users;"

	createGroupsTableStr = `
CREATE TABLE groups (
	id    SERIAL PRIMARY KEY,
	name  TEXT NOT NULL
);`

	dropGroupsTableStr = "DROP TABLE groups;"

	createGroupsUsersTableStr = `
CREATE TABLE groups_users(
	id         SERIAL PRIMARY KEY,
	user_id    INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
	group_id   INTEGER REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE,
	admin      BOOLEAN NOT NULL DEFAULT FALSE
);`

	dropGroupUserTableStr = "DROP TABLE groups_users;"

	createExpensesTableStr = `
CREATE TABLE expenses(
	id          SERIAL PRIMARY KEY,
	amount      INTEGER NOT NULL CHECK (amount >= 0),
	created_at  TIMESTAMP DEFAULT LOCALTIMESTAMP NOT NULL,
	group_id    REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE,
	payer_id    REFERENCES payers(id) ON UPDATE CASCADE ON DELETE CASCADE,
	category    category_t,
	description TEXT
);`

	dropExpensesTableStr = "DROP TABLE expenses;"

	createExpenseAssignmentsTableStr = `
CREATE TABLE expense_assignments (
	id         SERIAL PRIMARY KEY,
	user_id    INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
	amount     INTEGER NOT NULL CHECK (amount >= 0),
	expense_id INTEGER REFERENCES expenses(id) ON UPDATE CASCADE ON DELETE CASCADE
);`

	dropExpenseAssingmentsTableStr = "DROP TABLE expense_assignments;"

	createPaymentsTable = `
CREATE TABLE payments (
	id          SERIAL PRIMARY KEY,
	created_at  TIMESTAMP DEFAULT LOCALTIMESTAMP NOT NULL,
	amount      INTEGER NOT NULL CHECK (amount >= 0),
	giver_id    INTEGER REFERENCES users(id) NOT NULL,
	receiver_id INTEGER REFERENCES users(id) CHECK (payerId <> receiverId),
	group_id    INTEGER REFERENCES groups(id)
);`

	dropPaymentsTableStr = "DROP TABLE payments;"
)

// user query format strings

var (
	createSchemaArr = []string{
		createCategoriesStr,
		createUsersTableStr,
		createGroupsTableStr,
		createGroupsUsersTableStr,
		createExpensesTableStr,
		createExpenseAssignmentsTableStr,
		createPaymentsTable,
	}

	// Ensure reverse order to above
	dropSchemaArr = []string{
		dropPaymentsTableStr,
		dropExpenseAssingmentsTableStr,
		dropExpensesTableStr,
		dropGroupUserTableStr,
		dropGroupsTableStr,
		dropUsersTableStr,
		dropCategoriesStr,
	}
)

type postgresStore struct {
	db *sqlx.DB
}

func Create(d *sql.DB) *postgresStore {
	return &postgresStore{db: sqlx.NewDb(d, "postgres")}
}

func (p postgresStore) MustCreateSchema() {
	for _, s := range createSchemaArr {
		p.db.MustExec(s)
	}
}

func (p postgresStore) MustDropSchema() {
	for _, s := range dropSchemaArr {
		p.db.MustExec(s)
	}
}
