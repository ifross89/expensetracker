package models

import (
	"git.ianfross.com/expensetracker/auth"

	"github.com/juju/errors"

	"time"
)

// Group represents a group of users in which the expenses are shared. An
// example of this would be housemates sharing the expenses incurred while
// living together, such as shared meals and communal home items.
type Group struct {
	ID   int64
	Name string
}

// UserGroupMap represents the database structure mapping users and groups.
// This is a many-to-many relationship. The Admin flag represents whether
// the particular user is an admin of the group. An admin has the ability
// (soon) to be able to add and remove people from the group.
type UserGroupMap struct {
	ID      int64 `db:"id"`
	GroupID int64 `db:"group_id"`
	UserID  int64 `db:"user_id"`
	Admin   bool  `db:"admin"`
}

// Payment represent a transfer of money from one person to another in the
// group. This is typically performed when one person is at a deficit overall
// to the group and another has paid a surplus with expenses.
type Payment struct {
	ID         int64
	GroupID    int64
	Amount     Pence
	GiverID    int64
	ReceiverID int64
	CreatedAt  time.Time
}

// Storer is the interface required in order to perform the actions required
// to store persist the models defined in the package. Any particular behavior
// that is not described by the type system will be explained in the comments
// above each method.
type Storer interface {
	// Group storage functions
	InsertGroup(*Group) error
	UpdateGroup(*Group) error
	DeleteGroup(*Group) error
	GroupByID(int64) (*Group, error)
	AddUserToGroup(*Group, *auth.User, bool) error
	RemoveUserFromGroup(*Group, *auth.User) error

	// Expense storage functions
	InsertExpense(*Expense, []int64) error // Need to fill in Id and Assignments
	UpdateExpense(*Expense, []int64) error
	ExpenseByID(int64) (*Expense, error)
	DeleteExpense(*Expense) error

	// Payment storage functions
	InsertPayment(*Payment) error
	UpdatePayment(*Payment) error
	DeletePayment(*Payment) error
	PaymentByID(int64) (*Payment, error)
}

// Manager contains the methods that are available to the models in the. The
// manager needs to be created with a Storer interface, which deals with the
// persistence of the structs. Actions built on these persistence methods
// are available for use, for example in HTTP handlers.
type Manager struct {
	store Storer
}

// NewManager creates a new instance of the Manager object.
func NewManager(s Storer) *Manager {
	return &Manager{s}
}

// NewGroup creates and persists a new group with the name supplied.
func (m *Manager) NewGroup(name string) (*Group, error) {
	g := &Group{Name: name}
	return g, errors.Trace(m.store.InsertGroup(g))
}

// DeleteGroup removes the group from storage and any mappings to the members
// of the group.
func (m *Manager) DeleteGroup(g *Group) error {
	return errors.Trace(m.store.DeleteGroup(g))
}

// UpdateGroup saves any changes made to the group object. If the ID has been
// changed then this will fail or over-write another existing group!
func (m *Manager) UpdateGroup(g *Group) error {
	return errors.Trace(m.store.UpdateGroup(g))
}

// GroupByID retrieves a group from persistence by the ID supplied.
func (m *Manager) GroupByID(id int64) (*Group, error) {
	g, err := m.store.GroupByID(id)
	return g, errors.Trace(err)
}

// AddUserToGroup associates a user to the group. This is done internally by
// creating a mapping between the user and the group.
func (m *Manager) AddUserToGroup(g *Group, u *auth.User, admin bool) error {
	return errors.Trace(m.store.AddUserToGroup(g, u, admin))
}

// RemoveUserFromGroup dissociates a user from the group. Any expense
// assignments are deleted. Payments remain as these apply to other in
// the group, but this user must not be taken into account in any calculations
func (m *Manager) RemoveUserFromGroup(g *Group, u *auth.User) error {
	return errors.Trace(m.store.RemoveUserFromGroup(g, u))
}

// NewExpense creates, assigns and persists a new expense. The assignments
// should be created using AssignExpense.
// For consistency, the expense and the assignments need to occur
// transactionally. i.e. they must all be persisted, or all rollback. This
// is to ensure a consisent state in the database. This is the reason that
// AssignExpense cannot be used to persist the assignments, as this must
// be called within the transaction. This can only be guaranteed at the
// storage driver level (i.e. the implementation of the Storer interface)
func (m *Manager) NewExpense(g *Group, amount Pence, payer int64, cat Category, desc string, users []int64) (*Expense, error) {
	e := &Expense{
		Amount:      amount,
		PayerID:     payer,
		Category:    cat,
		Description: desc,
		GroupID:     g.ID,
	}

	if err := m.store.InsertExpense(e, users); err != nil {
		return nil, errors.Annotate(err, "Unable to insert expense")
	}

	return e, nil
}

// UpdateExpense saves any changes to the expense. If there are changes
// made to the amount or number of people, then all previous assignments
// must be removed and this must be reassigned. This must all happen within
// a transaction.
func (m *Manager) UpdateExpense(e *Expense, users []int64) error {
	// the storage function needs to remove all the assignments
	// and reassign the expense within a transaction. This
	// is to ensure consistency within the database.
	return errors.Trace(m.store.UpdateExpense(e, users))
}

// DeleteExpense removes an expense and any assignments associated with the
// expense. This must be transactional.
func (m *Manager) DeleteExpense(e *Expense) error {
	// Deletes all associated assignments
	return errors.Trace(m.store.DeleteExpense(e))
}

// InsertPayment persists a payment of money from one person to another within
// a group.
func (m *Manager) InsertPayment(g *Group, giver, receiver int64, amount Pence) (*Payment, error) {
	p := &Payment{
		Amount:     amount,
		GroupID:    g.ID,
		GiverID:    giver,
		ReceiverID: receiver,
	}
	err := m.store.InsertPayment(p)
	if err != nil {
		return nil, errors.Annotate(err, "Error inserting payment")
	}

	return p, nil
}

// DeletePayment removes a Payment from storage
func (m *Manager) DeletePayment(p *Payment) error {
	return errors.Trace(m.store.DeletePayment(p))
}

// UpdatePayment saves any modifications to the payment
func (m *Manager) UpdatePayment(p *Payment) error {
	return errors.Trace(m.store.UpdatePayment(p))
}

// PaymentByID returns a payment object with the given ID.
func (m *Manager) PaymentByID(id int64) (*Payment, error) {
	p, err := m.store.PaymentByID(id)
	if err != nil {
		return nil, errors.Annotate(err, "Could not retrieve payment")
	}
	return p, nil
}
