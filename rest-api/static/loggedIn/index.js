function LoggedIn(conf) {
  this.conf = conf;
  var button = document.getElementById('logout');
  button.onclick = this.handleLogout.bind(this);
  console.log("logged in area");
}

LoggedIn.prototype.handleLogout = function() {
  console.log('Logging out...');
};

(function() {
  'use strict';
  var loggedIn = new LoggedIn({
    backendUrl: '/api'
  });
})();
