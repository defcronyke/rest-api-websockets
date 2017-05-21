package ws

import (
  "github.com/gorilla/websocket"
  "net/http"
  "log"
  "fmt"
  "sync"
)

type Websockets struct {
  Users       map[string]*websocket.Conn  // A list of connected users
  UsersMx     *sync.Mutex
  BcastCh     chan Broadcast
  BcastUserCh chan Broadcast
}

type MsgWithArgs struct {
  Cmd     string    `json:"cmd"`
  StrArgs []string  `json:"strArgs"`
  IntArgs []int64   `json:"intArgs"`
}

type Broadcast struct {
  User  *websocket.Conn
  Res   BroadcastRes
}

type BroadcastRes struct {
  Cmd string  `json:"cmd"`
  Msg string  `json:"msg"`
}

type ErrRes struct {
  Code  int     `json:"code"`
  Msg   string  `json:"msg"`
}

func NewWebsockets() (*Websockets) {
  ws := &Websockets{
    Users:        make(map[string]*websocket.Conn),
    UsersMx:      &sync.Mutex{},
    BcastCh:      make(chan Broadcast),
    BcastUserCh:  make(chan Broadcast),
  }

  http.HandleFunc("/", ws.ConnectionHandler)

  // Listen for broadcasts to all users
  go func() {
    for {
      msg := <- ws.BcastCh
      ws.Broadcast(msg.User, msg.Res)
    }
  }()

  // Listen to broadcasts for each individual user
  go func() {
    var err error
    for {
      msg := <- ws.BcastUserCh
      if err = msg.User.WriteJSON(msg.Res); err != nil {
        log.Printf("WriteJSON failed: %v", err)
        continue
      }
    }
  }()

  return ws
}

func (ws *Websockets) Listen(addr string) {
  var err error
  log.Printf("Websocket server listening at: ws://%v", addr)
  if err = http.ListenAndServe(addr, nil); err != nil {
    log.Printf("Error: Listening failed: %v", err)
    return
  }
}

// Handle incoming websocket connections
func (ws *Websockets) ConnectionHandler(w http.ResponseWriter, r *http.Request) {
  // Upgrade HTTP connection to a websocket connection
  var err error

  // Check JWT access token to make sure user is authenticated

  upgrader := websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
      return true
    },
  }
  var c *websocket.Conn
  if c, err = upgrader.Upgrade(w, r, nil); err != nil {
    log.Printf("Error: Websocket upgrade failed: %v", err)
    return
  }
  defer c.Close()
  log.Printf("Websocket upgrade success. Connection from: %v", r.RemoteAddr)

  // Add user to list of connections
  ws.UsersMx.Lock()
  ws.Users[r.RemoteAddr] = c
  ws.UsersMx.Unlock()
  log.Printf("Number of connections: %v", len(ws.Users))

  // Listen for incoming websocket messages
  for {
    var msg MsgWithArgs
    if err = c.ReadJSON(&msg); err != nil {
      errMsg := fmt.Sprintf("Error: ReadJSON failed: %v", err)
      log.Printf("%v", errMsg)

      // Remove user from list of connections
      ws.UsersMx.Lock()
      delete(ws.Users, r.RemoteAddr)
      ws.UsersMx.Unlock()
      log.Printf("Number of connections: %v", len(ws.Users))

      if err = c.WriteJSON(ErrRes{
        Code: http.StatusInternalServerError,
        Msg: errMsg,
      }); err != nil {
        log.Printf("WriteJSON failed: %v", err)
        return
      }
      return
    }
    log.Printf("Msg: %+v", msg)

    switch msg.Cmd {
    case "hello":
      if err = c.WriteJSON(MsgWithArgs{
        Cmd: "hello",
      }); err != nil {
        log.Printf("WriteJSON failed: %v", err)
        return
      }

    case "broadcast":
      ws.BcastCh <- Broadcast{
        User: c,
        Res:  BroadcastRes{
          Cmd:  msg.Cmd,
          Msg:  msg.StrArgs[0],
        },
      }

    default:
      if err = c.WriteJSON(ErrRes{
        Code: http.StatusMethodNotAllowed,
        Msg: fmt.Sprintf("Error: Cmd not understood: %v", msg.Cmd),
      }); err != nil {
        log.Printf("WriteJSON failed: %v", err)
        return
      }
    }
  }
}

func (ws *Websockets) Broadcast(c *websocket.Conn, msg BroadcastRes) {
  for _, v := range ws.Users {
    res := Broadcast{
      User: v,
      Res: BroadcastRes{
        Cmd: "broadcast",
        Msg: msg.Msg,
      },
    }

    if (c != v) {
      ws.BcastUserCh <- res
    }
  }
}
