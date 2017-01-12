package sirpent

/*
"game" : {
  "uuid" : "f57da97d-685c-4eac-9089-df86b85ac5c4",
  "grid" : {
     "radius" : 25
  },
  "players" : [
     "46bit_",
     "46bit"
  ]
}
*/

type Game struct {
	ID   string `json:"uuid"`
	Grid Grid   `json:"grid"`
	// All players in this game.
	Players []string `json:"players"`
}
