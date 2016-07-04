package sirpent

import (
	"fmt"
)

type HexVector struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

func (v HexVector) Distance(v2 HexVector) int {
	return max(max(abs(v2.X-v.X), abs(v2.Y-v.Y)), abs(v2.Z-v.Z))
}

func (v HexVector) Neighbour(d Direction) HexVector {
	v2 := v
	switch d {
	case NorthEast:
		v2.X++
		v2.Z--
	case East:
		v2.X++
		v2.Y--
	case SouthEast:
		v2.Y--
		v2.Z++
	case SouthWest:
		v2.X--
		v2.Z++
	case West:
		v2.X--
		v2.Y++
	case NorthWest:
		v2.Y++
		v2.Z--
	}
	return v2
}

func (v HexVector) Neighbours() []HexVector {
	directions := Directions()
	neighbours := make([]HexVector, len(directions))
	for i := range directions {
		neighbours[i] = v.Neighbour(directions[i])
	}
	return neighbours
}

func (v HexVector) String() string {
	return fmt.Sprintf("HexVector(X=%d, Y=%d, Z=%d)", v.X, v.Y, v.Z)
}
