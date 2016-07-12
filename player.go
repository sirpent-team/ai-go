package sirpent

import (
	"fmt"
	"sync"
	"time"
)

type Player struct {
	ID UUID `json:"id"`
	// Address to open a TCP socket to.
	server_address string
	connection     *jsonSocket
	ready_wg       *sync.WaitGroup
	// Is the Player alive after the most recent tick?
	Alive bool `json:"alive"`
	// What killed the player?
	DiedFrom CauseOfDeath `json:"died_from"`
}

func NewPlayer(server_address string) *Player {
	// @TODO: Ensure CauseOfDeath will deny every cause of death unless modified.
	return &Player{
		ID:             NewUUID(),
		server_address: server_address,
		connection:     nil,
		Alive:          true,
	}
}

func (p *Player) Connect(game *Game, player_connection_timeout time.Duration, ended_wg *sync.WaitGroup) {
	// 1. Simultaneously all Players connect and send the player ID.
	connection, err := newJsonSocket(p.server_address, player_connection_timeout)
	if err != nil {
		p.errorKillPlayer(err)
		ended_wg.Done()
		return
	}
	p.connection = connection

	err = p.connection.sendOrTimeout(p.ID)
	if err != nil {
		p.errorKillPlayer(err)
		ended_wg.Done()
		return
	}

	err = p.connection.sendOrTimeout(game)
	if err != nil {
		p.errorKillPlayer(err)
		ended_wg.Done()
		return
	}

	ended_wg.Done()
}

func (p *Player) PlayTurn(game *Game, ended_wg *sync.WaitGroup) {
	if p.Alive {
		previous_game_state := game.Ticks[game.TickCount-2]
		next_game_state := game.Ticks[game.TickCount-1]

		// Player turn loop:
		// 1. Send current game state.
		// 2. Receive chosen move.
		// 3. Update player state.
		// 4. Wait for global turn operations.
		// 5. Go to 1 unless player is dead.

		// 1. Send current game state.
		err := p.connection.sendOrTimeout(previous_game_state)
		if err != nil {
			p.errorKillPlayer(err)
			ended_wg.Done()
			return
		}

		// 2. Receive chosen move.
		var direction Direction
		err = p.connection.receiveOrTimeout(&direction)
		if err == nil {
			err = game.Grid.ValidateDirection(direction)
		}
		if err != nil {
			p.errorKillPlayer(err)
			ended_wg.Done()
			return
		}

		// 3. Update player state.
		previous_player_state := previous_game_state.Plays[p.ID]
		next_player_state := &PlayerState{
			Player: p,
			Move:   direction,
			Alive:  p.Alive,
			Snake:  previous_player_state.Snake.Move(game.Grid, direction),
		}
		next_game_state.Plays[p.ID] = next_player_state
	}

	ended_wg.Done()
}

func (p *Player) errorKillPlayer(err error) {
	p.Alive = false
	p.DiedFrom.DiagnoseError(err)
	fmt.Printf("---\nDIED: Player %s died from %s---\n", p.ID, p.DiedFrom.Spew())
}
