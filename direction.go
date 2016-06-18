package sirpent

import (
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
		if strings.ToUpper(directions[i].String()) == d {
			return directions[i], nil
		}
	}
	return NorthEast, errors.New(fmt.Sprintf("Direction '%s' not found.", d))
}
