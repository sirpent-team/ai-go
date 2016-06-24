package sirpent

import (
	"time"
)

type Player struct {
	ID UUID
	// Address to open a TCP socket to.
	ServerAddress string
	connection    *jsonSocket
	// Is the Player alive after the most recent tick?
	Alive bool
	// What killed the player?
	DiedFrom CauseOfDeath
}

func NewPlayer(server_address string) *Player {
	// @TODO: Ensure CauseOfDeath will deny every cause of death unless modified.
	return &Player{
		ID:            NewUUID(),
		ServerAddress: server_address,
		connection:    nil,
		Alive:         true,
	}
}

func (p *Player) Connect(game *Game) error {
	connection, err := newJsonSocket(p.ServerAddress, time.Duration(5*time.Second))
	if err != nil {
		p.Alive = false
		p.DiedFrom.HandleError(err)
		return err
	}
	p.connection = connection

	// Send player ID to player.
	err = p.connection.sendOrTimeout(p.ID, 5*time.Second)
	if err != nil {
		p.Alive = false
		p.DiedFrom.HandleError(err)
		return err
	}

	// Send game grid to player.
	err = p.connection.sendOrTimeout(game.Grid, 5*time.Second)
	if err != nil {
		p.Alive = false
		p.DiedFrom.HandleError(err)
	}
	return err
}

func (p *Player) requestMove(gs *GameState) (Direction, error) {
	var direction Direction

	// @TODO: There has to be a better way to handle errors than this?
	err := p.connection.sendOrTimeout(gs, 5*time.Second)
	if err != nil {
		p.Alive = false
		p.DiedFrom.HandleError(err)
		return direction, err
	}

	err = p.connection.receiveOrTimeout(&direction, 5*time.Second)
	if err != nil {
		p.Alive = false
		p.DiedFrom.HandleError(err)
	}

	return direction, err
}
