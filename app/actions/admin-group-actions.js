/*
 * admin-group-actions.js implements the admin actions associated with storing and retrieving
 * groups
 */

var AppDispatcher = require('../dispatcher/app-dispatcher');
var Constants = require('../constants/expense-tracker-constants');
var AdminGroupClient = require('../clients/admin-group-client');

function dispatch(key, response, params) {
  var payload = { actionType: key, response: response };

  if (params) {
    payload.queryParams = params;
  }

  AppDispatcher.handleRequestAction(payload);
}

function groupsFromServerSuccess(groups) {
  dispatch(Constants.api.ADMIN_GROUP_LOAD_SUCCESS, groups);
}

function groupsFromServerFail(result) {
  console.log('groupsFromServerFail: result=' + result);
  dispatch(Constants.api.ADMIN_USERS_LOAD_FAIL, result);
}

module.exports = {
  getGroups: function() {
    dispatch(Constants.api.ADMIN_GROUP_LOAD);
    AdminGroupClient.getAll().then(
      groupsFromServerSuccess,
      groupsFromServerFail
    );
  }
};
