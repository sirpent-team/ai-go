package sirpent

import (
	"encoding/json"
	"errors"
)

type HexGridHexagonal struct {
	HexGridUnshaped
	Rings int
}

func (g HexGridHexagonal) origin() Vector {
	return Vector{0, 0, 0}
}

func (g HexGridHexagonal) Cells() ([]Vector, error) {
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

func (g HexGridHexagonal) CryptoRandomCell() (Vector, error) {
	v := Vector{0, 0, 0}
	v[0] = crypto_int(-g.Rings, g.Rings)
	v[1] = crypto_int(max(0-g.Rings, 0-g.Rings-v[0]), min(g.Rings, g.Rings-v[0]))
	v[2] = 0 - v[0] - v[1]
	return v, nil
}

func (g HexGridHexagonal) IsCellWithinBounds(v Vector) bool {
	return g.DistanceBetweenCells(v, g.origin()) <= g.Rings
}

func (g HexGridHexagonal) MarshalJSON() ([]byte, error) {
	g_for_json := struct {
		GridType string `json:"grid_type"`
		Rings    int    `json:"rings"`
	}{
		GridType: "hex_grid_hexagonal",
		Rings:    g.Rings,
	}

	return json.Marshal(g_for_json)
}

func (g *HexGridHexagonal) UnmarshalJSON(b []byte) error {
	g_for_json := struct {
		GridType string `json:"grid_type"`
		Rings    int    `json:"rings"`
	}{}
	err := json.Unmarshal(b, &g_for_json)
	if err != nil {
		return err
	}

	if g_for_json.GridType != "hex_grid_hexagonal" {
		return errors.New("Decoding wrong grid type.")
	}

	g.Rings = g_for_json.Rings

	return nil
}

var _ Grid = (*HexGridHexagonal)(nil)
