const utils = require('./ClientUtils');

export default {
  getAll: utils.createGetClientFunc({
    url: '/admin/users'
  }),

  save: utils.createPostClientFunc({
    url: '/admin/user'
  }),

  del: utils.createDeleteClientFunc({
    baseUrl: '/admin/user'
  }),

  edit: utils.createPutClientFunc({
    baseUrl: '/admin/user'
  })
};

