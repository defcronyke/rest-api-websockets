package ws

import (
  "github.com/gorilla/websocket"
  "net/http"
  "log"
  "fmt"
)

type Websockets struct {

}

type Msg struct {
  Cmd string  `json:"cmd"`
}

type MsgWithArgs struct {
  Cmd     string    `json:"cmd"`
  StrArgs []string  `json:"strArgs"`
  IntArgs []int64   `json:"intArgs"`
}

type ErrRes struct {
  Code  int     `json:"code"`
  Msg   string  `json:"msg"`
}

func NewWebsockets() (*Websockets) {
  return &Websockets{}
}

func (ws *Websockets) Listen(addr string) {
  http.HandleFunc("/", ws.ConnectionHandler)
  log.Printf("Websocket server listening at: ws://%v", addr)
  log.Fatal(http.ListenAndServe(addr, nil))
}

// Handle incoming websocket connections
func (ws *Websockets) ConnectionHandler(w http.ResponseWriter, r *http.Request) {
  // Upgrade HTTP connection to a websocket connection
  var err error
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

  // Listen for incoming websocket messages
  for {
    var msg MsgWithArgs
    if err = c.ReadJSON(&msg); err != nil {
      errMsg := fmt.Sprintf("Error: ReadJSON failed: %v", err)
      log.Printf("%v", errMsg)
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
