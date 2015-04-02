var AppDispatcher = require("../dispatcher/app-dispatcher");
var Constants = require("../constants/expense-tracker-constants");
var _ = require('underscore');
var EventEmitter = require('events').EventEmitter;


var users = {};
var pendingUsers = {};
var pendingDelete = {};
var j = 0;

var AdminUserStore = _.extend({}, EventEmitter.prototype,{
	getUsers: function() {
		return _.extend({}, users, pendingUsers);
	},

	emitChange: function() {
	this.emit(Constants.store.CHANGE_EVENT);
},

	/** * @param {function} callback */
	addChangeListener: function(callback){
		this.on(Constants.store.CHANGE_EVENT, callback);
	},

	/** * @param {function} callback */
	removeChangeListener: function(callback) {
		this.removeListener(Constants.store.CHANGE_EVENT, callback);
	}
});

function makeKey(userId) {
	return "user-" + userId;
}

function createKey(user) {
	var key = "";
	if (user.id) {
		key = makeKey(user.id);
	} else {
		j++;
		key = "pending-user-" + j;
	}
	user.key = key;
}

function persistUsers(userList) {
	for (var i=0; i < userList.length; i++) {
		var user = userList[i];
		createKey(user);
		users[user.key] = user;
	}
}

function saveNewUser(user) {
	createKey(user);
	pendingUsers[user.email] = user;
}

function saveExistingUser(user) {
	delete pendingUsers[user.email];
	createKey(user);
	users[user.key] = user;
}

function deleteUser(userId) {
	// Get the user from the store
	var key = makeKey(userId);
	var user = users[key];
	if (!user) {
		console.error("Could not get user with id=" + userId);
		return
	}
	pendingDelete[key] = user;
	delete users[key];
}

function purgeUser(userId) {
	delete pendingDelete[makeKey(userId)];
}

function restoreUser(userId) {
	var key = makeKey(userId);
	var user = pendingDelete[key];
	delete pendingDelete[key];
	users[key] = user;
}

AdminUserStore.appDispatch = AppDispatcher.register(function(payload) {
	var action = payload.action;
	console.log("ACTION: ", payload.action);
	switch(action.actionType) {
		case Constants.api.ADMIN_USERS_LOAD_SUCCESS:
			persistUsers(action.response);
			break;
		case Constants.api.ADMIN_USER_SAVE:
			var user = action.response;
			user.pending = true;
			if (user.id) {
				saveExistingUser(user);
			} else {
				saveNewUser(user);
			}
			break;

		case Constants.api.ADMIN_USER_SAVE_SUCCESS:
			var user = action.response;
			user.pending = false;
			saveExistingUser(user);
			break;

		case Constants.api.ADMIN_USER_DELETE:
			deleteUser(action.response.userId);
			break;

		case Constants.api.ADMIN_USER_DELETE_SUCCESS:
			purgeUser(action.response.userId);
			break;

		case Constants.api.ADMIN_USER_DELETE_FAIL:
			restoreUser(action.response.userId);
			break;

		default:
			return true;
	}
	AdminUserStore.emitChange();
	return true;
});

module.exports = AdminUserStore;