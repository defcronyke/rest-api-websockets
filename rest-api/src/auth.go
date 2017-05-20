package RestApiWebsockets

import (
  "log"
  "net/http"
  "encoding/json"
  "fmt"
  "errors"
  "crypto/rand"
  "crypto/sha256"
  "encoding/base64"
  jwt "github.com/dgrijalva/jwt-go"
  "time"
  "strconv"
  "io/ioutil"
  "crypto/rsa"
)

type Auth struct {

}

type AuthSuccessRes struct {
  Ok  bool    `json:"ok"`
  Jwt string  `json:"jwt"`
}

type AuthErrRes struct {
  Ok    bool    `json:"ok"`
  Msg   string  `json:"msg"`
  Code  int     `json:"code"`
}

type CreateAccountReqBody struct {
  Username  string  `json:"username"`
  Password  string  `json:"password"`
  Email     string  `json:"email"`
}

type LoginReqBody struct {
  Username  string  `json:"username"`
  Password  string  `json:"password"`
}

// A POST route for creating new user accounts
func (a *Auth) CreateAccount(w http.ResponseWriter, r *http.Request) {
  log.Printf("Creating a new user account: %v", r)
  w.Header().Set("Content-Type", "application/json")

  // Get data from POST body
  d := json.NewDecoder(r.Body)
  var reqBody CreateAccountReqBody
  if err := d.Decode(&reqBody); err != nil {
    a.Err(fmt.Sprintf("Error: %v", err), http.StatusBadRequest, w, r)
    return
  }
  defer r.Body.Close()

  // Make the new user account
  var u *User
  if u = NewUser(reqBody.Username, reqBody.Password, reqBody.Email); u.Username == "" {
    a.Err("Error: Creating new user account failed", http.StatusInternalServerError, w, r)
    return
  }

  // Save the new account to the database
  db := NewDb(w, r)
  if db.PutCreateAccount(u); u.Username == "" {
    return
  }

  // Give the user a JWT so they can access authenticated routes
  var loginJwt string
  if loginJwt = a.GetLoginJwt(u.Username, r); loginJwt == "" {
    a.Err("Error: Getting new JWT failed", http.StatusInternalServerError, w, r)
    return
  }

  resBody := AuthSuccessRes{
    Ok: true,
    Jwt: loginJwt,
  }
  log.Printf("Created a new user account: %v", u)
  json.NewEncoder(w).Encode(resBody)
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
  log.Printf("Attempting to log in user: %v", r)
  w.Header().Set("Content-Type", "application/json")

  // Get data from POST body
  d := json.NewDecoder(r.Body)
  var reqBody LoginReqBody
  if err := d.Decode(&reqBody); err != nil {
    a.Err(fmt.Sprintf("Error: %v", err), http.StatusBadRequest, w, r)
    return
  }
  defer r.Body.Close()

  // Check username and passhash to see if user has authenticated successfully
  db := NewDb(w, r)
  var authenticated bool
  if authenticated = db.QueryCheckUsernameAndPasshash(reqBody); !authenticated {
    a.Err("Error: Username and/or password is incorrect", http.StatusUnauthorized, w, r)
    return
  }

  // Give the user a JWT so they can access authenticated routes
  var loginJwt string
  if loginJwt = a.GetLoginJwt(reqBody.Username, r); loginJwt == "" {
    a.Err("Error: Getting new JWT failed", http.StatusInternalServerError, w, r)
    return
  }

  resBody := AuthSuccessRes{
    Ok: true,
    Jwt: loginJwt,
  }
  log.Printf("User logged in: %v", reqBody.Username)
  json.NewEncoder(w).Encode(resBody)
}

func (a *Auth) GetLoginJwt(username string, r *http.Request) string {
  var err error
  iss := r.Host + "/login"
  numRandBytes := 24
  randBytes := make([]byte, numRandBytes)
  if _, err := rand.Read(randBytes); err != nil {
    log.Printf("Error: %v", err)
    return ""
  }
  jti := base64.URLEncoding.EncodeToString(randBytes) + strconv.FormatInt(time.Now().UnixNano(), 10)
  token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
    "sub": username,
    "iss": iss,
    "nbf": time.Now().Unix(),
    "exp": time.Now().Unix() + 60 * 60,
    "aud": []string{iss, username},
    "jti": jti,
  })
  var privKey []byte
  if privKey, err = ioutil.ReadFile("keys/jwt.key"); err != nil {
    log.Printf("Error: Failed loading private key from disk: %v", err)
    return ""
  }
  var parsedPrivKey *rsa.PrivateKey
  if parsedPrivKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKey); err != nil {
    log.Printf("Error: Failed parsing RS256 private key from PEM file: %v", err)
    return ""
  }
  var tokenStr string
  if tokenStr, err = token.SignedString(parsedPrivKey); err != nil {
    log.Printf("Error: Failed signing JWT: %v", err)
    return ""
  }
  return tokenStr
}

func (a *Auth) Err(msg string, code int, w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(code)
  log.Printf("Error: %v, Code: %v, Req: %v", msg, code, r)
  json.NewEncoder(w).Encode(AuthErrRes{
    Ok: false,
    Msg: msg,
    Code: code,
  })
}

func HashPassword(password, salt string) string {
  h := sha256.New()
  h.Write([]byte(salt + password))
  return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func HashPasswordNewSalt(password string) (string, string, error) {
  numRandBytes := 24
  randBytes := make([]byte, numRandBytes)
  if _, err := rand.Read(randBytes); err != nil {
    log.Printf("Error: %v", err)
    return "", "", errors.New(fmt.Sprintf("Error: Failed getting random bytes: %v", err))
  }
  salt := base64.URLEncoding.EncodeToString(randBytes)
  return HashPassword(password, salt), salt, nil
}
