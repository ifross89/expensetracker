package postgrestore

import (
	"git.ianfross.com/ifross/expensetracker/auth"
	"git.ianfross.com/ifross/expensetracker/models"

	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"

	"sort"
)

const (
	// Group only strings
	insertGroupStr = `INSERT INTO groups (name) VALUES (:name) RETURNING *;`
	updateGroupStr = `UPDATE groups SET name=:name WHERE id=:id;`
	deleteGroupStr = `DELETE FROM groups where id=:id;`
	groupByIDStr   = `SELECT * FROM groups where id=:id;`
	groupByUserStr = `
SELECT groups.* FROM groups
	INNER JOIN groups_users
		ON groups_users.group_id=groups.id
	WHERE groups_users.user_id=:id;`
	allGroupsStr = `SELECT * FROM groups;`

	// Strings involving user group mappings
	addUserToGroupStr      = `INSERT INTO groups_users (group_id, user_id) VALUES (:group_id, :user_id) RETURNING *;`
	removeUserFromGroupStr = `DELETE FROM groups_users where user_id=:user_id AND group_id=:group_id;`

	// Payment strings
	insertPaymentStr = `
INSERT INTO payments (group_id, amount, giver_id, receiver_id)
	VALUES(:group_id, :amount, :giver_id, :receiver_id) RETURNING *;`
	updatePaymentStr = `
UPDATE payments SET
	group_id=:group_id,
	amount=:amount,
	giver_id=:giver_id,
	receiver_id=:receiver_id
WHERE id=:id;`
	deletePaymentStr = `DELETE FROM payments WHERE id=:id;`
	paymentByIDStr   = `SELECT * FROM payments WHERE id=:id;`

	// Expense strings
	insertExpeseStr = `
INSERT INTO expenses (amount, payer_id, group_id, category, description)
	VALUES (:amount, :payer_id, :group_id, :category, :description) RETURNING *;`
	insertExpenseAssignmentStr = `
INSERT INTO expense_assignments (amount, user_id, expense_id, group_id)
	VALUES (:amount, :user_id, :expense_id, :group_id) RETURNING *;`
	deleteExpenseStr            = `DELETE FROM expenses WHERE id=:id;`
	deleteExpenseAssignmentsStr = `DELETE FROM expense_assignments WHERE expense_id=:id;`
	updateExpenseStr            = `
UPDATE expenses set
		amount=:amount,
		payer_id=:payer_id,
		group_id=:group_id,
		category=:category
	WHERE id=:id;`

	expenseByIDStr          = `SELECT * FROM expenses WHERE id=:id;`
	assingmentsByExpenseStr = `SELECT * from expense_assignments WHERE expense_id=:id;`
	expensesByGroupStr      = `SELECT * FROM expenses WHERE group_id=:id;`
	assignmentsByGroupStr   = `SELECT * FROM expense_assignments WHERE group_id=:id;`
)

func (s *postgresStore) InsertGroup(g *models.Group) error {
	if g.ID != 0 {
		return models.ErrAlreadySaved
	}

	err := s.insertGroupStmt.Get(g, g)
	if err != nil {
		return errors.Annotate(err, "Error inserting group")
	}

	return nil
}

func (s *postgresStore) UpdateGroup(g *models.Group) error {
	r, err := s.updateGroupStmt.Exec(g)
	if err != nil {
		return errors.Annotate(err, "Error updating group")
	}

	n, _ := r.RowsAffected()
	if n != 1 {
		return errors.New("Invalid group ID")
	}

	return nil
}

func (s *postgresStore) DeleteGroup(g *models.Group) error {
	r, err := s.deleteGroupStmt.Exec(g)
	if err != nil {
		return errors.Annotate(err, "Error deleting group")
	}

	n, _ := r.RowsAffected()
	if n != 1 {
		return errors.New("No group deleted")
	}

	g.ID = 0
	return nil
}

func (s *postgresStore) GroupByID(id int64) (*models.Group, error) {
	var g = models.Group{ID: id}
	err := s.groupByIDStmt.Get(&g, g)
	if err != nil {
		return nil, errors.Annotate(err, "Error getting group by ID")
	}

	return &g, nil
}

func (s *postgresStore) GroupsByUser(u *auth.User) ([]*models.Group, error) {
	var groups []*models.Group
	err := s.groupsByUserStmt.Select(&groups, u)
	if err != nil {
		return nil, errors.Annotate(err, "Error getting user's groups")
	}

	return groups, nil
}

