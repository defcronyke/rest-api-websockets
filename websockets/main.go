package main

import (
  ws "./lib"
)

func main() {
  w := ws.NewWebsockets()
  w.Listen("localhost:8081")
}
