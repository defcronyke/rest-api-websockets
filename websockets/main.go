package main

import (
  ws "./lib"
  // "net"
  // "log"
)

func main() {
  // Get private IP address
  // var err error
  // var serveIp string
  // var ifaces []net.Interface
  // if ifaces, err = net.Interfaces(); err != nil {
  //   log.Printf("Error: Failed getting network interfaces: %v", err)
  //   return
  // }
  // for _, i := range ifaces {
  //   var addrs []net.Addr
  //   if addrs, err = i.Addrs(); err != nil {
  //     log.Printf("Error: Failed getting network addresses: %v", err)
  //     return
  //   }
  //   for _, addr := range addrs {
  //     var ip net.IP
  //     switch v := addr.(type) {
  //     case *net.IPNet:
  //       ip = v.IP
  //     case *net.IPAddr:
  //       ip = v.IP
  //     }
  //     log.Printf("Found local IP address: %v", ip)
  //   }
  // }

  w := ws.NewWebsockets()
  w.Listen("instance-1:8081")
  w.Listen("localhost:8081")
}
