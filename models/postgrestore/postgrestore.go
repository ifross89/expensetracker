package postgrestore

import (
	"git.ianfross.com/expensetracker/models"

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
	'tickets');`

	dropCategoriesStr = "DROP TYPE category_t;"

	createUsersTableStr = `
CREATE TABLE users (
	id                  SERIAL PRIMARY KEY,
	email               VARCHAR(64) NOT NULL CHECK (email <> '') UNIQUE,
	hashedPassword      VARCHAR(128),
	isAdmin             BOOLEAN NOT NULL DEFAULT false,
	isActive            BOOLEAN NOT NULL DEFAULT false,
	signupToken         TEXT,
	passwordChangeToken TEXT,
	isNew               BOOLEAN default true);`

	dropUsersTableStr = "DROP TABLE users;"

	createGroupsTableStr = `
CREATE TABLE groups (
	id    SERIAL PRIMARY KEY,
	name  TEXT NOT NULL);`

	dropGroupsTableStr = "DROP TABLE groups;"

	createGroupsUsersTableStr = `
CREATE TABLE groups_users(
	id        SERIAL PRIMARY KEY,
	userId    INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
	groupId   INTEGER REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE);`

	dropGroupUserTableStr = "DROP TABLE groups_users;"

	createExpensesTableStr = `
CREATE TABLE expenses(
	id        SERIAL PRIMARY KEY,
	createdAt TIMESTAMP DEFAULT LOCALTIMESTAMP NOT NULL,
	time      TIMESTAMP DEFAULT LOCALTIMESTAMP NOT NULL,
	groupId   REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE,
	payerId   REFERENCES payers(id) ON UPDATE CASCADE ON DELETE CASCADE,
	category  category_t,
	description TEXT);`

	dropExpensesTableStr = "DROP TABLE expenses;"

	createExpenseAssignmentsTableStr = `
CREATE TABLE expense_assignments (
	id        SERIAL PRIMARY KEY,
	userId    INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
	amount    INTEGER NOT NULL CHECK (amount >= 0),
	expenseId INTEGERE REFERENCES expenses(id) ON UPDATE CASCADE ON DELETE CASCADE);`

	dropExpenseAssingmentsTableStr = "DROP TABLE expense_assignments;"

	createPaymentsTable = `
CREATE TABLE payments (
	id         SERIAL PRIMARY KEY,
	createdOn  TIMESTAMP DEFAULT LOCALTIMESTAMP NOT NULL,
	amount     INTEGER NOT NULL CHECK (amount >= 0),
	payerId    INTEGER REFERENCES users(id) NOT NULL,
	receiverId INTEGER REFERENCES users(id) CHECK (payerId <> receiverId),
	circleID   INTEGER REFERENCES groups(id));`

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

func (p postgresStore)
