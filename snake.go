package sirpent

import (
	"fmt"
)

type CollidedWithPlayerError struct {
	CollidedWithPlayerID UUID
}

func (e CollidedWithPlayerError) Error() string {
	return fmt.Sprintf("Collided with player ID '%s'.", e.CollidedWithPlayerID)
}

type CollidedWithBoundsError struct{}

func (e CollidedWithBoundsError) Error() string {
	return fmt.Sprintf("Collided with bounds.")
}

type Snake []HexVector

func NewSnake(start_position HexVector) Snake {
	// @TODO: Ideal capacity of a snake?
	s := make([]HexVector, 1)
	s[0] = start_position
	return s
}

// Prevent updating the slice of the previous snake.
func (s Snake) Copy() Snake {
	s2 := make([]HexVector, len(s))
	copy(s2, s)
	return s2
}

func (s Snake) Move(direction Direction, grow_extra_segment bool) Snake {
	s2 := s.Copy()
	// Unless keeping an extra segment, discard final segment.
	if !grow_extra_segment {
		s2 = s2[:len(s)-1]
	}
	// Prepend the new, moved head.
	s2 = append(Snake{s[0].Neighbour(direction)}, s2...)

	return s2
}

// Collision detection.
func (s Snake) TailContains(av HexVector) bool {
	for i := 1; i < len(s); i++ {
		if s[i] == av {
			return true
		}
	}
	return false
}

func (s Snake) HeadIntersects(s2 Snake) bool {
	return s2.TailContains(s[0])
}

func (s Snake) HeadIntersectsSelf() bool {
	return s.HeadIntersects(s)
}
