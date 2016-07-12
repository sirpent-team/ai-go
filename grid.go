package sirpent

import (
	"encoding/json"
	"errors"
)

type Vector []int

func (v Vector) Eq(v2 Vector) bool {
	if len(v) != len(v2) {
		return false
	}
	for i := range v {
		if i >= len(v2) || v[i] != v2[i] {
			return false
		}
	}
	return true
}

type Grid interface {
	// Error allows for grids with an unbounded number of cells.
	Cells() ([]Vector, error)
	CryptoRandomCell() (Vector, error)

	Directions() []Direction
	// Error if direction invalid. Makes up for being unable to typecheck grid-specific directions.
	ValidateDirection(d Direction) error
	CellNeighbour(v Vector, d Direction) Vector
	CellNeighbours(v Vector) []Vector

	IsCellWithinBounds(v Vector) bool
	DistanceBetweenCells(v1, v2 Vector) int
}

func ParseGridJSON(b []byte) (Grid, error) {
	g_for_json := struct {
		GridType string `json:"grid_type"`
		Rings    int    `json:"rings"`
	}{}
	err := json.Unmarshal(b, &g_for_json)
	if err != nil {
		return nil, err
	}

	var grid Grid
	switch g_for_json.GridType {
	case "hex_grid_hexagonal":
		grid = &HexGridHexagonal{Rings: g_for_json.Rings}
	}

	if grid == nil {
		return nil, errors.New("Unknown Grid Type.")
	}
	return grid, nil
}
