const request = require('superagent');
const DEFAULT_TIMEOUT = 30000;
module.exports = {
  createGetClientFunc(options) {
    if (!options.url) {
      throw new Error('URL must be provided when creating a client GET func');
    }
    options.timeout = options.timeout || DEFAULT_TIMEOUT;
    return () => {
      return new Promise((resolve, reject) => {
        request
          .get(options.url)
          .timeout(options.timeout)
          .end((err, res) => {
            if (err) {
              return reject(JSON.parse(res.text));
            }
            resolve(JSON.parse(res.text));
          })
      });
    }
  },
  createPostClientFunc(options) {
    if (!options.url) {
      throw new Error('URL must be provided when creating a POST client func');
    }
    options.timeout = options.timeout || DEFAULT_TIMEOUT;
    return (payload) => {
      return new Promise((resolve, reject) => {
        request
          .post(options.url)
          .send(payload)
          .timeout(options.timeout)
          .set('Content-Type: application/json')
          .end((err, res) => {
            if (err) {
              return reject(JSON.parse(res.text));
            }
            resolve(JSON.parse(res.text));
          });
      });
    }
  },
  createDeleteClientFunc(options) {
    if (!options.baseUrl) {
      throw new Error('baseUrl must be provided when creating a DELETE client func');
    }
    options.timeout = options.timeout || DEFAULT_TIMEOUT;
    return (objWithId) => {
      return new Promise((resolve, reject) => {
        request
          .del(options.baseUrl + '/' + objWithId.id)
          .timeout(options.timeout)
          .end((err, res) => {
            if (err) {
              var payload = JSON.parse(res.text);
              payload.objectToDelete = objWithId;
              return reject(payload);
            }
            resolve(objWithId);
          });
      });
    }
  },
  createPutClientFunc(options) {
    if (!options.baseUrl) {
      throw new Error('baseUrl myst be provided when creating a PUT client func');
    }
    options.timeout = options.timeout || DEFAULT_TIMEOUT;
    return (objWithId) => {
      return new Promise((resolve, reject) => {
        request
          .put(options.baseUrl + '/' + objWithId.id)
          .send(objWithId)
          .timeout(options.timeout)
          .end((err, res) => {
            if (err) {
              var payload = JSON.parse(text);
              payload.objectToPut = objWithId;
              return reject(payload);
            }
            resolve(objWithId);
          });
      });
    }
  }
};
