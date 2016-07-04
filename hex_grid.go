package sirpent

type HexGrid interface {
	IsWithinBounds(v HexVector) bool
	CryptoRandomCell() HexVector
	Cells() []HexVector
}

type HexHexGrid struct {
	Rings int `json:"rings"`
}

func NewHexHexGrid(rings int) HexGrid {
	return &HexHexGrid{Rings: rings}
}

func (hhg HexHexGrid) Origin() HexVector {
	return HexVector{0, 0, 0}
}

func (hhg HexHexGrid) IsWithinBounds(v HexVector) bool {
	return v.Distance(hhg.Origin()) <= hhg.Rings
}

func (hhg HexHexGrid) CryptoRandomCell() HexVector {
	v := HexVector{0, 0, 0}
	v.X = crypto_int(-hhg.Rings, hhg.Rings)
	v.Y = crypto_int(max(0-hhg.Rings, 0-hhg.Rings-v.X), min(hhg.Rings, hhg.Rings-v.X))
	v.Z = 0 - v.X - v.Y
	return v
}

func (hhg HexHexGrid) Cells() []HexVector {
	cells := make([]HexVector, 0)
	for x := -hhg.Rings; x <= hhg.Rings; x++ {
		for y := -hhg.Rings; y <= hhg.Rings; y++ {
			z := 0 - x - y
			if z == 0 {
				cells = append(cells, HexVector{x, y, z})
			}
		}
	}
	return cells
}

var _ HexGrid = (*HexHexGrid)(nil)
