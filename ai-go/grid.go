package sirpent

import (
	"fmt"
)

type DirectionError struct {
	DirectionValue Direction
}

func (e DirectionError) Error() string {
	return fmt.Sprintf("Direction '%s' not found.", e.DirectionValue)
}

type Direction string

type Vector struct {
	X int
	Y int
}

func (v Vector) Eq(v2 Vector) bool {
	return v.X == v2.X && v.Y == v2.Y
}

type HexagonalGrid struct {
  Radius int `json:"radius"`
}

func (g *HexagonalGrid) origin() Vector {
  return Vector{0, 0}
}

func (g *HexagonalGrid) Cells() ([]Vector, error) {
  cells := make([]Vector, 0)
  for x := -g.Radius; x <= g.Radius; x++ {
    for y := -g.Radius; y <= g.Radius; y++ {
      cells = append(cells, Vector{x, y})
    }
  }
  return cells, nil
}

func (g *HexagonalGrid) CryptoRandomCell() (Vector, error) {
  v := Vector{
    X: crypto_int(-g.Radius, g.Radius),
    Y: crypto_int(-g.Radius, g.Radius),
  }
  return v, nil
}

func (g *HexagonalGrid) Directions() []Direction {
  return []Direction{"north", "northeast", "southeast", "south", "southwest", "northwest"}
}

func (g *HexagonalGrid) ValidateDirection(d Direction) error {
  directions := g.Directions()
  for i := range directions {
    if directions[i] == d {
      return nil
    }
  }
  return DirectionError{DirectionValue: d}
}

func (g *HexagonalGrid) CellNeighbour(v Vector, d Direction) Vector {
  neighbour := Vector{v.X, v.Y}
  switch d {
  case "north":
    neighbour.Y--
  case "northeast":
    neighbour.X++
    neighbour.Y--
  case "southeast":
    neighbour.X++
  case "south":
    neighbour.Y++
  case "southwest":
    neighbour.X--
    neighbour.Y++
  case "northwest":
    neighbour.X--
  }
  return neighbour
}

func (g *HexagonalGrid) CellNeighbours(v Vector) []Vector {
  directions := g.Directions()
  neighbours := make([]Vector, len(directions))
  for i := range directions {
    neighbours[i] = g.CellNeighbour(v, directions[i])
  }
  return neighbours
}

func (g *HexagonalGrid) IsCellWithinBounds(v Vector) bool {
  return g.DistanceBetweenCells(v, g.origin()) <= g.Radius
}

func (g *HexagonalGrid) DistanceBetweenCells(v1, v2 Vector) int {
  dx := abs(v2.X - v1.X)
  dy := abs(v2.Y - v1.Y)
  dz := abs((v1.X + v1.Y) - (v2.X + v2.Y))
  return max(max(dx, dy), dz)
}
