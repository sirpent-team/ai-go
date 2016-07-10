package sirpent

type HexGridUnshaped struct{}

// Cube coordinates, everywhere.
// http://www.redblobgames.com/grids/hexagons/#coordinates

func (g HexGridUnshaped) CellNeighbour(v Vector, d Direction) Vector {
	v2 := v
	switch d {
	case NorthEast:
		v2[0]++
		v2[2]--
	case East:
		v2[0]++
		v2[1]--
	case SouthEast:
		v2[1]--
		v2[2]++
	case SouthWest:
		v2[0]--
		v2[2]++
	case West:
		v2[0]--
		v2[1]++
	case NorthWest:
		v2[1]++
		v2[2]--
	}
	return v2
}

func (g HexGridUnshaped) DistanceBetweenCells(v1, v2 Vector) int {
	dx := abs(v2[0] - v1[0])
	dy := abs(v2[1] - v1[1])
	dz := abs(v2[2] - v1[2])
	return max(max(dx, dy), dz)
}
