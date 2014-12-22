/* AdminUserActions.js implements the admin actions associated with storing and
 * retrieving users.
 */

var AppDispatcher = require("../dispatcher/AppDispatcher");
var Constants = require("../constants/ExpenseTrackerConstants");
var AdminUserClient = require("../clients/UserClient").AdminUserClient;

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

function userDeleteSuccess(response) {
	dispatch(Constants.api.ADMIN_USER_DELETE_SUCCESS);
}

function userDeleteFail(response, params) {
	dispatch(Constants.api.ADMIN_USER_DELETE_FAIL);
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
		dispatch(Constants.api.ADMIN_USER_DELETE, userId);
		AdminUserClient.del(
			userDeleteSuccess,
			userDeleteFail);
	}
}

module.exports = AdminUserActions;