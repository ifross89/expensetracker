/* UserClient.js: Encapsulates interactions with the server when performing user
 * storage operations.
 */

var request = require("superagent");
var _ = require("underscore");
var $ = require("jquery");

var AppDispatcher = require("../dispatcher/app-dispatcher");

var API_URL = '/admin';
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
						var resBody = JSON.parse(res.text);
						reject(resBody);
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
						if (res.ok && res.body.status === "success") {
							var user = res.body.data;
							resolve(user);
						} else {
							var resBody = JSON.parse(res.text);
							resBody.user = user;
							reject(resBody);
						}
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
						var resBody = JSON.parse(res.text);
						resBody.user = user;
						reject(resBody);
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
						var resBody = JSON.parse(res.text);
						reject(resBody, {userId: userId});
					}
				});
		});
		return promise;
	}


};

module.exports = {
	AdminUserClient: AdminUserClient
};
