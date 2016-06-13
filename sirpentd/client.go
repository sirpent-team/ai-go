package main

import (
  "fmt"
  "net"
)

func main() {
  domain := "127.0.0.1"
  port := "8080"

  target := fmt.Sprintf("%s:%s", domain, port)
  conn, err := net.Dial("tcp", target)
  if err != nil {
    panic(fmt.Sprintf("Could not establish connection to %s.", target))
  }

  fmt.Fprintf(conn, "MESSAGE.\n")
}
