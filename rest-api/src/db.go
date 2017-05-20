package RestApiWebsockets

import (
  "github.com/qedus/nds"
  "google.golang.org/appengine"
  "google.golang.org/appengine/datastore"
  "fmt"
  "log"
  "net/http"
  "golang.org/x/net/context"
  "encoding/json"
)

type Db struct {
  r *http.Request
  w http.ResponseWriter
  c context.Context
}

type DbErrRes struct {
  Ok    bool    `json:"ok"`
  Msg   string  `json:"msg"`
  Code  int     `json:"code"`
}

func NewDb(w http.ResponseWriter, r *http.Request) (*Db) {
  return &Db{
    r: r,
    w: w,
    c: appengine.NewContext(r),
  }
}

func (d *Db) PutCreateAccount(u *User) {
  var err error
  // Check db to see if username or email address are already taken
  if accountExists := d.QueryCheckIfAccountExists(u); accountExists {
    u.Username = ""
    return
  }
  // Put new user into db
  k := datastore.NewKey(d.c, "user", "", 0, nil)
  if k, err = nds.Put(d.c, k, u); err != nil {
    d.Err(fmt.Sprintf("Putting new user account into database failed: %v", err), http.StatusInternalServerError)
    return
  }
  log.Printf("Put new user account into database: %v", u)
}

func (d *Db) QueryCheckIfAccountExists(u *User) bool {
  if exists := d.QueryCheckIfUsernameExists(u); exists {
    return true
  }
  if exists := d.QueryCheckIfEmailExists(u); exists {
    return true
  }
  return false
}

// Check if username already exists in db
func (d *Db) QueryCheckIfUsernameExists(u *User) bool {
  var err error
  q := datastore.NewQuery("user").
  Filter("Username =", u.Username)
  for res := q.Run(d.c); ; {
    var resUser User
    if _, err = res.Next(&resUser); err == datastore.Done {
      return false
    } else if err != nil {
      d.Err(fmt.Sprintf("Error: Checking if username already exists failed: %v", err), http.StatusInternalServerError)
      return true
    }
    d.Err("Error: A user with that username already exists", http.StatusUnauthorized)
    return true
  }
}

// Check if email address already exists in db
func (d *Db) QueryCheckIfEmailExists(u *User) bool {
  var err error
  q := datastore.NewQuery("user").
  Filter("Email =", u.Email)
  for res := q.Run(d.c); ; {
    var resUser User
    if _, err = res.Next(&resUser); err == datastore.Done {
      return false
    } else if err != nil {
      d.Err(fmt.Sprintf("Error: Checking if email address already exists failed: %v", err), http.StatusInternalServerError)
      return true
    }
    d.Err("Error: A user with that email address already exists", http.StatusUnauthorized)
    return true
  }
}

func (d *Db) QueryCheckUsernameAndPasshash(reqBody LoginReqBody) bool {
  var user *User
  if user = d.QueryGetUser(reqBody.Username); user.Username == "" {
    return false
  }

  if passhash := HashPassword(reqBody.Password, user.Salt); passhash != user.Passhash {
    return false
  }
  return true
}


func (d *Db) QueryGetUser(username string) (*User) {
  var err error
  var resUser User
  q := datastore.NewQuery("user").
  Filter("Username =", username)
  for res := q.Run(d.c); ; {
    if _, err = res.Next(&resUser); err == datastore.Done {
      return &User{}
    } else if err != nil {
      d.Err(fmt.Sprintf("Error: Failed getting user from database: %v", err), http.StatusInternalServerError)
      return &User{}
    }
    return &resUser
  }
}

func (d *Db) Err(msg string, code int) {
  d.w.WriteHeader(code)
  log.Printf("Error: %v, Code: %v, Req: %v", msg, code, d.r)
  json.NewEncoder(d.w).Encode(DbErrRes{
    Ok: false,
    Msg: msg,
    Code: code,
  })
}
