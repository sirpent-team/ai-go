package sirpent

type Vector []int

type Grid interface {
	// Error allows for grids with an unbounded number of cells.
	Cells() ([]Vector, error)
	CryptoRandomCell() (Vector, error)
	CellNeighbour(v Vector, d Direction) Vector
	CeilNeighbours(v Vector) []Vector
	IsCellWithinBounds(v Vector) bool
	DistanceBetweenCells(v1, v2 Vector) int
}

func (g Grid) CeilNeighbours(v Vector) []Vector {
	directions := Directions()
	neighbours := make([]Vector, len(directions))
	for i := range directions {
		neighbours[i] = g.CellNeighbour(v, directions[i])
	}
	return neighbours
}
