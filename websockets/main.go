package main

import (
  ws "./lib"
)

func main() {
  w := ws.NewWebsockets()
  w.Listen("instance-1:8081")
  w.Listen("localhost:8081")
}
