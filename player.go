package sirpent

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Player struct {
	Name string `json:"name"`
	// Address to open a TCP socket to.
	ServerAddress string `json:"server_address"`
	connection    *jsonSocket
	ready_wg      *sync.WaitGroup
	// Is the Player alive after the most recent tick?
	Alive bool `json:"alive"`
	// What killed the player?
	DiedFrom CauseOfDeath `json:"died_from"`
}

func NewPlayer(name string, server_address string) *Player {
	// @TODO: Ensure CauseOfDeath will deny every cause of death unless modified.
	return &Player{
		Name:          name,
		ServerAddress: server_address,
		connection:    nil,
		Alive:         true,
	}
}

func (p *Player) Connect(game *Game, player_connection_timeout time.Duration, err_chan chan error) {
	// 1. Simultaneously all Players connect and send the player ID.
	connection, err := newJsonSocket(p.ServerAddress, player_connection_timeout)
	if err != nil {
		err_chan <- err
		return
	}
	p.connection = connection

	err = p.connection.sendOrTimeout(p.Name)
	if err != nil {
		err_chan <- err
		return
	}

	err = p.connection.sendOrTimeout(game.World)
	if err != nil {
		err_chan <- err
		return
	}

	err = p.connection.sendOrTimeout(game.Players)
	if err != nil {
		err_chan <- err
		return
	}

	err_chan <- nil
}

func (p *Player) PlayTurn(game *Game, action_chan chan *PlayerAction, err_chan chan error) {
	if !p.Alive {
		err_chan <- errors.New("Player cannot take a turn for they are already dead.")
		return
	}

	previous_game_state := game.LatestTick()

	// Player turn loop:
	// 1. Send current game state.
	// 2. Receive chosen move.
	// 3. Update player state.
	// 4. Wait for global turn operations.
	// 5. Go to 1 unless player is dead.

	// 1. Send current game state.
	err := p.connection.sendOrTimeout(previous_game_state)
	if err != nil {
		err_chan <- err
		return
	}

	// 2. Receive chosen action.
	var action *PlayerAction
	err = p.connection.receiveOrTimeout(&action)
	if err != nil {
		err_chan <- err
		return
	}
	action_chan <- action
}

func (p *Player) ErrorKillPlayer(err error) {
	p.Alive = false
	p.DiedFrom.DiagnoseError(err)
	fmt.Printf("---\nDIED: Player %s died from %s---\n", p.Name, p.DiedFrom.Spew())
}
