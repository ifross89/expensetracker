var AppDispatcher = require("../dispatcher/app-dispatcher");
var Constants = require("../constants/expense-tracker-constants");
var AuthClient = require('../clients/auth-client').AuthClient;

function dispatch(key, response, params) {
	var payload = {actionType: key, response: response};
	if (params) {
		payload.queryParams = params;
	}
	AppDispatcher.handleRequestAction(payload);
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

function changePasswordSuccess() {
	dispatch(Constants.api.CHANGE_PASSWORD_SUCCESS);
	console.log('Change password success');
}

function changePasswordFail(message) {
	dispatch(Constants.api.CHANGE_PASSWORD_FAIL, message);
	console.log('Change password fail. Message:', message);
}

var AuthActions = {
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
	},

	changePassword: function(
		oldPassword, newPassword, confirmPassword
	) {
		dispatch(Constants.api.CHANGE_PASSWORD);
		AuthClient.changePassword(oldPassword, newPassword, confirmPassword)
			.then(changePasswordSuccess, changePasswordFail);
	}
};

module.exports = AuthActions;
