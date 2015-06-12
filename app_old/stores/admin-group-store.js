var _ = require('underscore');

var AppDispatcher = require('../dispatcher/app-dispatcher');
var Constants = require('../constants/expense-tracker-constants');
var EventEmitter = require('events').EventEmitter;

var groups = {};
var pendingGroups = {};
var pendingDelete = {};

var j = 0;

var AdminGroupStore = _.extend({}, EventEmitter.prototype, {
  getGroups: function() {
    return _.extend({}, groups, pendingGroups);
  },

  emitChange: function() {
    this.emit(Constants.store.ADMIN_GROUP_CHANGE_EVENT);
  },

  addChangeListener: function(callback) {
    this.on(Constants.store.ADMIN_GROUP_CHANGE_EVENT, callback);
  },

  removeChangeListener: function(callback) {
    this.removeListener(Constants.store.ADMIN_GROUP_CHANGE_EVENT);
  }
});

function makeKey(groupId) {
  return 'group-' + groupId;
}

function createKey(group) {
  var key = '';
  if (group.id) {
    key = makeKey(group.id);
  } else {
    j++;
    key = 'pending-group-' + j;
  }

  group.key = key;
}

function persistSavedGroups(groups) {
  groups.forEach( function(group) {
    createKey(group);
    groups[group.key] = user;
  });
}

function saveNewGroup(group) {
  createKey(group);
  pendingGroups[group.name] = group;
}

function saveExistingGroup(group) {
  delete pendingGroups[group.name];
  persistSavedGroups([group]);
}

function purgeGroup(group) {
  delete pendingDelete[group.key];
  delete pendingGroups[group.name];
  delete groups[group.key];
}

function restoreGroup(groupId) {
  var key = makeKey(groupId);
  var group = pendingDelete[key];
  delete pendingDelete[key];
  groups[key] = group;
}


function deletePending(name) {
  delete pendingGroups[name];
}

AdminGroupStore.appDispatch = AppDispatcher.register(function(payload) {
  var action = payload.action;
  console.log('ADMIN GROUP STORE: Action=', payload.action);
  switch (action.actionType) {
    case Constants.api.ADMIN_GROUP_LOAD_SUCCESS:
          persistGroups(action.response);
          break;
    default:
          return true;

  }

  AdminGroupStore.emitChange();
  return true;
});

module.exports = AdminGroupStore;
