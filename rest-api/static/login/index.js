function Login(conf) {
  this.conf = conf;
  var button = document.getElementById('login');
  button.onclick = this.handleLogin.bind(this);
  this.checkAuth();
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
    if (!res.ok) {
      console.log('Error: Login failed');
      return;
    }
    localStorage.setItem("accessToken", res.jwt);
    document.location.href = '/loggedIn';
  });
};

Login.prototype.checkAuth = function() {
  var accessToken = localStorage.getItem('accessToken');
  if (accessToken === null) {
    return;
  }
  var headers = new Headers();
  headers.append('Authorization', 'Bearer ' + accessToken);
  fetch(this.conf.backendUrl + '/loggedIn', {
    headers: headers
  }).
  then(function(res) {
    return res.json();
  }).
  then(function(res) {
    console.log(res);
    if (!res.loggedIn) {
      return;
    }
    document.location.href = '/loggedIn';
  });
};

(function() {
  'use strict';
  var login = new Login({
    backendUrl: '/api'
  });
})();
