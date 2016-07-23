package sirpent

type HexGridHexagonal struct {
	*HexGridUnshaped
	Rings int
}

func (g *HexGridHexagonal) origin() Vector {
	return Vector{0, 0, 0}
}

func (g *HexGridHexagonal) Cells() ([]Vector, error) {
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

func (g *HexGridHexagonal) CryptoRandomCell() (Vector, error) {
	v := Vector{0, 0, 0}
	v[0] = crypto_int(-g.Rings, g.Rings)
	v[1] = crypto_int(max(0-g.Rings, 0-g.Rings-v[0]), min(g.Rings, g.Rings-v[0]))
	v[2] = 0 - v[0] - v[1]
	return v, nil
}

func (g *HexGridHexagonal) IsCellWithinBounds(v Vector) bool {
	return g.DistanceBetweenCells(v, g.origin()) <= g.Rings
}

var _ Grid = (*HexGridHexagonal)(nil)
