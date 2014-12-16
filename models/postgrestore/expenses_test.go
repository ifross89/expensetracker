package postgrestore

import (
	"git.ianfross.com/expensetracker/auth"
	"git.ianfross.com/expensetracker/models"

	"testing"
)

func testPaymentCrud(st *postgresStore, t *testing.T) {
	t.Log("Create necessary group and users")
	u := &auth.User{
		Email:  "hello@example.com",
		PwHash: "hash",
	}

	err := st.Insert(u)

	if err != nil {
		t.Fatalf("Error creating user: %v", err)
		return
	}

	u2 := &auth.User{
		Email:  "hello2@example.com",
		PwHash: "hash",
	}

	err = st.Insert(u2)
	if err != nil {
		t.Fatalf("Error creating user %v", err)
		return
	}

	g := &models.Group{
		Name: "test group",
	}

	err = st.InsertGroup(g)
	if err != nil {
		t.Fatalf("Error creating group %v", err)
		return
	}

	t.Log("Creating payment")
	p := &models.Payment{
		GroupID:    g.ID,
		GiverID:    u.ID,
		ReceiverID: u2.ID,
		Amount:     100,
	}

	err = st.InsertPayment(p)
	if err != nil {
		t.Fatalf("Error inserting payment: %v", err)
		return
	}

	t.Log("Create and save a second payment with same information")
	p.ID = 0
	err = st.InsertPayment(p)
	if err != nil {
		t.Fatalf("Error inserting payment: %v", err)
		return
	}

	p.Amount = 200
	err = st.UpdatePayment(p)
	if err != nil {
		t.Fatalf("Error updating payment: %v", err)
		return
	}

	p2, err := st.PaymentByID(p.ID)
	if err != nil {
		t.Fatalf("Error getting payment by ID: %v", err)
		return
	}

	if p.Amount != p2.Amount {
		t.Fatalf("Expected payments to have the same amount")
		return
	}

	err = st.DeletePayment(p)
	if err != nil {
		t.Fatalf("Error deleting payment: %v", err)
		return
	}

	p, err = st.PaymentByID(p2.ID)
	if err == nil {
		t.Fatalf("Expected error when getting deleted payment")
		return
	}

	err = st.UpdatePayment(p2)
	if err == nil {
		t.Fatalf("Expected error updating deleted payment")
		return
	}
}

func testGroupCrud(st *postgresStore, t *testing.T) {
	t.Log("Creating group")
	g := &models.Group{
		Name: "Test group",
	}

	err := st.InsertGroup(g)
	if err != nil {
		t.Fatalf("Error inserting group: %v", err)
		return
	}

	if g.ID == 0 {
		t.Fatalf("ID of saved group==0")
		return
	}

	g2, err := st.GroupByID(g.ID)
	if err != nil {
		t.Fatalf("Could not get group by err=%v", err)
		return
	}

	if g.ID != g2.ID || g.Name != g2.Name {
		t.Fatalf("Group retrieved by ID does not match group inserted. g1=%+v, g2=%+v", g, g2)
		return
	}

	g.Name = "Updated group name"
	err = st.UpdateGroup(g)
	if err != nil {
		t.Fatalf("Error updating group: %v", err)
		return
	}

	g2, err = st.GroupByID(g.ID)
	if err != nil {
		t.Fatalf("Error getting group with updated name: %v", err)
		return
	}

	if g2.Name != g.Name {
		t.Fatalf("Name different when getting updated group: g1=%+v, g2=%+v", g, g2)
		return
	}

	u := &auth.User{
		Email:  "test@example.com",
		PwHash: "hash",
		Token:  "token",
		Active: true,
		Admin:  false,
	}

	err = st.AddUserToGroup(g, u, true)
	if err == nil {
		t.Fatalf("Expected error adding non saved user to group")
		return
	}

	err = st.Insert(u)
	if err != nil {
		t.Fatalf("Error saving user: %v", err)
		return
	}

	err = st.AddUserToGroup(g, u, true)
	if err != nil {
		t.Fatalf("Error adding user to group: %v", err)
		return
	}

	err = st.AddUserToGroup(g, u, true)
	if err == nil {
		t.Fatalf("Expected error adding user to group twice")
		return
	}

	err = st.RemoveUserFromGroup(g, u)
	if err != nil {
		t.Fatalf("Error removing user from group: %v", err)
		return
	}

	err = st.RemoveUserFromGroup(g, u)
	if err == nil {
		t.Fatalf("Expected error removing user from group twice", err)
		return
	}

	err = st.DeleteGroup(g)
	if err != nil {
		t.Fatalf("Error deleting group: %v", err)
		return
	}

	g.Name = ""
	g.ID = 0
	g, err = st.GroupByID(g2.ID)
	if err == nil {
		t.Errorf("No error getting deleted group")
		return
	}
}

func TestGroupCrud(t *testing.T) {
	wrapDbTest(s, testGroupCrud)(t)
}

func TestPaymentCrud(t *testing.T) {
	wrapDbTest(s, testPaymentCrud)(t)
}
