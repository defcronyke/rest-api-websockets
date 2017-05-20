package RestApiWebsockets

import (
  "log"
)

type User struct {
  Username  string  `json:"username"`
  Passhash  string  `json:"-"`
  Email     string  `json:"email"`
  Salt      string  `json:"-"`
}

func NewUser(username, password, email string) (*User) {
  var err error
  var passhash string
  var salt string
  if passhash, salt, err = HashPasswordNewSalt(password); err != nil {
    log.Printf("Error: Failed hashing password with new salt: %v", err)
    return &User{}
  }

  return &User{
    Username: username,
    Passhash: passhash,
    Email:    email,
    Salt:     salt,
  }
}
