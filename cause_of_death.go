package sirpent

import (
	"github.com/davecgh/go-spew/spew"
	"net"
)

// Useful for debugging and statistics.
type CauseOfDeath struct {
	// Death error.
	Err string `json:"error"`
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

func (cod *CauseOfDeath) DiagnoseError(err error) {
	cod.Err = err.Error()
	if net_err, ok := err.(net.Error); ok {
		if net_err.Timeout() {
			cod.SocketTimeout = true
		} else {
			cod.SocketProblem = true
		}
	} else if _, ok := err.(DirectionError); ok {
		cod.InvalidMove = true
	} else if collision_err, ok := err.(CollisionError); ok {
		if collision_err.CollidedWithBounds() {
			cod.CollidedWithBounds = true
		}
		if collision_err.CollidedWithPlayer() {
			cod.CollisionWithPlayerID = collision_err.CollidedWithPlayerID
		}
	}
}

// Spew is useful to neatly present Cause Of Death. But it cannot be the String() method.
// https://github.com/davecgh/go-spew/issues/45
func (cod CauseOfDeath) Spew() string {
	return spew.Sdump(cod)
}
