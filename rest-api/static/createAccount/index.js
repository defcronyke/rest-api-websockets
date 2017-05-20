function CreateAccount(conf) {
  this.conf = conf;
  var button = document.getElementById('createAccount');
  button.onclick = this.handleCreateAccount.bind(this);
}

CreateAccount.prototype.handleCreateAccount = function() {
  var body = {
    email: document.getElementById('email').value,
    username: document.getElementById('username').value,
    password: document.getElementById('password').value,
    confirmPassword: document.getElementById('confirmPassword').value
  };
  this.postCreateAccount(body);
};

CreateAccount.prototype.postCreateAccount = function(body) {
  fetch(this.conf.backendUrl + '/createAccount', {
    method: 'POST',
    body: JSON.stringify(body)
  }).
  then(function(res) {
    return res.json();
  }).
  then(function(res) {
    console.log(res);
    if (!res.ok) {
      console.log('Error: Create new account failed');
      return;
    }
    localStorage.setItem("accessToken", res.jwt);
    document.location.href = '/loggedIn';
  });
};

(function() {
  'use strict';
  var createAccount = new CreateAccount({
    backendUrl: '/api'
  });
})();
