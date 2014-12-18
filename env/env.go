package env

import (
	"git.ianfross.com/expensetracker/auth"
	"git.ianfross.com/expensetracker/models"
)

type Config struct {
	Port int
}

type Env struct {
	*models.Manager
	*auth.UserManager
	Conf Config
}
