package sirpent

import (
	"encoding/json"
	"fmt"
)

type World struct {
	GridKind GridKind `json:"grid_kind"`
	Grid     Grid     `json:"grid"`
}

func (w *World) UnmarshalJSON(b []byte) error {
	var gk GridKind
	if err := json.Unmarshal([]byte("\"hex_grid_hexagonal\""), &gk); err != nil {
		return err
		//log.Fatal(err)
	}

	var raw2 json.RawMessage
	w3 := struct {
		GridKind GridKind    `json:"grid_kind"`
		Grid     interface{} `json:"grid"`
	}{Grid: &raw2}
	if err := json.Unmarshal(b, &w3); err != nil {
		fmt.Println("afff")
		fmt.Println(string(b))
		return err
		//log.Fatal(err)
	}

	var raw json.RawMessage
	w2 := struct {
		GridKind GridKind    `json:"grid_kind"`
		Grid     interface{} `json:"grid"`
	}{Grid: &raw}

	if err := json.Unmarshal(b, &w2); err != nil {
		return err
		//log.Fatal(err)
	}
	w.GridKind = w2.GridKind

	grid := gridKingHandlers[w2.GridKind]()
	if err := json.Unmarshal(raw, &grid); err != nil {
		return err
	}
	w.Grid = grid

	return nil
}
