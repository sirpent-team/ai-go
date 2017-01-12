package sirpent

type HexagonalGrid struct {
	Rings int
}

func (g *HexagonalGrid) origin() Vector {
	return Vector{0, 0, 0}
}

func (g *HexagonalGrid) Cells() ([]Vector, error) {
	cells := make([]Vector, 0)
	for x := -g.Rings; x <= g.Rings; x++ {
		for y := -g.Rings; y <= g.Rings; y++ {
			z := 0 - x - y
			if z == 0 {
				cells = append(cells, Vector{x, y, z})
			}
		}
	}
	return cells, nil
}

func (g *HexagonalGrid) CryptoRandomCell() (Vector, error) {
	v := Vector{0, 0, 0}
	v[0] = crypto_int(-g.Rings, g.Rings)
	v[1] = crypto_int(max(0-g.Rings, 0-g.Rings-v[0]), min(g.Rings, g.Rings-v[0]))
	v[2] = 0 - v[0] - v[1]
	return v, nil
}

func (g *HexagonalGrid) Directions() []Direction {
	return []Direction{"NORTHEAST", "EAST", "SOUTHEAST", "SOUTHWEST", "WEST", "NORTHWEST"}
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

func (g *HexagonalGrid) CellNeighbours(v Vector) []Vector {
	directions := g.Directions()
	neighbours := make([]Vector, len(directions))
	for i := range directions {
		neighbours[i] = g.CellNeighbour(v, directions[i])
	}
	return neighbours
}

func (g *HexagonalGrid) IsCellWithinBounds(v Vector) bool {
	return g.DistanceBetweenCells(v, g.origin()) <= g.Rings
}

func (g *HexagonalGrid) DistanceBetweenCells(v1, v2 Vector) int {
	dx := abs(v2[0] - v1[0])
	dy := abs(v2[1] - v1[1])
	dz := abs(v2[2] - v1[2])
	return max(max(dx, dy), dz)
}

var _ Grid = (*HexagonalGrid)(nil)
