package handlers

import (
	"git.ianfross.com/ifross/expensetracker/env"

	"github.com/julienschmidt/httprouter"

	"encoding/json"
	"net/http"
	"fmt"
)

const (
	ErrInvalidUsernamePw = "Invalid username or password supplied"
)

type loginHandler struct {
	*HandlerVars
}

func CreateLoginHandler(
	e *env.Env,
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) (http.Handler, int, error) {
	return loginHandler{createHandlerVars(e, ps)}, 200, nil
}

type loginInfo struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (h loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	info := loginInfo{}
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		fmt.Println(err)
		jsonError(w, http.StatusBadRequest, "Username and password must be supplied")
		return
	}

	u, err := h.env.UserManager.ByEmail(info.Email)
	if err != nil {
		fmt.Println(err)
		jsonError(w, http.StatusUnauthorized, ErrInvalidUsernamePw)
		return
	}

	err = h.env.UserManager.Authenticate(u, info.Password)
	if err != nil {
		fmt.Println(err)
		jsonError(w, http.StatusUnauthorized, ErrInvalidUsernamePw)
		return
	}

	err = h.env.UserManager.LogIn(w, r, u)
	if err != nil {
		fmt.Println(err)
		jsonErrorWithCodeText(w, http.StatusInternalServerError)
		return
	}

	jsonSuccess(w, u)
}

type logoutHandler struct {
	*HandlerVars
}

func CreateLogoutHandler(
	e *env.Env,
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) (http.Handler, int, error) {
	return logoutHandler{createHandlerVars(e, ps)}, 200, nil
}

func (h logoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := h.env.UserManager.FromSession(w, r)
	if err != nil {
		jsonErrorWithCodeText(w, http.StatusUnauthorized)
		return
	}

	err = h.env.UserManager.LogOut(w, r)
	if err != nil {
		jsonErrorWithCodeText(w, http.StatusInternalServerError)
		return
	}

	jsonSuccess(w, nil)
}
