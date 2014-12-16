package postgrestore

import (
	"git.ianfross.com/expensetracker/auth"
	"git.ianfross.com/expensetracker/models"

	"fmt"
	"github.com/juju/errors"
)

const (
	// Group only strings
	insertGroupStr = `INSERT INTO groups (name) VALUES (:name) RETURNING *;`
	updateGroupStr = `UPDATE groups SET name=:name WHERE id=:id;`
	deleteGroupStr = `DELETE FROM groups where id=:id;`
	groupByIDStr   = `SELECT * FROM groups where id=:id;`

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
	VALUES (:amount, :payer_id, :group_id, :category) RETURNING *;`
	insertExpenseAssignmentStr = `
INSERT INTO expense_assignments (amount, user_id, expense_id)
	VALUES (:amount, :user_id, :expense_id) RETURNING *;`
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
	fmt.Printf("InsertPayment p=%+v\n", p)
	err := s.insertPaymentStmt.Get(p, p)
	if err != nil {
		return errors.Annotate(err, "Error inserting payment")
	}
	return nil
}

func (s *postgresStore) UpdatePayment(p *models.Payment) error {
	fmt.Printf("UpdatePayment p=%+v\n", p)
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
	fmt.Println("DeletePayment called with ID=", p.ID)
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
	fmt.Printf("PaymentByID: got payment=%+v\n", p)
	return &p, nil
}

func (s *postgresStore) InsertExpense(e *models.Expense, userIds []int64) error {
	// Assign expense and commit everything to the db within the same transaction
	if e.ID != 0 {
		return models.ErrAlreadySaved
	}

	eas, err := e.Assign(userIds)
	if err != nil {
		return errors.Annotate(err, "Error assigning expense")
	}

	tx, err := db.Beginx()
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

	stmt, err = tx.PrepareNamed(insertExpenseAssignmentStr)
	if err != nil {
		_ = tx.Rollback()
		return errors.Annotate(err, "Error preparing insert expense assigment statement")
	}

	for _, ea := range eas {
		err = stmt.Get(ea, ea)
		if err != nil {
			_ = tx.Rollback()
			return errors.Annotate(err, "Error inserting expense assignment")
		}
	}
	// Sucessfully inserted expense and assignments.
	err = tx.Commit()
	if err != nil {
		return errors.Annotate(err, "Error committing to database")
	}

	return nil
}
