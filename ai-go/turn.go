package sirpent

import "encoding/json"

/*
"turn" : {
  "turn_number" : 0,
  "directions" : {},
  "food" : [
     {"x" : 1, "y" : 0}
  ],
  "snakes" : {
     "46bit_" : {
        "segments" : [
           {"y" : 22, "x" : -17 }
        ]
     },
     "46bit" : {
        "segments" : [
           {"y" : 10, "x" : -8}
        ]
     }
  },
  "eaten" : {}
  "casualties" : {},
}
*/

type Turn struct {
	TurnNumber uint                       `json:"turn_number"`
	Directions map[string]Direction       `json:"directions"`
	Food       []Vector                   `json:"food"`
	Snakes     map[string]Snake           `json:"snakes"`
	Eaten      map[string]Vector          `json:"eaten"`
	Casualties map[string]json.RawMessage `json:"casualties"`
}
