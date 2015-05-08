package handlers

import (
	"git.ianfross.com/ifross/expensetracker/env"
	"git.ianfross.com/ifross/expensetracker/auth"

	"github.com/julienschmidt/httprouter"
	"github.com/golang/glog"
	"github.com/juju/errors"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"git.ianfross.com/ifross/expensetracker/models"
)

type jsonResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
}

func jsonSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(jsonResponse{"success", data, "", http.StatusOK})
	if err != nil {
		glog.Errorf("Error encoding json in successful respone:%v", err)
	}
}

func jsonError(w http.ResponseWriter, code int, message string, err error) error {
	if err != nil {
		glog.Errorf("Error in handler: error=%v\nmessage=%s", errors.ErrorStack(err), message)
	} else {
		glog.Errorf("Error in handler: message=%s", message)
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(jsonResponse{"error", nil, message, code})
}

func jsonErrorWithCodeText(w http.ResponseWriter, code int, err error) error {
	return jsonError(w, code, http.StatusText(code), err)
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
		jsonErrorWithCodeText(w, http.StatusInternalServerError, errors.Trace(err))
		return
	}

	user, err := a.env.New(u.Name, u.Email, u.Password, u.Password, u.Active, u.Admin)
	if err != nil {
		err = errors.Trace(err)
		jsonError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = a.env.Insert(user)
	if err != nil {
		jsonErrorWithCodeText(w, http.StatusInternalServerError, errors.Trace(err))
		return
	}

	jsonSuccess(w, user)
	if err != nil {
		glog.Errorf("Error encoding json: %v\n", err)
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
		jsonErrorWithCodeText(w, http.StatusInternalServerError, errors.Trace(err))
		return
	}

	fmt.Println("Users:", users)

	jsonSuccess(w, users)
	if err != nil {
		glog.Errorf("Error encoding json: %v\n", err)
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
		jsonError(w, http.StatusInternalServerError, err.Error(), errors.Trace(err))
		return
	}

	err = h.env.DeleteUserById(int64(uid))
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error(), errors.Trace(err))
		return
	}

	jsonSuccess(w, nil)
	if err != nil {
		fmt.Printf("Error sending success: %v", err)
	}
}


type adminGroupsGETHandler struct {
	*HandlerVars
}

func CreateAdminGroupsGETHandler(
	e *env.Env,
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) (http.Handler, int, error) {
	return adminGroupsGETHandler{createHandlerVars(e, ps)}, http.StatusOK, nil
}

func (h adminGroupsGETHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get user
	_, err := h.env.AdminFromSession(w, r)
	if err != nil {
		jsonErrorWithCodeText(w, http.StatusUnauthorized, errors.Trace(err))
		return
	}

	// User is authenticated
	groups, err := h.env.AllGroups()
	if err != nil {
		jsonErrorWithCodeText(w, http.StatusInternalServerError, errors.Trace(err))
		return
	}

	jsonSuccess(w, groups)
}

type adminGroupPOSTHandler struct {
	*HandlerVars
}

func CreateAdminGroupPOSTHandler(
	e *env.Env,
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) (http.Handler, int, error) {
	return adminGroupPOSTHandler{createHandlerVars(e, ps)}, http.StatusOK, nil
}

func (h adminGroupPOSTHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := h.env.AdminFromSession(w, r)

	if err != nil {
		jsonErrorWithCodeText(w, http.StatusUnauthorized, errors.Trace(err))
		return
	}

	newGroup := struct {
		Name string `json:"name"`
		Emails []string `json:"emails"`
		}{}

	err = json.NewDecoder(r.Body).Decode(&newGroup)

	if err != nil && err != io.EOF {
		jsonErrorWithCodeText(w, http.StatusInternalServerError, errors.Trace(err))
		return
	}


	var users []*auth.User
	// First make sure all of them are users
	for _, email := range newGroup.Emails {
		u, err := h.env.UserManager.ByEmail(email)
		if err != nil {
			jsonError(w, http.StatusBadRequest, fmt.Sprintf("User with email %s does not exist", email), errors.Trace(err))
			return
		}
		users = append(users, u)
	}

	g, err := h.env.NewGroup(newGroup.Name)
	if err != nil {
		jsonErrorWithCodeText(w, http.StatusInternalServerError, errors.Trace(err))
		return
	}

	for _, user := range users {
		err := h.env.AddUserToGroup(g, user, false)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "Error creating group", errors.Trace(err))
			return
		}
	}

	jsonSuccess(w, g)
}

type adminGroupDELETEHandler struct {
	*HandlerVars
}

func CreateAdminGroupDELETEHandler(
	e *env.Env,
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) (http.Handler, int, err) {
	return adminGroupDELETEHandler{createHandlerVars(e, ps)}
}

func (h adminGroupDELETEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := h.env.AdminFromSession(w, r)

	if err != nil {
		jsonErrorWithCodeText(w, http.StatusUnauthorized, errors.Trace(err))
		return
	}

	// User is admin

	groupId := struct {
		Id int64 `json:"id"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&groupId)
	if err != nil {
		jsonErrorWithCodeText(w, http.StatusInternalServerError, errors.Trace(err))
		return
	}

	err = h.env.DeleteGroup(&models.Group{Id:groupId.Id})
	if err != nil {
		jsonErrorWithCodeText(w, http.StatusInternalServerError, errors.Trace(err))
		return
	}

	jsonSuccess(w, nil)
}

type adminGroupPUTHandler struct {
	*HandlerVars
}

func CreateAdminGroupPUTHandler(
	e *env.Env,
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) (http.Handler, int, error) {
	return adminGroupPUTHandler{createHandlerVars(e, ps)}, http.StatusOK, nil
}

// Need to figure out what happens to the expenses
func (h adminGroupPUTHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := h.env.AdminFromSession(w, r)

	if err != nil {
		jsonErrorWithCodeText(w, http.StatusUnauthorized, errors.Trace(err))
		return
	}

	group := struct {
		Id int64 `json:"id"`
		Name string `json:"name"`
		Emails []string `json:"emails"`
		}{}

	err = json.NewDecoder(r.Body).Decode(&group)

	if err != nil && err != io.EOF {
		jsonErrorWithCodeText(w, http.StatusInternalServerError, errors.Trace(err))
		return
	}

	// TODO: finish this function

	jsonError(w, http.StatusServiceUnavailable, "Unimplemented", nil)
}
