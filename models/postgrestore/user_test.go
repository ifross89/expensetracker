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
		defer func() {
			if r := recover(); r != nil {
				st.MustDropTables()
				st.MustDropTypes()

				t.Failed()
			}
		}()

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

func wrapDbBenchmark(st *postgresStore, bench func(*postgresStore, *testing.B)) func(*testing.B) {
	return func(b *testing.B) {
		defer func() {
			if r := recover(); r != nil {
				st.MustDropTables()
				st.MustDropTypes()

				b.Failed()
			}
		}()

		st.MustCreateTypes()
		defer st.MustDropTypes()
		st.MustCreateTables()
		st.MustPrepareStmts()
		defer st.MustDropTables()

		bench(st, b)
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
		Name:   "TEST",
	}

	t.Log("Attempting to insert user into store")
	err := st.Insert(u)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
		return
	}

	t.Logf("User object after insertion: %#v", *u)
	if u.ID == 0 {
		t.Fatalf("user ID == 0 after insert")
		return
	}

	err = st.Insert(u)
	if err == nil {
		t.Fatalf("No error trying to insert user twice")
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

	if !isMatchingUser(t, u, u2) {
		t.Fatalf("Users do not match: user1=%+v, user2=%+v\n", u, u2)
		return
	}

	t.Log("Attempt to retrieve a non-existant user")
	u2, err = st.UserByEmail("noone@example.com")
	if err == nil {
		t.Fatalf("Expected no user found error, got nil")
		return
	}

	t.Log("Attempt to update a user")
	u.Token = "NEW TOKEN"
	err = st.Update(u)
	if err != nil {
		t.Fatalf("Error during user update: %v", err)
		return
	}

	// Now try and retrieve the user by the token
	u2, err = st.UserByToken("NEW TOKEN")
	if err != nil {
		t.Fatalf("Error getting updated user by token: %v", err)
		return
	}

	if !isMatchingUser(t, u, u2) {
		t.Fatalf("Users do not match: user1=%+v, user2=%+v\n", u, u2)
		return
	}

	u2, err = st.UserByID(u.ID)
	if err != nil {
		t.Fatalf("Error retrieving user with ID=%d", u.ID)
		return
	}
	if !isMatchingUser(t, u, u2) {
		t.Fatalf("Users do not match: user1=%+v, user%+v\n", u, u2)
	}

	//Now delete the user
	err = st.Delete(u)
	if err != nil {
		t.Fatalf("Error deleting user: %v", err)
		return
	}

	// Now attempt to delete the user again
	err = st.Delete(u2)
	if err == nil {
		t.Fatalf("No error when deleting a user twice.")
		return
	}

	// Try and get the deleted user
	_, err = st.UserByToken("NEW TOKEN")
	if err == nil {
		t.Fatalf("No error getting deleted user")
	}

}

func isMatchingUser(t *testing.T, u1, u2 *auth.User) bool {
	if u1.ID != u2.ID {
		t.Logf("Id of users differ, want %d, got %d", u1.ID, u2.ID)
		return false
	}

	if u1.Email != u2.Email {
		t.Logf("Email of users differ, want %s, got %s", u1.Email, u2.Email)
		return false
	}

	if u1.Token != u2.Token {
		t.Logf("Token of users differ, want %s got %s", u1.Token, u2.Token)
		return false
	}

	if u1.Admin != u2.Admin {
		t.Logf("Admin of users differ, want %v, got %v", u1.Admin, u2.Admin)
		return false
	}

	return true
}

func TestUserCrud(t *testing.T) {
	wrapDbTest(s, testUserCrud)(t)
}
