package sirpent

import (
	"encoding/json"
	//"errors"
	"fmt"
)

type PlayerState struct {
	Player *Player
	// Player-chosen Direction this tick.
	Move Direction
	// Player alive after this tick?
	Alive bool
	// The current state of the snake.
	Snake Snake
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

func (ps *PlayerState) MarshalJSON() ([]byte, error) {
	ps_for_json := struct {
		PlayerID UUID
		Move     Direction
		Alive    bool
		Snake    Snake
	}{
		PlayerID: ps.Player.ID,
		Move:     ps.Move,
		Alive:    ps.Alive,
		Snake:    ps.Snake,
	}

	return json.Marshal(ps_for_json)
}

func (ps *PlayerState) UnmarshalJSON(b []byte) error {
	ps_for_json := struct {
		PlayerID UUID
		Move     Direction
		Alive    bool
		Snake    Snake
	}{}

	err := json.Unmarshal(b, &ps_for_json)
	if err == nil {
		// @TODO: Need to document that this creates barebones player structs.
		ps.Player = &Player{ID: ps_for_json.PlayerID}
		ps.Move = ps_for_json.Move
		ps.Alive = ps_for_json.Alive
		ps.Snake = ps_for_json.Snake
	}

	return err
}
