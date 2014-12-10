package postgrestore

import (
	"git.ianfross.com/expensetracker/auth"

	_ "github.com/juju/errors"

	"testing"
	"time"
)

// wrapDbTest wraps a test to ensure the database is created before the test
// and destroyed after the end of the test. This means there is no data left
// in the database between test runs.runs
func wrapDbTest(st *postgresStore, test func(*postgresStore, *testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		st.debug = false
		defer func() { st.debug = false }()
		// Create database schema
		st.MustCreateTypes()
		defer st.MustDropTypes()
		st.MustCreateTables()
		st.MustPrepareStmts()
		defer st.MustDropTables()

		// perform the test
		test(st, t)
	}
}

func testUserCrud(st *postgresStore, t *testing.T) {
	t.Log("Creating user")
	u := &auth.User{
		Email:  "hello@example.com",
		PwHash: "exampleHash",
		Admin:  true,
		Active: true,
		Token:  "TOKEN",
	}

	t.Log("Attempting to insert user into store")
	err := st.Insert(u)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
		return
	}
	t.Logf("User object after insertion: %#v", *u)
	if u.Id == 0 {
		t.Fatalf("user ID == 0 after insert")
		return
	}

	if u.CreatedAt == nil {
		t.Fatalf("user created time not updated after insert")
		return
	}
	// check that the time is in the past, but not longer than 5s ago
	if time.Now().UTC().Before(*(u.CreatedAt)) ||
		(*(u.CreatedAt)).Add(5*time.Second).Before(time.Now().UTC()) {
		t.Fatalf("CreatedAt time too different from now. Now=%v, CreatedAt=%v",
			time.Now().UTC(), *(u.CreatedAt))
	}

	t.Log("Attempt to retrieve the user stored")
	u2, err := st.UserByEmail("hello@example.com")
	if err != nil {
		t.Fatalf("Error retrieving user: %v", err)
		return
	}

	if u.Id != u2.Id {
		t.Fatalf("Id of users differ, want %d, got %d", u.Id, u2.Id)
		return
	}

	if u.Email != u2.Email {
		t.Fatalf("Email of users differ, want %s, got %s", u.Email, u2.Email)
		return
	}

	if u.Token != u2.Token {
		t.Fatalf("Token of users differ, want %s got %s", u.Token, u2.Token)
		return
	}

	if u.Admin != u2.Admin {
		t.Fatalf("Admin of users differ, want %v, got %v", u.Admin, u2.Admin)
		return
	}

	t.Log("Attempt to retrieve a non-existant user")
	u2, err = st.UserByEmail("noone@example.com")
	if err == nil {
		t.Fatalf("Expected no user found error, got nil")
		return
	}
}

func TestUserCrud(t *testing.T) {
	wrapDbTest(s, testUserCrud)(t)
}
