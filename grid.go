package sirpent

// @TODO: Consider how to store grids. Want to support arbitrary shapes, be able to get a
// list of cells in the grid, check if a cell is within bounds.
type Grid struct {
	Radius int
	Origin AxialVector
}

func (g Grid) IsWithinBounds(av AxialVector) bool {
	return av.Distance(g.Origin) <= g.Radius
}

func (g Grid) CryptoRandomCell() AxialVector {
	cv := CubeVector{0, 0, 0}
	cv.X = crypto_int(-g.Radius, g.Radius)
	cv.Y = crypto_int(max(-g.Radius, -g.Radius-cv.X), min(g.Radius, g.Radius-cv.X))
	cv.Z = -cv.X - cv.Y
	return cv.AsAxialVector()
}
