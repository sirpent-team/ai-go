package sirpent

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"net"
	"time"
)

type Player struct {
	ID uuid.UUID
	// Address to open a TCP socket to.
	ServerAddress string
	// Opened socket.
	ServerSocket net.Conn
	// Bufio readers and writers.
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
	// Is the Player alive after the most recent tick?
	Alive bool
	// What killed the player?
	DiedFrom CauseOfDeath
}

// Useful for debugging and statistics.
type CauseOfDeath struct {
	// Tick that the player died in.
	TimeOfDeath TickID
	// Did the player died because the connection to their server dropped or timed out?
	SocketDisconnected bool
	SocketTimeout      bool
	// Did the player die because they specified an unrecognised direction?
	InvalidMove       bool
	InvalidMoveString string
	// Did the player die because their head moved into the segment behind it?
	// @TODO: Does this case really need a separate flag to CollisionWithPlayerID?
	//TurnedBackOnTail bool
	// Did the player die because they collided with a Player (including themself)? Or nil.
	// @TODO: Maybe remove TurnedBackOnTail in favour of inferring from this?
	//        Alternatively store the syndrome of the collision.
	CollisionWithPlayerID uuid.UUID
	// Did the player die because they went beyond the boundaries of the world?
	CollisionWithBounds bool
}

func NewPlayer(server_address string) *Player {
	// @TODO: Ensure CauseOfDeath will deny every cause of death unless modified.
	return &Player{ID: NewUUID(), ServerAddress: server_address, Alive: true}
}

func (p *Player) ConnectToPlayer() error {
	c, err := net.DialTimeout("tcp", p.ServerAddress, time.Duration(5*time.Second))
	if err == nil {
		p.ServerSocket = c
		p.Scanner = bufio.NewScanner(p.ServerSocket)
		p.Writer = bufio.NewWriter(p.ServerSocket)
	}
	return err
}

/*func (p Player) SendWorld(w *World) error {
  fmt.Fprintf(p.Connection, "SENDING_WORLD\r\n")

  return nil//err
}*/

func (p *Player) RequestMove() (Direction, error) {
	p.ServerSocket.SetDeadline(time.Now().Add(5 * time.Second))
	// @TODO: Need a communication protocol.
	fmt.Fprintf(p.ServerSocket, "MOVE PLS\n")
	if !p.Scanner.Scan() {
		err := p.Scanner.Err()
		// @TODO: If a socket times out partway through receiving a direction message,
		// this will presently lead to an InvalidMove error. Rework the communication code
		// to avoid.
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			p.Alive = false
			p.DiedFrom.SocketTimeout = true
		} else {
			// @TODO: Unsure if a non-Timeout error strictly indicates the socket disconnected.
			// But it may indicate recalling Scan won't help... maybe recode to be sure?
			// However this handles the case when Scan() returns false to indicate EOF,
			// in which case Err() is nil.
			if err == nil {
				err = errors.New("EOF received from Player socket.")
			}
			p.Alive = false
			p.DiedFrom.SocketDisconnected = true
		}
		return NorthEast, err
	}
	direction_msg := p.Scanner.Text()
	p.ServerSocket.SetDeadline(time.Time{})

	direction, err := DirectionByString(direction_msg)
	if err != nil {
		p.Alive = false
		p.DiedFrom.InvalidMove = true
		p.DiedFrom.InvalidMoveString = direction_msg
	}
	return direction, err
}

func (p *Player) NewPlayerState(snake Snake) *PlayerState {
	return &PlayerState{
		Player:       p,
		Alive:        true,
		CurrentSnake: snake,
	}
}

type PlayerState struct {
	Player *Player
	// Player-chosen Direction this tick.
	ChosenMove Direction
	// Player alive after this tick?
	Alive bool
	// The current state of the snake.
	CurrentSnake Snake
	// Food item eaten this tick.
	//EatenFood *Food
}

func (ps PlayerState) Successor() (*PlayerState, error) {
	direction, err := ps.Player.RequestMove()

	// Get Alive from the Player, in case the player was killed by RequestMove/etc.
	ps2 := &PlayerState{
		Player:       ps.Player,
		ChosenMove:   direction,
		Alive:        ps.Player.Alive,
		CurrentSnake: ps.CurrentSnake,
	}

	// If there was an error in RequestMove or the player isn't alive (which often go together)
	// then no further change of state needed.
	if err != nil || !ps2.Alive {
		return ps2, err
	}

	// Update snake position.
	// @TODO: Where should StepInDirection grow_extra_segment argument come in?
	ps2.CurrentSnake.StepInDirection(ps2.ChosenMove, false)

	return ps2, err
}

func (ps PlayerState) MarshalJSON() ([]byte, error) {
	ps_for_json := struct {
		PlayerID     string
		ChosenMove   string
		Alive        bool
		CurrentSnake Snake
	}{
		PlayerID:     ps.Player.ID.String(),
		ChosenMove:   ps.ChosenMove.String(),
		Alive:        ps.Alive,
		CurrentSnake: ps.CurrentSnake,
	}

	return json.Marshal(ps_for_json)
}
