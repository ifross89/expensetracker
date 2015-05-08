/*
 * client-utils.js contains utility functions for creating clients
 */

var DEFAULT_TIMEOUT = 30000;

function createGetClientFunc(options) {
  if (!options.url) {
    throw new Error('URL must be provided when creating a client func.');
  }

  options.timeout = options.timeout || DEFAULT_TIMEOUT;

  return function() {
    var promise = new Promise((resolve, reject) => {
      request
        .get(options.url)
        .timeout(options.timeout)
        .end((res) => {
          if (res.ok && res.body.status === 'success') {
            resolve();
          } else {
            console.log('GET request to ' + options.url + ' failed.');
            var resBody = JSON.parse(res.text);
            reject(resBody.message);
          }
        });
    });
    return promise;
  };
}


module.exports = {
  createGetClientFinc: createGetClientFunc
};
