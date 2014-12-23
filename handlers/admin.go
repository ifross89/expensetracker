package handlers

import (
	"git.ianfross.com/expensetracker/env"

	"github.com/julienschmidt/httprouter"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type jsonResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"message,omitempty"`
}

func jsonSuccess(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(jsonResponse{"success", data, "", http.StatusOK})
}

func jsonError(w http.ResponseWriter, code int, message string) error {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(jsonResponse{"error", nil, message, code})
}

type HandlerVars struct {
	env *env.Env
	ps  httprouter.Params
}

func createHandlerVars(e *env.Env, ps httprouter.Params) *HandlerVars {
	return &HandlerVars{e, ps}
}

type adminUsersPOSTHandler struct {
	*HandlerVars
}

func (a adminUsersPOSTHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	u := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Admin    bool   `json:"admin"`
		Active   bool   `json:"active"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil && err != io.EOF {
		fmt.Printf("Error: %v\n", err)
		return
	}

	user, err := a.env.New(u.Name, u.Email, u.Password, u.Password, u.Active, u.Admin)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	err = a.env.Insert(user)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = jsonSuccess(w, user)
	if err != nil {
		fmt.Printf("Error encoding json: %v\n", err)
	}
}

func CreateAdminUsersPOSTHandler(
	e *env.Env,
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) (http.Handler, int, error) {
	return adminUsersPOSTHandler{createHandlerVars(e, ps)}, 200, nil
}

type adminUsersGETHandler struct {
	*HandlerVars
}

func CreateAdminUsersGETHandler(
	e *env.Env,
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) (http.Handler, int, error) {
	return adminUsersGETHandler{createHandlerVars(e, ps)}, 200, nil
}

func (a adminUsersGETHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	users, err := a.env.Users()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Users:", users)

	err = jsonSuccess(w, users)
	if err != nil {
		fmt.Printf("Error encoding json: %v\n", err)
	}
}

type adminUserDELETEHandler struct {
	*HandlerVars
}

func CreateAdminUserDELETEHandler(
	e *env.Env,
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) (http.Handler, int, error) {
	return adminUserDELETEHandler{createHandlerVars(e, ps)}, http.StatusOK, nil
}

func (h adminUserDELETEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uidStr := h.ps.ByName("user_id")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.env.DeleteUserById(int64(uid))
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
	}

	err = jsonSuccess(w, nil)
	if err != nil {
		fmt.Printf("Error sending success: %v", err)
	}
}
