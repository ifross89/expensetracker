'use strict';

import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import BaseStore from './BaseStore';
import assign from 'object-assign';

// data storage
let _data = [];

// add private functions to modify data
function add(user) {
  _data.push(user);
}

function updateUser(user) {
  for (let i=0; i<_data.length; i++ ) {
    if (_data[i].id == user.id) {
      _data[i] = user;
    }
  }
}

// Facebook style store creation.
const AdminUserStore = assign({}, BaseStore, {
  // public methods used by Controller-View to operate on data
  getAll() {
    return {
      users: _data
    };
  },

  // register store with dispatcher, allowing actions to flow through
  dispatcherIndex: Dispatcher.register(function(payload) {
    let action = payload.action;

    console.log(payload);
    switch(action.type) {
      case Constants.ActionTypes.ADD_ADMIN_USER_SUCCESS:
        let user = action.payload;
        add(user);
        AdminUserStore.emitChange();
        break;
      case Constants.ActionTypes.GET_ADMIN_USERS_SUCCESS:
        _data = [];
        let users = action.payload.data;
        users.forEach(user => {
          add(user);
        });
        AdminUserStore.emitChange();
        break;
      case Constants.ActionTypes.EDIT_ADMIN_USER_SUCCESS:
        updateUser(action.payload);
        AdminUserStore.emitChange();
        break;
    }
  })
});

export default AdminUserStore;
