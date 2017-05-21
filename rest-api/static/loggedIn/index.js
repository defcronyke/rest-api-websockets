function LoggedIn(conf) {
  this.conf = conf;
  var helloButton = document.getElementById('wsHello');
  helloButton.onclick = this.handleSendHello.bind(this);
  var notImplementedButton = document.getElementById('wsNotImplemented');
  notImplementedButton.onclick = this.handleSendNotImplemented.bind(this);
  var logoutButton = document.getElementById('logout');
  logoutButton.onclick = this.handleLogout.bind(this);
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
  this.ws = new WebSocket(this.conf.wsUrl);
  if (this.ws === null) {
    console.log('Error: Websocket connection failed');
    return;
  }
  console.log('Connected to websocket server');
  this.ws.onmessage = function(e) {
    console.log('websocket msg received:', JSON.parse(e.data));
  }
};

LoggedIn.prototype.handleSendHello = function() {
  this.ws.send(JSON.stringify({
    cmd: 'hello'
  }));
};

LoggedIn.prototype.handleSendNotImplemented = function() {
  this.ws.send(JSON.stringify({
    cmd: 'garlic'
  }));
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
