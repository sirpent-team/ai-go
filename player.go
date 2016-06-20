package sirpent

import (
	"encoding/json"
	//"errors"
	"net"
	"time"
	"fmt"
)

type Player struct {
	ID UUID
	// Address to open a TCP socket to.
	ServerAddress string
	// Opened socket.
	ServerSocket net.Conn
	// Bufio readers and writers.
	SocketEncoder *json.Encoder
	SocketDecoder *json.Decoder
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
	CollisionWithPlayerID UUID
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
		p.SocketEncoder = json.NewEncoder(p.ServerSocket)
		p.SocketDecoder = json.NewDecoder(p.ServerSocket)

		// Send player ID to player.
		// @TODO: Send more information? Needs a protocol defining.
		p.ServerSocket.SetDeadline(time.Now().Add(5 * time.Second))
		err = p.SocketEncoder.Encode(p.ID)
		p.ServerSocket.SetDeadline(time.Time{})
	}

	return err
}

/*func (p Player) SendWorld(w *World) error {
  fmt.Fprintf(p.Connection, "SENDING_WORLD\r\n")

  return nil//err
}*/

func (p *Player) RequestMove(gs *GameState) (Direction, error) {
	p.ServerSocket.SetDeadline(time.Now().Add(5 * time.Second))

	err := p.SocketEncoder.Encode(gs)
	if err != nil {
		panic(err)
	}

	var direction Direction
	err = p.SocketDecoder.Decode(&direction)
	if err != nil {
		panic(err)
	}

	p.ServerSocket.SetDeadline(time.Time{})
	return direction, err

	/* if !p.Scanner.Scan() {
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

	direction, err := DirectionByString(direction_msg)
	if err != nil {
		p.Alive = false
		p.DiedFrom.InvalidMove = true
		p.DiedFrom.InvalidMoveString = direction_msg
	}
	return direction, err*/
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

func (ps PlayerState) Successor(gs *GameState) (*PlayerState, error) {
	direction, err := ps.Player.RequestMove(gs)

	// Get Alive from the Player, in case the player was killed by RequestMove/etc.
	ps2 := &PlayerState{
		Player:       ps.Player,
		ChosenMove:   direction,
		Alive:        ps.Player.Alive,
		CurrentSnake: ps.CurrentSnake,
	}
	fmt.Printf("player id=%s direction=%s previous snake=%+v ", ps.Player.ID, direction, ps.CurrentSnake)

	// If there was an error in RequestMove or the player isn't alive (which often go together)
	// then no further change of state needed.
	if err != nil || !ps2.Alive {
		return ps2, err
	}

	// Update snake position.
	// @TODO: Where should StepInDirection grow_extra_segment argument come in?
	ps2.CurrentSnake.StepInDirection(ps2.ChosenMove, false)
	fmt.Printf("next snake=%+v\n", ps2.CurrentSnake)

	return ps2, err
}

func (ps *PlayerState) MarshalJSON() ([]byte, error) {
	ps_for_json := struct {
		PlayerID     UUID
		ChosenMove   string
		Alive        bool
		CurrentSnake Snake
	}{
		PlayerID:     ps.Player.ID,
		ChosenMove:   ps.ChosenMove.String(),
		Alive:        ps.Alive,
		CurrentSnake: ps.CurrentSnake,
	}

	return json.Marshal(ps_for_json)
}

func (ps *PlayerState) UnmarshalJSON(b []byte) error {
	ps_for_json := struct {
		PlayerID     UUID
		ChosenMove   Direction
		Alive        bool
		CurrentSnake Snake
	}{}

	err := json.Unmarshal(b, &ps_for_json)
	if err == nil {
		// @TODO: Need to document that this creates barebones player structs.
		ps.Player = &Player{ID: ps_for_json.PlayerID}
		ps.ChosenMove = ps_for_json.ChosenMove
		ps.Alive = ps_for_json.Alive
		ps.CurrentSnake = ps_for_json.CurrentSnake
	}

	return err
}
