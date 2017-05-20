function LoggedIn(conf) {
  this.conf = conf;
  var button = document.getElementById('logout');
  button.onclick = this.handleLogout.bind(this);
  this.checkAuth();
  console.log("logged in area");
  this.connectWs();
}

LoggedIn.prototype.checkAuth = function() {
  var accessToken = localStorage.getItem('accessToken');
  if (accessToken === null) {
    document.location.href = '/login';
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
      document.location.href = '/login';
      return;
    }
  });
};

LoggedIn.prototype.connectWs = function() {
  var ws = new WebSocket(this.conf.wsUrl);
};

LoggedIn.prototype.handleLogout = function() {
  console.log('Logging out...');
};

(function() {
  'use strict';
  var loggedIn = new LoggedIn({
    backendUrl: '/api',
    wsUrl: 'ws://localhost:8081'
  });
})();