func (s *postgresStore) AllGroups() ([]*models.Group, error) {
	var groups []*models.Group
	err := s.db.Select(&groups, allGroupsStr)
	if err != nil {
		return nil, errors.Annotate(err, "Error getting all groups")
	}

	return groups, nil
}

func (s *postgresStore) AddUserToGroup(g *models.Group, u *auth.User, admin bool) error {
	m := models.UserGroupMap{
		GroupID: g.ID,
		UserID:  u.ID,
		Admin:   admin,
	}

	err := s.addUserToGroupStmt.Get(&m, m)
	if err != nil {
		return errors.Annotate(err, "Error adding user to group")
	}

	if m.ID == 0 {
		return errors.New("UserGroupMap has ID=0 after insertion")
	}
	return nil
}

func (s *postgresStore) RemoveUserFromGroup(g *models.Group, u *auth.User) error {
	m := models.UserGroupMap{
		GroupID: g.ID,
		UserID:  u.ID,
	}

	r, err := s.removeUserFromGroupStmt.Exec(m)
	if err != nil {
		return errors.Annotate(err, "Could not remove user from group")
	}
	n, _ := r.RowsAffected()
	if n != 1 {
		return errors.New("User not in group")
	}

	return nil
}

func (s *postgresStore) InsertPayment(p *models.Payment) error {
	err := s.insertPaymentStmt.Get(p, p)
	if err != nil {
		return errors.Annotate(err, "Error inserting payment")
	}
	return nil
}

func (s *postgresStore) UpdatePayment(p *models.Payment) error {
	r, err := s.updatePaymentStmt.Exec(p)
	if err != nil {
		return errors.Annotate(err, "Could not update payment")
	}
	n, _ := r.RowsAffected()
	if n != 1 {
		return errors.New("No payment with ID")
	}

	return nil
}

func (s *postgresStore) DeletePayment(p *models.Payment) error {
	r, err := s.deletePaymentStmt.Exec(p)
	if err != nil {
		return errors.Annotate(err, "Error deleting payment")
	}

	n, _ := r.RowsAffected()
	if n != 1 {
		return errors.New("Payment does not exist")
	}

	return nil
}

func (s *postgresStore) PaymentByID(id int64) (*models.Payment, error) {
	var p = models.Payment{
		ID: id,
	}

	err := s.paymentByIDStmt.Get(&p, p)
	if err != nil {
		return nil, errors.Annotate(err, "Error getting payment by ID")
	}
	return &p, nil
}

func (s *postgresStore) InsertExpense(e *models.Expense, userIDs []int64) error {
	// Assign expense and commit everything to the db within the same transaction
	if e.ID != 0 {
		return models.ErrAlreadySaved
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return errors.Annotate(err, "Could not create transaction")
	}

	stmt, err := tx.PrepareNamed(insertExpeseStr)
	if err != nil {
		_ = tx.Rollback()
		return errors.Annotate(err, "Error preparing insert expense statement")
	}
	err = stmt.Get(e, e)
	if err != nil {
		_ = tx.Rollback()
		return errors.Annotate(err, "Error inserting expense")
	}

	eas, err := e.Assign(userIDs)
	if err != nil {
		_ = tx.Rollback()
		return errors.Annotate(err, "Error assigning expense")
	}

	err = s.insertExpenseAssignments(eas, tx)
	if err != nil {
		_ = tx.Rollback()
		return errors.Trace(err)
	}

	// Sucessfully inserted expense and assignments.
	err = tx.Commit()
	if err != nil {
		return errors.Annotate(err, "Error committing to database")
	}

	e.Assignments = eas

	return nil
}

