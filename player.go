package sirpent

import (
  "net"
  "fmt"
  "time"
  "bufio"
  "errors"
)

type Player struct {
  ServerLocation string
  Connection net.Conn
  S *Snake
}

func NewPlayer(server_location string) *Player {
  p := Player{ServerLocation: server_location, Connection: nil, S: nil}
  return &p
}

func (p *Player) ConnectToPlayer() error {
  c, err := net.DialTimeout("tcp", p.ServerLocation, time.Duration(10 * time.Second))
  p.Connection = c
  return err
}

func (p Player) SendWorld(w *World) error {
  fmt.Fprintf(p.Connection, "SENDING_WORLD\r\n")

  return nil//err
}

func (p Player) ReceiveMove() (int, error) {
  fmt.Fprintf(p.Connection, "MOVE PLS\n")
  direction_msg, err := bufio.NewReader(p.Connection).ReadString('\n')
  fmt.Printf("\n### '%s' / '%s' ###\n\n", direction_msg, err)
  direction := 0
  if err != nil {
    return direction, err
  }

  // @TODO: Have an actual protocol to handle framing and be sane.
  switch direction_msg {
  case "NORTH\n":
    direction = NORTH
  case "NORTHEAST\n":
    direction = NORTHEAST
  case "SOUTHEAST\n":
    direction = SOUTHEAST
  case "SOUTH\n":
    direction = SOUTH
  case "SOUTHWEST\n":
    direction = SOUTHWEST
  case "NORTHWEST\n":
    direction = NORTHWEST
  default:
    err = errors.New(fmt.Sprintf("Unrecognised direction from player '%s'.", direction_msg))
  }

  return direction, err
}
