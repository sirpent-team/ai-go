package sirpent

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Direction int

// Stringified with `stringer -type=Direction`, direction_string.go.
const (
	NorthEast Direction = iota
	East
	SouthEast
	SouthWest
	West
	NorthWest
)

func Directions() []Direction {
	return []Direction{NorthEast, East, SouthEast, SouthWest, West, NorthWest}
}

func DirectionByString(d string) (Direction, error) {
	directions := Directions()
	for i := range directions {
		if strings.ToLower(directions[i].String()) == strings.ToLower(d) {
			return directions[i], nil
		}
	}
	return NorthEast, errors.New(fmt.Sprintf("Direction '%s' not found.", d))
}

func (d Direction) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Direction) UnmarshalJSON(b []byte) error {
	var direction_str string
	err := json.Unmarshal(b, &direction_str)
	if err == nil {
		var d2 Direction
		d2, err = DirectionByString(direction_str)
		*d = d2
	}
	return err
}
