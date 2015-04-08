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

type changePasswordInfo struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

type changePasswordHandler struct {
	*HandlerVars
}

func CreateChangePasswordHandler(
	e *env.Env,
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) (http.Handler, int, error) {
	return changePasswordHandler{createHandlerVars(e, ps)}, http.StatusOK, nil
}

func (h changePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, err := h.env.UserManager.FromSession(w, r)
	if err != nil {
		jsonErrorWithCodeText(w, http.StatusUnauthorized)
		return
	}
	info := changePasswordInfo{}
	err = json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		fmt.Println(err)
		jsonError(w, http.StatusBadRequest, "Old password, new password and password confirmation must be supplied")
		return
	}

	err = h.env.UserManager.Authenticate(u, info.OldPassword)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "Old password supplied was incorrect")
		return
	}

	err = h.env.UserManager.UpdatePw(u, info.NewPassword, info.ConfirmPassword)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonSuccess(w, nil)

}
