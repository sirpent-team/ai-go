package sirpent

import (
	"fmt"
)

// @TODO: Wrap these in a Direction label for easy iteration?
const (
	NORTH     = iota
	NORTHEAST = iota
	SOUTHEAST = iota
	SOUTH     = iota
	SOUTHWEST = iota
	NORTHWEST = iota
)

type HexagonalVector struct {
	X, Y int
}

// Coord, Distance and Neighbour all implement methods defined upon the Vector interface. The aim is to allow
// arbitrary coordinate systems to be dropped in later, whether that be 2D orthogonal grids or 5D triangular grids.
//
// Defining these methods upon the interface has the advantage of forcing all Vectors to implement them. However
// it also requires all vector systems to accept vectors from other systems in order to typecheck. Furthermore
// as Vector.X is not defined, we can't use those as the field isn't necessarily defined.
//
// As such the easiest workaround is a function to retrieve coordinates. X and Y are being kept exported for now
// until it's decided how to instantiate a HexagonalVector. It might make more sense to use an array or slice
// instead of X and Y, provided the compiler will still optimise all the overhead out (questionable).
func (v HexagonalVector) Coord(offset int) int {
	switch offset {
	case 0:
		return v.X
	case 1:
		return v.Y
	default:
		return 0
	}
}

func (v HexagonalVector) Distance(v2 Vector) Vector {
	return HexagonalVector{X: v2.Coord(0) - v.Coord(0), Y: v2.Coord(1) - v.Coord(1)}
}

func (v HexagonalVector) Neighbour(direction int) Vector {
	v2 := HexagonalVector{v.Coord(0), v.Coord(1)}
	switch direction {
	case NORTH:
		v2.Y--
	case NORTHEAST:
		v2.X++
		v2.Y--
	case SOUTHEAST:
		v2.X++
		v2.Y += 1 - (v.X & 1)
	case SOUTH:
		v2.Y++
	case SOUTHWEST:
		v2.X--
		v2.Y++
	case NORTHWEST:
		v2.X--
		v2.Y -= v.X & 1
	default: // @TODO: Return an actual error if the direction was invalid.
		v2.X = -1
		v2.Y = -1
	}
	return v2
}

func (v HexagonalVector) String() string {
	return fmt.Sprintf("(%d,%d)", v.X, v.Y)
}

var _ Vector = (*HexagonalVector)(nil)

type HexagonalGrid struct {
	Width, Height int
}

func (g HexagonalGrid) Dimensions() []int {
	return []int{g.Width, g.Height}
}

func (g HexagonalGrid) IsWithinBounds(v Vector) bool {
	return (v.Coord(0) >= 0 && v.Coord(0) < g.Width && v.Coord(1) >= 0 && v.Coord(1) < g.Height)
}

var _ Grid = (*HexagonalGrid)(nil)
