package env

import (
	"git.ianfross.com/ifross/expensetracker/auth"
	"git.ianfross.com/ifross/expensetracker/models"
)

type Config struct {
	Port int
}

type Env struct {
	*models.Manager
	*auth.UserManager
	Conf Config
}
