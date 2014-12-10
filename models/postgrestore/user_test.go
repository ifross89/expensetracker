package postgrestore

import (
	"git.ianfross.com/expensetracker/auth"

	"testing"
	"time"
)

// wrapDbTest wraps a test to ensure the database is created before the test
// and destroyed after the end of the test. This means there is no data left
// in the database between test runs.runs
func wrapDbTest(st *postgresStore, test func(*postgresStore, *testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		st.debug = true
		defer func() { st.debug = false }()
		// Create database schema
		st.MustCreateTypes()
		defer st.MustDropTypes()
		st.MustCreateTables()
		defer st.MustDropTables()

		// perform the test
		test(st, t)
	}
}

func testUserInsert(st *postgresStore, t *testing.T) {
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
	if time.Now().UTC().Before(*(u.CreatedAt)) || (*(u.CreatedAt)).Add(5*time.Second).Before(time.Now().UTC()) {
		t.Fatalf("CreatedAt time too different from now. Now=%v, CreatedAt=%v", time.Now().UTC(), *(u.CreatedAt))
	}
}

func TestUserInsert(t *testing.T) {
	wrapDbTest(s, testUserInsert)(t)
}
