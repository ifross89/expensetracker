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

func (a actionsMap) perform(action string) {
	err := a[action]()
	if err != nil {
		panic(err)
	}
}

var (
	dbUser = flag.String("db_user", "expensetracker", "database user to connect with")
	dbName = flag.String("db_name", "expensetracker", "name of the database to connect to")
	dbPw   = flag.String("db_pw", "", "user's database password")
	port   = flag.Int("port", 8181, "HTTP port to listen on")
	action = flag.String("action", "start", "action to perform. Available: "+actions.available())
)

func DBConn() *sqlx.DB {
	return sqlx.MustOpen("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", dbUser, dbName, dbPw))
}

func start() error {
	db := DBConn()
	store := postgrestore.MustCreate(db)
	sessionStore := auth.NewCookieSessionStore(
		[]byte("new-authentication-key"),
		[]byte("new-encryption-key"))

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

	fmt.Println("Server started on port", e.Conf.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", e.Conf.Port), router)
}

func createSchema() error {
	db := DBConn()
	store := postgresstore.MustCreate(db)

	store.MustCreateTypes()
	store.MustCreateTables()
	return nil
}

func dropSchema() error {
	db := DBConn()
	store := postgresstore.MustCreate(db)

	store.MustDropTables()
	store.MustDropTypes()
	return nil
}

func addAdmin() error {
	fmt.Println("Adding admin")
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

	actions.perform(*action)
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
