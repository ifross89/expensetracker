package models

import (
	"github.com/juju/errors"

	"time"
)

var (
	ErrAlreadySaved = errors.New("Cannot insert as model as already saved")
)

// Group represents a group of users in which the expenses are shared. An
// example of this would be housemates sharing the expenses incurred while
// living together, such as shared meals and communal home items.
type Group struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// UserGroupMap represents the database structure mapping users and groups.
// This is a many-to-many relationship. The Admin flag represents whether
// the particular user is an admin of the group. An admin has the ability
// (soon) to be able to add and remove people from the group.
type UserGroupMap struct {
	ID      int64 `db:"id" json:"id"`
	GroupID int64 `db:"group_id" json:"groupId"`
	UserID  int64 `db:"user_id" json:"userId"`
	Admin   bool  `db:"admin" json:"admin"`
}

// Payment represent a transfer of money from one person to another in the
// group. This is typically performed when one person is at a deficit overall
// to the group and another has paid a surplus with expenses.
type Payment struct {
	ID         int64     `db:"id" json:"id"`
	GroupID    int64     `db:"group_id" json:"groupId"`
	Amount     Pence     `db:"amount" json:"amount"`
	GiverID    int64     `db:"giver_id" json:"giverId"`
	ReceiverID int64     `db:"receiver_id" json:"receieverId"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
}
