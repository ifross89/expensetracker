package models

import (
	"git.ianfross.com/expensetracker/auth"

	"github.com/juju/errors"

	"time"
)

type Group struct {
	Id   int64
	Name string
}

type UserGroupMap struct {
	Id      int64
	GroupId int64
	UserId  int64
	IsAdmin bool
}

type Payment struct {
	Id         int64
	GroupId    int64
	Amount     Pence
	GiverId    int64
	ReceiverId int64
	CreatedAt  time.Time
}

type Storer interface {
	// Group storage functions
	InsertGroup(*Group) error
	UpdateGroup(*Group) error
	DeleteGroup(*Group) error
	GroupById(int64) (*Group, error)
	AddUserToGroup(*Group, *auth.User, bool) error
	RemoveUserFromGroup(*Group, *auth.User) error

	// Expense storage functions
	InsertExpense(*Expense) error // Need to fill in Id and Assignments
	UpdateExpense(*Expense) error
	ExpenseById(int64) (*Expense, error)
	DeleteExpense(*Expense) error

	// Payment storage functions
	InsertPayment(*Payment) error
	UpdatePayment(*Payment) error
	DeletePayment(*Payment) error
	PaymentById(int64) (*Payment, error)
}

type Manager struct {
	store Storer
}

func NewManager(s Storer) *Manager {
	return &Manager{s}
}

func (m *Manager) NewGroup(name string) (*Group, error) {
	g := &Group{Name: name}
	return g, errors.Trace(m.store.InsertGroup(g))
}

func (m *Manager) DeleteGroup(g *Group) error {
	return errors.Trace(m.store.DeleteGroup(g))
}

func (m *Manager) UpdateGroup(g *Group) error {
	return errors.Trace(m.store.UpdateGroup(g))
}

func (m *Manager) GroupById(id int64) (*Group, error) {
	g, err := m.store.GroupById(id)
	return g, errors.Trace(err)
}

func (m *Manager) AddUserToGroup(g *Group, u *auth.User, admin bool) error {
	return errors.Trace(m.store.AddUserToGroup(g, u, admin))
}

func (m *Manager) RemoveUserFromGroup(g *Group, u *auth.User) error {
	return errors.Trace(m.store.RemoveUserFromGroup(g, u))
}

func (m *Manager) NewExpense(g *Group, amount Pence, payer int64, cat Category, desc string, users []int64) (*Expense, error) {
	e := &Expense{
		Amount:      amount,
		PayerId:     payer,
		Category:    cat,
		Description: desc,
		GroupId:     g.Id,
	}

	if err := m.store.InsertExpense(e); err != nil {
		return nil, errors.Annotate(err, "Unable to insert expense")
	}

	return e, nil
}

func (m *Manager) UpdateExpense(e *Expense) error {
	// the storage function needs to remove all the assignments
	// and reassign the expense within a transaction. This
	// is to ensure consistency within the database.
	return errors.Trace(m.store.UpdateExpense(e))
}

func (m *Manager) DeleteExpense(e *Expense) error {
	// Deletes all associated assignments
	return errors.Trace(m.store.DeleteExpense(e))
}

func (m *Manager) InsertPayment(g *Group, giver, receiver int64, amount Pence) (*Payment, error) {
	p := &Payment{
		Amount:     amount,
		GroupId:    g.Id,
		GiverId:    giver,
		ReceiverId: receiver,
	}
	err := m.store.InsertPayment(p)
	if err != nil {
		return nil, errors.Annotate(err, "Error inserting payment")
	}

	return p, nil
}

func (m *Manager) DeletePayment(p *Payment) error {
	return errors.Trace(m.store.DeletePayment(p))
}

func (m *Manager) UpdatePayment(p *Payment) error {
	return errors.Trace(m.store.UpdatePayment(p))
}

func (m *Manager) PaymentById(id int64) (*Payment, error) {
	p, err := m.store.PaymentById(id)
	if err != nil {
		return nil, errors.Annotate(err, "Could not retrieve payment")
	}
	return p, nil
}
