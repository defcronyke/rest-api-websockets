function Login(conf) {
  this.conf = conf;
  var button = document.getElementById('login');
  button.onclick = this.handleLogin.bind(this);
}

Login.prototype.handleLogin = function() {
  var body = {
    username: document.getElementById('username').value,
    password: document.getElementById('password').value
  };
  this.postLogin(body);
};

Login.prototype.postLogin = function(body) {
  fetch(this.conf.backendUrl + '/login', {
    method: 'POST',
    body: JSON.stringify(body)
  }).
  then(function(res) {
    return res.json();
  }).
  then(function(res) {
    console.log(res);
  });
};

(function() {
  'use strict';
  var login = new Login({
    backendUrl: '/api'
  });
})();