func (s *postgresStore) UpdateExpense(e *models.Expense, userIDs []int64) error {
	if e.ID == 0 {
		return models.ErrStructNotSaved
	}

	eas, err := e.Assign(userIDs)
	if err != nil {
		return errors.Annotate(err, "Could not assign expense")
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return errors.Annotate(err, "Could not create transaction")
	}

	stmt, err := tx.PrepareNamed(deleteExpenseAssignmentsStr)
	if err != nil {
		_ = tx.Rollback()
		return errors.Annotate(err, "Error preparing delete expense assignments statement")
	}

	r, err := stmt.Exec(e)
	if err != nil {
		_ = tx.Rollback()
		return errors.Annotate(err, "Error deleting expense assignments")
	}

	n, _ := r.RowsAffected()
	if n == 0 {
		_ = tx.Rollback()
		return errors.New("expense does not have any associated assignments")
	}

	stmt, err = tx.PrepareNamed(updateExpenseStr)
	if err != nil {
		_ = tx.Rollback()
		return errors.Annotate(err, "Error preparing update expense statement")
	}

	r, err = stmt.Exec(e)
	if err != nil {
		_ = tx.Rollback()
		return errors.Annotate(err, "Error updating expense")
	}

	err = s.insertExpenseAssignments(eas, tx)
	if err != nil {
		_ = tx.Rollback()
		return errors.Trace(err)
	}

	// updated expense and created new assignments
	err = tx.Commit()
	if err != nil {
		return errors.Annotate(err, "error committing expense update")
	}

	e.Assignments = eas
	return nil
}

func (s *postgresStore) insertExpenseAssignments(eas []*models.ExpenseAssignment, tx *sqlx.Tx) error {
	stmt, err := tx.PrepareNamed(insertExpenseAssignmentStr)
	if err != nil {
		return errors.Annotate(err, "Error preparing insert expense assigment statement")
	}

	for _, ea := range eas {
		err = stmt.Get(ea, ea)
		if err != nil {
			return errors.Annotate(err, "Error inserting expense assignment")
		}
	}

	return nil
}

func (s *postgresStore) ExpenseByID(id int64) (*models.Expense, error) {
	e := models.Expense{
		ID: id,
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errors.Annotate(err, "could not create transaction")
	}

	stmt, err := tx.PrepareNamed(expenseByIDStr)
	if err != nil {
		_ = tx.Rollback()
		return nil, errors.Annotate(err, "could not create expense by ID statement")
	}

	err = stmt.Get(&e, e)
	if err != nil {
		_ = tx.Rollback()
		return nil, errors.Annotate(err, "could not get expense by id")
	}

	var eas []*models.ExpenseAssignment
	stmt, err = tx.PrepareNamed(assingmentsByExpenseStr)
	if err != nil {
		_ = tx.Rollback()
		return nil, errors.Annotate(err, "could not prepare assignments by expense statement")
	}

	err = stmt.Select(&eas, e)
	if err != nil {
		_ = tx.Rollback()
		return nil, errors.Annotate(err, "could not get assignments for expense")
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Annotate(err, "could not commit")
	}

	e.Assignments = eas

	return &e, nil

}

func (s *postgresStore) ExpensesByGroup(g *models.Group) ([]*models.Expense, error) {
	var es []*models.Expense
	var eas []*models.ExpenseAssignment
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errors.Annotate(err, "could not create transaction")
	}

	stmt, err := tx.PrepareNamed(expensesByGroupStr)
	if err != nil {
		_ = tx.Rollback()
		return nil, errors.Trace(err)
	}

	err = stmt.Select(&es, g)
	if err != nil {
		_ = tx.Rollback()
		return nil, errors.Trace(err)
	}

	stmt, err = tx.PrepareNamed(assignmentsByGroupStr)
	if err != nil {
		_ = tx.Rollback()
		return nil, errors.Trace(err)
	}

	err = stmt.Select(&eas, g)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// got expenses and assignments
	err = tx.Commit()
	if err != nil {
		return nil, errors.Trace(err)
	}

	// Pair the assignments with the expense
	for _, e := range es {
		e.Assignments = make([]*models.ExpenseAssignment, 0, 0)
	}

	// Assignments need to be sorted by group
	sort.Sort(models.ByExpense(eas))

	eIndex := int64(-1)
	eID := int64(-1)
	for _, ea := range eas {
		if eID != ea.ExpenseID {
			//
			eIndex++
			eID = ea.ExpenseID
		}

		es[eIndex].Assignments = append(es[eIndex].Assignments, ea)
	}

	return es, nil
}

func (s *postgresStore) DeleteExpense(e *models.Expense) error {
	r, err := s.deleteExpenseStmt.Exec(e)

	// Delete expense deletes all associated assignments due to CASCADE
	if err != nil {
		return errors.Annotatef(err, "Could not delete expense with ID=%d", e.ID)
	}

	n, _ := r.RowsAffected()
	if n == 0 {
		return errors.New("Expense does not exist")
	}

	e.ID = 0
	return nil
}
