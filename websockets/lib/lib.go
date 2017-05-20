package ws

import (
  "github.com/gorilla/websocket"
  "net/http"
  "log"
)

type Websockets struct {

}

func NewWebsockets() (*Websockets) {
  return &Websockets{}
}

func (ws *Websockets) Listen(addr string) {
  http.HandleFunc("/", ws.ConnectionHandler)
  log.Printf("Websocket server listening at: ws://%v", addr)
  log.Fatal(http.ListenAndServe(addr, nil))
}

func (ws *Websockets) ConnectionHandler(w http.ResponseWriter, r *http.Request) {
  var err error
  // var c *websocket.Conn
  upgrader := websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
      return true
    },
  }
  if _, err = upgrader.Upgrade(w, r, nil); err != nil {
    log.Printf("Error: Websocket upgrade failed: %v", err)
    return
  }
  log.Printf("Websocket upgrade success. Connection from: %v", r.RemoteAddr)
}
