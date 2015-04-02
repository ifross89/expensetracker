var AppDispatcher = require('../dispatcher/app-dispatcher');
var EventEmitter = require('events').EventEmitter;
var Constants = require('../constants/ExpenseTracjerConstants');
var assign = require('object-assign');

var CHANGE_EVENT = 'change';

var _users = {};

var _pendingIndex = 1;

function create(user) {
	if (user.id) {
		// Can generate the key
		user.key = "user-" + id;
	} else {
		// create pending key
		user.key = "pending-user-" + _pendingIndex;
		_pendingIndex++;
	}
	_users[key] = assign({}, user);
}

function update(key, updates) {
	var old = _users[key];
	if (updates.id) {
		// No longer pending so must change key.
		delete _users[key]; // fails silently if doesn't exist
		key = "user-" + id;
	}
	_users[key] = assign({}, old, updates);
}

function updateAll(updates) {
	for (var key in _users) {
		update(key, updates);
	}
}

function destroy(key) {
	delete _users[key];
}

var UserStore = assign({}, EventEmitter.prototype, {
	getAll: function() {
		return _users;
	},

	emitChange: function() {
		this.emit(CHANGE_EVENT);
	},

	addChangeListener: function(callback) {
		this.on(CHANGE_EVENT, callback);
	},

	removeChangeListener: function() {
		this.removeListener(CHANGE_EVENT, callback);
	}
});

AppDispatcher.register(function(payload) {
	var action = payload.action;
	var user;

	switch(action.actionType) {
		case Constants.USER_CREATE:
			user = action.user;
			create(user);
			break;

		case Constants.USER_DESTROY:
			destroy(action.key);
			break;

		case Constants.MARK_PERSISTED:
			if (action.user.id) {
				update(action.key, action.user);
			}
			break;
		case Constants.USER_UPDATE:
			update(action.key, action.user);
			break;
		default:
			return true;
	}
	UserStore.emitChange();

	return true;
});