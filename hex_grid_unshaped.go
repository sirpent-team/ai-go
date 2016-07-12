package sirpent

type HexGridUnshaped struct{}

// Cube coordinates, everywhere.
// http://www.redblobgames.com/grids/hexagons/#coordinatestype Direction interface {

func (g *HexGridUnshaped) Directions() []Direction {
	return []Direction{"NORTHEAST", "EAST", "SOUTHEAST", "SOUTHWEST", "WEST", "NORTHWEST"}
}

func (g *HexGridUnshaped) ValidateDirection(d Direction) error {
	directions := g.Directions()
	for i := range directions {
		if directions[i] == d {
			return nil
		}
	}
	return DirectionError{DirectionValue: d}
}

func (g *HexGridUnshaped) CellNeighbour(v Vector, d Direction) Vector {
	neighbour := Vector{v[0], v[1], v[2]}
	switch d {
	case "NORTHEAST":
		neighbour[0]++
		neighbour[2]--
	case "EAST":
		neighbour[0]++
		neighbour[1]--
	case "SOUTHEAST":
		neighbour[1]--
		neighbour[2]++
	case "SOUTHWEST":
		neighbour[0]--
		neighbour[2]++
	case "WEST":
		neighbour[0]--
		neighbour[1]++
	case "NORTHWEST":
		neighbour[1]++
		neighbour[2]--
	}
	return neighbour
}

func (g *HexGridUnshaped) CellNeighbours(v Vector) []Vector {
	directions := g.Directions()
	neighbours := make([]Vector, len(directions))
	for i := range directions {
		neighbours[i] = g.CellNeighbour(v, directions[i])
	}
	return neighbours
}

func (g *HexGridUnshaped) DistanceBetweenCells(v1, v2 Vector) int {
	dx := abs(v2[0] - v1[0])
	dy := abs(v2[1] - v1[1])
	dz := abs(v2[2] - v1[2])
	return max(max(dx, dy), dz)
}
