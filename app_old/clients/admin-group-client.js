/* admin-group-client.js: Encapsulates interactions with the server when performing
 * group storage actions.
 */

var request = require('superagent');

var Utils = require('./client-utils');

var API_URL = '/admin';

var AdminGroupClient = {
  getAll: Utils.createGetClientFunc({url: API_URL + '/groups'})
};

module.exports = AdminGroupClient;

