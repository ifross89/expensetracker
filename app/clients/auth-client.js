var request = require('superagent');

var API_URL = '/auth';
var TIMEOUT = 10000;

var AuthClient = {
  login: (email, password) => {
    var promise = new Promise((resolve, reject) => {
      request
        .post(API_URL + '/login')
        .timeout(TIMEOUT)
        .set('Content-Type: application/json')
        .send({email: email, password: password})
        .end((res) => {
          if (res.ok && res.body.status === 'success') {
            var user = res.body.data;
            resolve(user);
          } else {
            var resBody = JSON.parse(res.text);
            reject(resBody.message);
          }
        });
    });
    return promise;
  },
  logout: () => {
    var promise = new Promise((resolve, reject) => {
      request
        .get(API_URL + '/logout')
        .timeout(TIMEOUT)
        .end((res) => {
          if (res.ok && res.body.status === 'success') {
            resolve();
          } else {
            console.log('logout fail, res=', res);
            var resBody = JSON.parse(res.text);
            reject(resBody.message);
          }
        })
    });
    return promise;
  },

  changePassword: (oldPassword, newPassword, confirmPassword) => {
    var promise = new Promise((resolve, reject) => {
      request
        .post(API_URL + '/change_password')
        .timeout(TIMEOUT)
        .set('Content-Type: application/json')
        .send({
          oldPassword: oldPassword,
          newPassword: newPassword,
          confirmPassword: confirmPassword
        })
        .end((res) => {
          if (res.ok && res.body.status === 'success') {
            resolve();
          } else {
            console.log('changePassword fail, res=', res);
            var resBody = JSON.parse(res.text);
            reject(resBody.message);
          }
        })
    });
    return promise;
  }

};

module.exports = {
  AuthClient: AuthClient
};
