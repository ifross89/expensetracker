/* AdminUserActions.js implements the admin actions associated with storing and
 * retrieving users.
 */

var AppDispatcher = require("../dispatcher/AppDispatcher");
var Constants = require("../constants/ExpenseTrackerConstants");
var AdminUserClient = require("../clients/UserClients").AdminUserClient;

var AdminUserActions = {
	LOADING_TOKEN: "LOADING";
	getUsers: function() {
		AdminUserClient.getAll().then(
			this.usersFromServerSuccess,
			this.usersFromServerFail);

	}
};