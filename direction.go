package sirpent

import (
	"fmt"
)

type DirectionError struct {
	DirectionValue Direction
}

func (e DirectionError) Error() string {
	return fmt.Sprintf("Direction '%s' not found.", e.DirectionValue)
}

type Direction string
