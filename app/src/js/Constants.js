import keyMirror from 'react/lib/keyMirror';

export default {
  // event name triggered from store, listened to by views
  CHANGE_EVENT: 'change',

  // Each time you add an action, add it here
  ActionTypes: keyMirror({
    ADD_ADMIN_USER: null,
    ADD_ADMIN_USER_SUCCESS: null,
    ADD_ADMIN_USER_FAIL: null,

    GET_ADMIN_USERS: null,
    GET_ADMIN_USERS_SUCCESS: null,
    GET_ADMIN_USERS_FAIL: null,

    EDIT_ADMIN_USER: null,
    EDIT_ADMIN_USER_SUCCESS: null,
    EDIT_ADMIN_USER_FAIL: null,

    DELETE_ADMIN_USER: null,
    DELETE_ADMIN_USER_SUCCESS: null,
    DELETE_ADMIN_USER_FAIL: null
  }),

  ActionSources: keyMirror({
    SERVER_ACTION: null,
    VIEW_ACTION: null
  })
};
