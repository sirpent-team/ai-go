package sirpent

import (
  "fmt"
)

// @TODO: Wrap these in a Direction label for easy iteration?
const (
  NORTH = iota
  NORTHEAST = iota
  SOUTHEAST = iota
  SOUTH = iota
  SOUTHWEST = iota
  NORTHWEST = iota
)

type HexagonalVector struct {
  X int
  Y int
}

func (v HexagonalVector) Distance(v2 HexagonalVector) HexagonalVector {
  return HexagonalVector{X: v2.X - v.X, Y: v2.Y - v.Y}
}

func (v HexagonalVector) Neighbour(direction int) HexagonalVector {
  switch direction {
  case NORTH:
    v.Y--
  case NORTHEAST:
    v.X++
    v.Y--
  case SOUTHEAST:
    v.X++
  case SOUTH:
    v.Y++
  case SOUTHWEST:
    v.X--
    v.Y++
  case NORTHWEST:
    v.X--
  default: // @TODO: Return an actual error if the direction was invalid.
    v.X = -1
    v.Y = -1
  }
  return v
}

func (v HexagonalVector) String() string {
  return fmt.Sprintf("(%d,%d)", v.X, v.Y)
}

type HexagonalGrid struct {
  Width int
  Height int
}

func (g HexagonalGrid) Dimensions() []int {
  return []int{g.Width, g.Height}
}

func (g HexagonalGrid) IsWithinBounds(v HexagonalVector) bool {
  return (v.X >= 0 && v.X < g.Width && v.Y >= 0 && v.Y < g.Height)
}
