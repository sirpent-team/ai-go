package sirpent

import (
	"fmt"
)

type CollisionError struct {
	CollidedWithPlayerID *UUID
	collided_with_bounds bool
}

func NewPlayerCollisionError(collided_with_player_id UUID) error {
	return CollisionError{&collided_with_player_id, false}
}

func NewBoundsCollisionError() error {
	return CollisionError{nil, true}
}

func (err CollisionError) Error() string {
	if err.CollidedWithPlayerID != nil {
		return fmt.Sprintf("Collided with player ID '%s'.", err.CollidedWithPlayerID)
	}
	if err.collided_with_bounds {
		return "Collided with bounds."
	}
	return "Unknown collision."
}

func (err *CollisionError) CollidedWithPlayer() bool {
	return (err.CollidedWithPlayerID != nil)
}

func (err *CollisionError) CollidedWithBounds() bool {
	return err.collided_with_bounds
}
