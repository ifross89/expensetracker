package handlers

import (
	"git.ianfross.com/ifross/expensetracker/auth"
	"git.ianfross.com/ifross/expensetracker/models"
	"git.ianfross.com/ifross/expensetracker/routeindex"
)

type ExpenseGetHandler struct {
	Um    *auth.UserManager
	M     *models.Manager
	index routeindex.Interface
}
