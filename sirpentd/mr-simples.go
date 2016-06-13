package main

import (
  "fmt"
  "net"
  "bufio"
  "github.com/Taneb/sirpent"
)

type World struct {
  G sirpent.Grid
  Snakes []sirpent.Snake
}

type Food struct {
  Position sirpent.Vector
}

func main() {
  port := "8080"
  ln, err := net.Listen("tcp", ":" + port)
  if err != nil {
    panic(fmt.Sprintf("Could not establish TCP server on port %s.", port))
    // handle error
  }
  for {
    conn, err := ln.Accept()
    if err != nil {
      // handle error
    }
    go handleConnection(conn)
  }
}

func handleConnection(conn net.Conn) {
  fmt.Fprintf(conn, "Hi.\r\n")
  status, err := bufio.NewReader(conn).ReadString('\n')
  if err != nil || status != "Hi back." {
    fmt.Printf("Bad status err='%s' status='%s'\n", err, status)
    return
  }

  width, err := bufio.NewReader(conn).ReadString('\n')
  height, err := bufio.NewReader(conn).ReadString('\n')
  snakes_len, err := bufio.NewReader(conn).ReadString('\n')

  start_x, err := bufio.NewReader(conn).ReadString('\n')
  start_y, err := bufio.NewReader(conn).ReadString('\n')

  /*for {
    status, err := bufio.NewReader(conn).ReadString('\n')
  }

  w := World{G: sirpent.HexagonalGrid{Width: 30, Height: 30}, Snakes: make([]sirpent.Snake, 1)}

  w.Snakes[0] = *sirpent.NewSnake(sirpent.HexagonalVector{X: 1, Y: 2})
  fmt.Println(w)
  w.Snakes[0].StepInDirection(sirpent.NORTHEAST)
  fmt.Println(w)
  w.Snakes[0].StepInDirection(sirpent.NORTHEAST)
  fmt.Println(w)
  w.Snakes[0].StepInDirection(sirpent.SOUTH)
  fmt.Println(w)
  w.Snakes[0].StepInDirection(sirpent.SOUTHEAST)
  fmt.Println(w)*/
}

/*
SERVER --> PLAYER
{
  "sirpent-server": ["sirpent-go", 0.0],
  "player": {
    "id": 2
  }
}
PLAYER --> SERVER
{
  "sirpent-player": ["mr-simples", 0.0]
}

SERVER --> PLAYER
{
  "method": "tick",
  "food": [
    {
      "position": [1, 2]
    }
  ],
  "snakes": {
    0 => {
      "length": 1,
      "dead": false,
      "segments": {
        "position": [3, 24]
      }
    },
    1 => {
      "length": 1,
      "dead": false,
      "segments": {
        "position": [14, 9]
      }
    },
    2 => {
      "length": 1,
      "dead": false,
      "segments": {
        "position": [21, 3]
      }
    }
  }
}

PLAYER --> SERVER
{
  "method": "tick-move",
  "direction": "NORTHEAST"
}
*/
