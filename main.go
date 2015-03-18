package main

import (
	"git.ianfross.com/ifross/expensetracker/auth"
	"git.ianfross.com/ifross/expensetracker/env"
	"git.ianfross.com/ifross/expensetracker/handlers"
	"git.ianfross.com/ifross/expensetracker/models"
	"git.ianfross.com/ifross/expensetracker/models/postgrestore"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"

	"fmt"
	"io"
	"net/http"
	"os"
	"flag"
	"time"
)

var (
	dbUser = flag.String("db-user", "expensetracker", "database user to connect with")
	dbName = flag.String("db-name", "expensetracker", "name of the database to connect to")
	dbPw = flag.String("db-pw", "", "user's database password")
	port = flag.Int("port", 8181, "HTTP port to listen on")

)

func main() {
	db := sqlx.MustOpen("postgres", "user=ian dbname=expensetrackerv2 password=wedge89 sslmode=disable")
	store := postgrestore.MustCreate(db)
	sessionStore := auth.NewCookieSessionStore(
		[]byte("new-authentication-key"),
		[]byte("new-encryption-key"))

	um := auth.NewUserManager(nil, store, nil, sessionStore)
	m := models.NewManager(store)
	//store.MustCreateTypes()
	//store.MustCreateTables()
	store.MustPrepareStmts()

	e := &env.Env{
		m,
		um,
		env.Config{
			Port: 8182,
		},
	}

	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		f, _ := os.Open("html/index.html")
		io.Copy(w, f)
		f.Close()
	})

	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.GET("/admin/users", CreateHandlerWithEnv(e, handlers.CreateAdminUsersGETHandler))
	router.POST("/admin/user", CreateHandlerWithEnv(e, handlers.CreateAdminUsersPOSTHandler))
	router.DELETE("/admin/user/:user_id", CreateHandlerWithEnv(e, handlers.CreateAdminUserDELETEHandler))

	fmt.Println("Server started on port", e.Conf.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", e.Conf.Port), router)
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

type TestHandler struct {
	e  *env.Env
	ps httprouter.Params
}

func CreateTestHandler(
	e *env.Env,
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) (http.Handler, int, error) {
	return &TestHandler{e, ps}, 0, nil
}

func (h *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s", h.ps.ByName("name"))
}
