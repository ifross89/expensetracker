package main

import (
	"git.ianfross.com/ifross/expensetracker/auth"
	"git.ianfross.com/ifross/expensetracker/env"
	"git.ianfross.com/ifross/expensetracker/handlers"
	"git.ianfross.com/ifross/expensetracker/models"
	"git.ianfross.com/ifross/expensetracker/models/postgrestore"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"

	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type actionsMap map[string]func() error

func (a actionsMap) available() string {
	var actions []string
	for k, _ := range a {
		actions = append(actions, k)
	}
	return "[" + strings.Join(actions, ", ") + "]"
}

func (a actionsMap) validAction(action string) bool {
	_, ok := a[action]
	return ok
}

func (a actionsMap) perform(action string) error {
	return a[action]()
}

var (
	dbUser = flag.String("db_user", "expensetracker", "database user to connect with")
	dbName = flag.String("db_name", "expensetracker", "name of the database to connect to")
	dbPw   = flag.String("db_pw", "", "user's database password")
	dbHost = flag.String("db_host", "localhost", "host the database is running on")
	dbPort = flag.Int("db_port", 5432, "port the database is listening on")

	port   = flag.Int("port", 8181, "HTTP port to listen on")
	action = flag.String("action", "start", "action to perform. Available: "+actions.available())

	adminName  = flag.String("admin_name", "", "Name of admin to add")
	adminEmail = flag.String("admin_email", "", "Email of admin to add")
	adminPw    = flag.String("admin_pw", "", "Password of admin to add")
)

func DBConn() (*sqlx.DB, error) {
	return sqlx.Open("postgres",
		fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%d sslmode=disable",
			*dbUser, *dbName, *dbPw, *dbHost, *dbPort))
}

func start() error {
	db, err := DBConn()
	if err != nil {
		return err
	}

	store := postgrestore.MustCreate(db)
	store.MustPrepareStmts()
	sessionStore := auth.NewCookieSessionStore(
		[]byte("newauthenticatio"),
		[]byte("newencryptionkey"))

	um := auth.NewUserManager(nil, store, nil, sessionStore)
	m := models.NewManager(store)

	e := &env.Env{
		m,
		um,
		env.Config{
			Port: *port,
		},
	}

	router := httprouter.New()
	// Main React route
	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		f, _ := os.Open("html/index.html")
		io.Copy(w, f)
		f.Close()
	})

	// CSS, JS, etc
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	// Admin Routes
	router.GET("/admin/users", CreateHandlerWithEnv(e, handlers.CreateAdminUsersGETHandler))
	router.POST("/admin/user", CreateHandlerWithEnv(e, handlers.CreateAdminUsersPOSTHandler))
	router.DELETE("/admin/user/:user_id", CreateHandlerWithEnv(e, handlers.CreateAdminUserDELETEHandler))

	router.POST("/auth/login", CreateHandlerWithEnv(e, handlers.CreateLoginHandler))
	router.GET("/auth/logout", CreateHandlerWithEnv(e, handlers.CreateLogoutHandler))

	fmt.Println("Server started on port", e.Conf.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", e.Conf.Port), router)
}

func createSchema() error {
	db, err := DBConn()
	if err != nil {
		return err
	}
	store := postgrestore.MustCreate(db)

	store.MustCreateTypes()
	store.MustCreateTables()
	return nil
}

func dropSchema() error {
	db, err := DBConn()
	if err != nil {
		return err
	}

	store := postgrestore.MustCreate(db)

	store.MustDropTables()
	store.MustDropTypes()
	return nil
}

func addAdmin() error {
	// TODO: Parameter checking
	db, err := DBConn()
	if err != nil {
		return err
	}
	store := postgrestore.MustCreate(db)
	store.MustPrepareStmts()
	sessionStore := auth.NewCookieSessionStore(
		[]byte("new-authentication-key"),
		[]byte("new-encryption-key"))

	um := auth.NewUserManager(nil, store, nil, sessionStore)

	user, err := um.New(*adminName, *adminEmail, *adminPw, *adminPw, true, true)
	if err != nil {
		return err
	}

	err = um.Insert(user)
	if err != nil {
		return err
	}
	return nil
}

var actions = actionsMap{
	"start":         start,
	"create_schema": createSchema,
	"drop_schema":   dropSchema,
	"add_admin":     addAdmin,
}

func main() {
	flag.Parse()
	if !actions.validAction(*action) {
		fmt.Println("Please choose a valid action. Available: " + actions.available())
		os.Exit(1)
	}

	err := actions.perform(*action)
	if err != nil {
		fmt.Printf("Error performing %s: %v", *action, err)
	}
}

type InitHandler func(*env.Env, http.ResponseWriter, *http.Request, httprouter.Params) (http.Handler, int, error)

func CreateHandlerWithEnv(e *env.Env, ih InitHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		start := time.Now()
		h, status, err := ih(e, w, r, ps)
		fmt.Printf("HTTP %d: %v\n", status, err)
		if err != nil {
			switch status {
			case http.StatusNotFound:
				http.NotFound(w, r)
			case http.StatusInternalServerError:
				http.Error(w, http.StatusText(status), status)
			default:
				http.Error(w, http.StatusText(status), status)
			}
		}

		h.ServeHTTP(w, r)
		fmt.Printf("Request server in %v\n", time.Since(start))
	}
}
