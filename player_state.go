package sirpent

import (
	//"errors"
	"fmt"
)

type PlayerState struct {
	Player *Player `json:"player"`
	// Player-chosen Direction this tick.
	Move Direction `json:"move"`
	// Player alive after this tick?
	Alive bool `json:"alive"`
	// The current state of the snake.
	Snake Snake `json:"snake"`
	// Food item eaten this tick.
	//EatenFood *Food
}

func NewPlayerState(p *Player, s Snake) *PlayerState {
	return &PlayerState{
		Player: p,
		Alive:  true,
		Snake:  s,
	}
}

func (ps PlayerState) Successor(gs *GameState) (*PlayerState, error) {
	direction, err := ps.Player.requestMove(gs)

	// Get Alive from the Player, in case the player was killed by RequestMove/etc.
	ps2 := &PlayerState{
		Player: ps.Player,
		Move:   direction,
		Alive:  ps.Player.Alive,
	}
	fmt.Printf("player id=%s direction=%s previous snake=%+v ", ps2.Player.ID, ps2.Move, ps.Snake)

	// If there was an error in RequestMove or the player isn't alive (which often go together)
	// then no further change of state needed.
	if err != nil || !ps2.Alive {
		return ps2, err
	}

	// Update snake position.
	// @TODO: How should StepInDirection grow_extra_segment argument come in?
	ps2.Snake = ps.Snake.Move(ps2.Move, false)
	fmt.Printf("next snake=%+v\n", ps2.Snake)

	return ps2, err
}
