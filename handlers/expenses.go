package handlers

import (
	"git.ianfross.com/expensetracker/auth"
	"git.ianfross.com/expensetracker/models"
	"git.ianfross.com/expensetracker/routeindex"
)

type ExpenseGetHandler struct {
	Um    *auth.UserManager
	M     *models.Manager
	index routeindex.Interface
}
