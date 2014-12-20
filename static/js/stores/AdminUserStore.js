var AppDispatcher = require("../dispatcher/AppDispatcher");
var Constants = require("../constants/ExpenseTrackerConstants");
var AdminUserClient = require("../clients/UserClients").AdminUserClient;

var users = {};

var AdminUserStore = {
	getUsers: function() {
		AdminUserClient.getAll().then(
			this.usersFromServerSuccess,
			this.usersFromServerFail);
		return users;
	}
};