package sirpent

import (
	"net"
)

// Useful for debugging and statistics.
type CauseOfDeath struct {
	// Death error.
	Err error `json:"err"`
	// Tick that the player died in.
	// @TODO: Implement setting this.
	//DiedAt TickID `json:"died_at"`
	// Did the player died because the connection to their server dropped or timed out?
	SocketProblem bool `json:"socket_problem"`
	SocketTimeout bool `json:"socket_timeout"`
	// Did the player die because they specified an unrecognised direction?
	InvalidMove bool `json:"invalid_move"`
	// Did the player die because they collided with a Player (including themself)? Or nil.
	CollisionWithPlayerID *UUID `json:"collision_with_player"`
	// Did the player die because they went beyond the boundaries of the world?
	CollidedWithBounds bool `json:"collided_with_bounds"`
}

func (cod *CauseOfDeath) HandleError(err error) {
	cod.Err = err
	if neterr, ok := err.(net.Error); ok {
		if neterr.Timeout() {
			cod.SocketTimeout = true
		} else {
			cod.SocketProblem = true
		}
	} else if _, ok := err.(DirectionError); ok {
		cod.InvalidMove = true
	} else if playercollisionerr, ok := err.(CollidedWithPlayerError); ok {
		cod.CollisionWithPlayerID = &playercollisionerr.CollidedWithPlayerID
	} else if _, ok := err.(CollidedWithBoundsError); ok {
		cod.CollidedWithBounds = true
	}
}
