package main

import (
  "fmt"
  "net"
  "bufio"
)

func main() {
  port := "8901"

  ln, err := net.Listen("tcp", ":" + port)
  if err != nil {
    panic(fmt.Sprintf("Could not establish TCP server on port %s.", port))
  }

  for {
    conn, err := ln.Accept()
    if err != nil {
      panic(fmt.Sprintf("Could not listen on port %s.", port))
    }

    go handleConnection(conn)
  }
}

func handleConnection(conn net.Conn) {
  //r := bufio.NewReader(conn)

  for {
    for i := 0; i < 3; i++ {
      bufio.NewReader(conn).ReadString('\n')
      conn.Write([]byte("SOUTHEAST\n"))
    }

    for i := 0; i < 3; i++ {
      bufio.NewReader(conn).ReadString('\n')
      conn.Write([]byte("NORTH\n"))
    }

    for i := 0; i < 3; i++ {
      bufio.NewReader(conn).ReadString('\n')
      conn.Write([]byte("NORTHWEST\n"))
    }

    for i := 0; i < 3; i++ {
      bufio.NewReader(conn).ReadString('\n')
      conn.Write([]byte("SOUTH\n"))
    }
  }

  /*s, err := r.ReadString('\n')
  if err != nil {
    panic(fmt.Sprintf("Error reading from connection: %s", err))
  }
  fmt.Printf("s = %s\n", s)*/
}
