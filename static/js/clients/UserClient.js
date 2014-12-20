/* UserClient.js: Encapsulates interactions with the server when performing user
 * storage operations.
 */

var request = require("superagent");
var _ = require("underscore");

var AppDispatcher = require("../dispatcher/AppDispatcher");

var API_URL = '/admin/users'

var = AdminUserClient {
	getAll: function() {
		var promise = new Promise((resolve, reject) => {
			request
				.get(API_URL)
				.end((res) => {
					if (res.ok) {
						// Server responded successfully, check to see if there is an error
						// in the JSON response.
						if (res.body.status === "success") {
							var users = res.body.data;
							_.each(users, (user) => user.pending = false;)
							resolve(res.body.data);
						} else {
							reject(res.body.message);
						}
					} else {
						reject(res.body.message);
					}
				});
		});
		return promise;
	},

	save: function(user) {
		user.pending = true;
		var promise = new Promise((resolve, reject) => {
			request
				.post(API_URL)
				.send(user)
				.end((res) => {
					if (res.ok && res.body.status == "success") {
						var user = res.body.data;
						user.pending = false;
						resolve(user);
					} else {
						reject(res.body);
					}
				});
		});
	}
}

module.exports = AdminUserClient;