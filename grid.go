package sirpent

import (
  "golang.org/x/net/websocket"
)

type Vector interface {
	Coord(offset int) int
	Distance(v2 Vector) Vector
	Neighbour(direction int) Vector
}

type Grid interface {
	Dimensions() []int
	IsWithinBounds(v Vector) bool
}

type World struct {
  G Grid
  Players []*Player
  Websockets []*websocket.Conn
}
