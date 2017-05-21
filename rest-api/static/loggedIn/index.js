function LoggedIn(conf) {
  this.conf = conf;

  var helloButton = document.getElementById('wsHello');
  helloButton.onclick = this.handleSendHello.bind(this);

  var notImplementedButton = document.getElementById('wsNotImplemented');
  notImplementedButton.onclick = this.handleSendNotImplemented.bind(this);

  var wsBroadcastButton = document.getElementById('wsBroadcast');
  wsBroadcastButton.onclick = this.handleBroadcast.bind(this);

  var logoutButton = document.getElementById('logout');
  logoutButton.onclick = this.handleLogout.bind(this);

  this.msgBox = document.getElementById('msgBox');

  this.checkAuth();
  console.log("logged in area");
  this.connectWs(this.conf.wsUrl);
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

LoggedIn.prototype.connectWs = function(url) {
  this.ws = new WebSocket(url);
  if (this.ws === null) {
    console.log('Error: Websocket connection failed');
    return;
  }
  this.ws.onopen = this.onopen.bind(this);
  this.ws.onmessage = this.onmessage.bind(this);
  this.ws.onclose = this.onclose.bind(this);
  this.ws.onerror = this.onerror.bind(this);
};

LoggedIn.prototype.onopen = function() {
  var connectedMsg = 'Connected to websocket server';
  console.log(connectedMsg);
  this.msgBox.value += connectedMsg + '\n';
  this.msgBox.scrollTop = this.msgBox.scrollHeight;
};

LoggedIn.prototype.onclose = function() {
  var disconnectMsg = 'Connection to websocket server closed, will attempt to reconnect in 5 seconds...';
  console.log(disconnectMsg);
  this.msgBox.value += disconnectMsg + '\n';
  this.msgBox.scrollTop = this.msgBox.scrollHeight;

  setTimeout((function() {
    this.connectWs(this.conf.wsUrl);
  }).bind(this), 1000 * 5);
};

LoggedIn.prototype.onerror = function(e) {
  console.log('Connection error:', e);
  console.log('Attempting to connect to production server...');
  this.ws.onclose = null;
  this.connectWs(this.conf.wsProdUrl);
};

LoggedIn.prototype.onmessage = function(e) {
  var msg = JSON.parse(e.data);
  switch (msg.cmd) {
  case 'hello':
    console.log('websocket hello received:', msg);
    this.msgBox.value += 'websocket msg received: ' + msg.cmd + '\n';
    this.msgBox.scrollTop = this.msgBox.scrollHeight;
    break;
  default:
    console.log('websocket msg received:', msg);
    this.msgBox.value += 'websocket msg received: ' + msg.msg + '\n';
    this.msgBox.scrollTop = this.msgBox.scrollHeight;
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

LoggedIn.prototype.handleBroadcast = function() {
  var bcast = 'Broadcasting to all users: ';
  var msg = 'Hi everyone!';
  console.log(bcast, msg);
  this.msgBox.value += bcast + msg + '\n';
  this.msgBox.scrollTop = this.msgBox.scrollHeight;
  this.ws.send(JSON.stringify({
    cmd: 'broadcast',
    strArgs: [msg]
  }));
};

LoggedIn.prototype.handleLogout = function() {
  console.log('Logging out...');
};

(function() {
  'use strict';
  var loggedIn = new LoggedIn({
    backendUrl: '/api',
    wsUrl: 'ws://localhost:8081',
    wsProdUrl: 'ws://35.185.61.156:8081'
  });
})();
