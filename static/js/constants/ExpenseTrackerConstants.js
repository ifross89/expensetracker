var keyMirror = require('keymirror');

var consts = {
	

}

module.exports = {
	api: keyMirror({
		ADMIN_USERS_LOAD: null,
		ADMIN_USERS_LOAD_SUCCESS: null,
		ADMIN_USERS_LOAD_FAIL: null,

		ADMIN_USER_SAVE: null,
		ADMIN_USER_SAVE_SUCCESS: null,
		ADMIN_USER_SAVE_FAIL: null,

		ADMIN_USER_DELETE: null,
		ADMIN_USER_DELETE_SUCCESS: null,
		ADMIN_USER_DELETE_FAIL: null
	}),

	request: keyMirror({
		TIMEOUT: null,
		ERROR: null,
		PENDING: null
	}),

	store: keyMirror({
		CHANGE_EVENT: null
	})
}
