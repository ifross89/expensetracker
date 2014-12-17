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

func testExpenseCrud(st *postgresStore, t *testing.T) {
	g := &models.Group{
		Name: "TestGroup",
	}

	err := st.InsertGroup(g)
	if err != nil {
		t.Fatalf("Error inserting group: %v", err)
		return
	}

	u1 := &auth.User{
		Email:  "u1@example.com",
		PwHash: "hash",
	}

	u2 := &auth.User{
		Email:  "u2@example.com",
		PwHash: "hash",
	}

	err = st.Insert(u1)

	if err != nil {
		t.Fatalf("Could not insert user: %v", err)
		return
	}

	err = st.Insert(u2)
	if err != nil {
		t.Fatalf("Could not insert user: %v", err)
		return
	}

	err = st.AddUserToGroup(g, u1, false)
	if err != nil {
		t.Fatalf("Error adding user to group: %v", err)
		return
	}

	err = st.AddUserToGroup(g, u2, false)
	if err != nil {
		t.Fatalf("error adding user to group: %v", err)
		return
	}

	e1 := &models.Expense{
		Category:    models.CategoryDrugs,
		Amount:      100,
		GroupID:     g.ID,
		Description: "Test Expense 1",
		PayerID:     u1.ID,
	}

	var allIDs []int64 = []int64{u1.ID, u2.ID}
	t.Log(allIDs)
	var oneID []int64 = []int64{u1.ID}

	err = st.InsertExpense(e1, allIDs)
	if err != nil {
		t.Fatalf("error inserting expense: %v", err)
		return
	}

	if len(e1.Assignments) != 2 {
		t.Fatalf("Expected 2 assignments, got %d", len(e1.Assignments))
		return
	}

	if e1.Assignments[0].Amount != 50 && e1.Assignments[1].Amount != 50 {
		t.Fatalf("Assigned amounts should be 50, got: %d and %d", e1.Assignments[0].Amount, e1.Assignments[0].Amount)
		return
	}

	err = st.UpdateExpense(e1, oneID)
	if err != nil {
		t.Fatalf("Error updating expense: %v", err)
		return
	}

	if len(e1.Assignments) != 1 {
		t.Fatalf("Expected 1 assignment, got %d", len(e1.Assignments))
	}

	e2 := &models.Expense{
		Category:    models.CategoryPresents,
		Amount:      100,
		GroupID:     g.ID,
		Description: "Test expense 2",
		PayerID:     u2.ID,
	}

	err = st.InsertExpense(e2, allIDs)
	if err != nil {
		t.Fatalf("Error inserting expense: %v", err)
		return
	}

	es, err := st.ExpensesByGroup(g)
	if err != nil {
		t.Fatalf("Error getting group expenses: %+v", err)
		return
	}

	if len(es) != 2 {
		t.Fatalf("Expected 2 expenses, got %d", len(es))
		return
	}

	if es[0].Amount != 100 {
		t.Fatalf("Expense should have Amount £1, got %s", es[0].Amount)
		return
	}

	if es[1].Amount != 100 {
		t.Fatalf("Expense should have Amount £1, got %s", es[1].Amount)
		return
	}
	e1ID := e1.ID
	err = st.DeleteExpense(e1)
	if err != nil {
		t.Fatalf("Error deleting expense: %v", err)
		return
	}

	err = st.DeleteExpense(&models.Expense{ID: e1ID})
	if err == nil {
		t.Fatalf("Should have error when deleting deleted expense")
		return
	}

	e3, err := st.ExpenseByID(e2.ID)
	if err != nil {
		t.Fatalf("Error getting expense: %v", err)
		return
	}

	if e3.GroupID != e2.GroupID {
		t.Fatalf("Expense group IDs do not match: %d vs %d", e2.GroupID, e3.GroupID)
		return
	}

	_, err = st.ExpenseByID(e1ID)
	if err == nil {
		t.Fatalf("Expected error getting deleted expense")
		return
	}

	es, err = st.ExpensesByGroup(g)
	if err != nil {
		t.Fatalf("Error getting group expenses :%v", err)
		return
	}

	if len(es) != 1 {
		t.Fatalf("Should be 1 expense in group, got %v", err)
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

func benchmarkExpenseCreation(st *postgresStore, b *testing.B) {
	g := &models.Group{
		Name: "Benchmark group",
	}
	st.InsertGroup(g)

	u := &auth.User{
		Email: "Benchmark User",
	}

	st.Insert(u)
	st.AddUserToGroup(g, u, false)
	b.ResetTimer()
	uIDs := []int64{u.ID}
	for i := 0; i < b.N; i++ {
		st.InsertExpense(&models.Expense{
			PayerID:     u.ID,
			Amount:      models.Pence(int64(b.N)),
			GroupID:     g.ID,
			Category:    models.CategoryGroceries,
			Description: "TEST EXPENSE",
		}, uIDs)
	}
}

func benchmarkExpenseRetrieval(st *postgresStore, b *testing.B) {
	g := &models.Group{
		Name: "Benchmark group",
	}
	st.InsertGroup(g)

	u1 := &auth.User{
		Email: "Benchmark User 1",
	}

	u2 := &auth.User{
		Email: "Benchmark User 2",
	}

	u3 := &auth.User{
		Email: "Benchmark User 3",
	}

	st.Insert(u1)
	st.Insert(u2)
	st.Insert(u3)
	st.AddUserToGroup(g, u1, false)
	st.AddUserToGroup(g, u2, false)
	st.AddUserToGroup(g, u3, false)
	uIDs := []int64{u1.ID, u2.ID, u3.ID}
	for i := 0; i < 1000; i++ {
		st.InsertExpense(&models.Expense{
			PayerID:     uIDs[i%3],
			Amount:      5000,
			GroupID:     g.ID,
			Category:    models.CategoryGroceries,
			Description: "TEST EXPENSE",
		}, uIDs)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.ExpensesByGroup(g)
	}
}

func BenchmarkExpenseCreation(b *testing.B) {
	wrapDbBenchmark(s, benchmarkExpenseCreation)(b)
}

func BenchmarkExpenseRetrieval(b *testing.B) {
	wrapDbBenchmark(s, benchmarkExpenseRetrieval)(b)
}

func TestGroupCrud(t *testing.T) {
	wrapDbTest(s, testGroupCrud)(t)
}

func TestPaymentCrud(t *testing.T) {
	wrapDbTest(s, testPaymentCrud)(t)
}

func TestExpenseCrud(t *testing.T) {
	wrapDbTest(s, testExpenseCrud)(t)
}
