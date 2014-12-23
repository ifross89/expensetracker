/* UserClient.js: Encapsulates interactions with the server when performing user
 * storage operations.
 */

var request = require("superagent");
var _ = require("underscore");
var $ = require("jquery");

var AppDispatcher = require("../dispatcher/AppDispatcher");

var API_URL = '/admin'
var TIMEOUT = 10000;

var  AdminUserClient = {
	getAll: function() {
		var promise = new Promise((resolve, reject) => {
			request
				.get(API_URL + "/users")
				.timeout(TIMEOUT)
				.end((res) => {
					if (res.ok) {
						// Server responded successfully, check to see if there is an error
						// in the JSON response.
						if (res.body.status === "success") {
							var users = res.text.data;
							_.each(users, (user) => user.pending = false);
							resolve(res.body.data);
						} else {
							reject(res);
						}
					} else {
						reject(res);
					}
				});
		});
		return promise;
	},

	save: function(user) {
		if (user.id) {
			promise = new Promise((resolve, reject) => {
				request.
					put(API_URL + "user/" + user.id)
					.timeout(TIMEOUT)
					.send(user)
					.end((res) => {
						if (res.ok && res.body.status === "success" {
							var user = res.body.data;
							resolve(user);
						} else {
							res.user = user;
							reject(res);
						})
					})
			});
			return promise;
		}

		var promise = new Promise((resolve, reject) => {
			request
				.post(API_URL + "/user")
				.timeout(TIMEOUT)
				.send(user)
				.end((res) => {
					if (res.ok && res.body.status == "success") {
						var user = res.body.data;
						resolve(user);
					} else {
						res.user = user;
						reject(res);
					}
				});
		});
		return promise
	},

	del: function(userId) {
		var promise = new Promise((resolve, reject) => {
			request
				.del(API_URL + "/user/" + userId)
				.timeout(TIMEOUT)
				.end(res => {
					if (res.ok && res.body.status === "success") {
						resolve({userId: userId});
					} else {
						reject(res, {userId: userId});
					}
				});
		});
		return promise;
	}


}

module.exports = {
	AdminUserClient: AdminUserClient
};