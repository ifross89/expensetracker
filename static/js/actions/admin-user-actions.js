/* AdminUserActions.js implements the admin actions associated with storing and
 * retrieving users.
 */

var AppDispatcher = require("../dispatcher/app-dispatcher");
var Constants = require("../constants/expense-tracker-constants");
var AdminUserClient = require("../clients/user-client").AdminUserClient;
var AuthClient = require('../clients/auth-client').AuthClient;

function dispatch(key, response, params) {
	var payload = {actionType: key, response: response};
	if (params) {
		payload.queryParams = params;
	}
	AppDispatcher.handleRequestAction(payload);
}

function usersFromServerSuccess(users) {
	dispatch(Constants.api.ADMIN_USERS_LOAD_SUCCESS, users);
}

function usersFromServerFail(result) {
	dispatch(Constants.api.ADMIN_USERS_LOAD_FAIL);
}

function userSaveSuccess(user) {
	dispatch(Constants.api.ADMIN_USER_SAVE_SUCCESS, user);
}

function userSaveFail(response) {
	dispatch(Constants.api.ADMIN_USER_SAVE_FAIL, response);
}

function userDeleteSuccess(response, params) {
	dispatch(Constants.api.ADMIN_USER_DELETE_SUCCESS, response);
}

function userDeleteFail(response, params) {
	dispatch(Constants.api.ADMIN_USER_DELETE_FAIL, response, params);
}

function loginSuccess(user) {
	dispatch(Constants.api.LOGIN_SUCCESS, user);
}

function loginFail(message) {
	dispatch(Constants.api.LOGIN_FAIL, message);
}

function logoutSuccess() {
	dispatch(Constants.api.LOGOUT_SUCCESS);
	console.log('logout success');
}

function logoutFail(message) {
	dispatch(Constants.api.LOGOUT_FAIL, message);
	console.log('logout fail');
}

var AdminUserActions = {
	getUsers: function() {
		dispatch(Constants.api.ADMIN_USERS_LOAD);
		AdminUserClient.getAll().then(
			usersFromServerSuccess,
			usersFromServerFail);
	},

	saveUser: function(user) {
		dispatch(Constants.api.ADMIN_USER_SAVE, user);
		AdminUserClient.save(user).then(
			userSaveSuccess,
			userSaveFail);
	},

	deleteUser: function(userId) {
		console.log("deleteUser");
		console.log(userId);
		dispatch(Constants.api.ADMIN_USER_DELETE, {userId: userId});
		AdminUserClient.del(userId).then(
			userDeleteSuccess,
			userDeleteFail);
	},

	login: function(email, password) {
		console.log('login');
		dispatch(Constants.api.LOGIN);
		AuthClient.login(email, password)
			.then(loginSuccess, loginFail);
	},

	logout: function() {
		console.log('logout');
		dispatch(Constants.api.LOGOUT);
		AuthClient.logout()
			.then(logoutSuccess, logoutFail);
	}
};

module.exports = AdminUserActions;