package sirpent

import (
	"encoding/json"
	"fmt"
	//	"errors"
	"github.com/satori/go.uuid"
)

type Game struct {
	ID uuid.UUID
	// Unrendered grid.
	Grid *Grid
	// All players in this game.
	Players map[uuid.UUID]*Player
	// Most recent tick number.
	TickCount TickID
	// Game state for all (loaded) ticks.
	Ticks map[TickID]*GameState
}

func NewGame(grid *Grid, players map[uuid.UUID]*Player) *Game {
	game := &Game{
		ID:        NewUUID(),
		Grid:      grid,
		Players:   players,
		TickCount: TickID(0),
		Ticks:     make(map[TickID]*GameState),
	}

	// Create first Tick, for starting positions etc.
	gs := game.NewGameState()
	// Generate player starting positions.
	for player_id, player := range game.Players {
		// @TODO: Ensure player starting positions do not intersect.
		snake := NewSnake(game.Grid.CryptoRandomCell())
		gs.Plays[player_id] = player.NewPlayerState(snake)
	}
	game.Ticks[game.TickCount] = gs
	game.TickCount++

	return game
}

func (game *Game) NewGameState() *GameState {
	return &GameState{
		ID: game.TickCount,
		//CurrentGrid: *game.Grid,
		Plays: make(map[uuid.UUID]*PlayerState),
	}
}

// @TODO: While NewGame creates an initial state, no need for a zero-state error.
func (game *Game) LatestState() *GameState {
	return game.Ticks[game.TickCount-1]
}

func (game *Game) Tick() *GameState {
	gs := game.Ticks[game.TickCount-1].Successor(game)
	game.Ticks[game.TickCount] = gs
	game.TickCount++
	return gs
}

type TickID uint

type GameState struct {
	// Tick number this is.
	ID TickID
	// Rendered grid computing this tick.
	RenderedGrid Grid
	// Player (and thus Snake) states computed this tick.
	Plays map[uuid.UUID]*PlayerState
	// Food states computed this tick.
	//Foods map[FoodID]FoodTick
}

func (gs *GameState) Successor(game *Game) *GameState {
	gs2 := &GameState{
		ID: game.TickCount,
		//RenderedGrid: *game.Grid,
		Plays: make(map[uuid.UUID]*PlayerState),
	}

	for player_id, player_state := range gs.Plays {
		if !player_state.Alive {
			continue
		}
		player := player_state.Player
		next_player_state, err := gs.Plays[player_id].Successor()
		if err != nil {
			fmt.Printf("ERROR: player id %s, player %+v, play %+v, err %s\n", player_id, player, next_player_state, err)
		}
		gs2.Plays[player_id] = next_player_state
	}

	gs2.CollisionDetection(game)

	return gs2
}

func (gs *GameState) CollisionDetection(game *Game) {
	// @TODO: Goroutine this, considering time to compare vs overhead.

	for player_id, player_state := range gs.Plays {
		func() {
			if !game.Grid.IsWithinBounds(player_state.CurrentSnake.Segments[0]) {
				fmt.Printf("%s outside bounds\n", player_state.Player.ID)
				player_state.Player.Alive = false
				player_state.Player.DiedFrom.CollisionWithBounds = true
				player_state.Alive = false
			}

			if player_state.CurrentSnake.IsHeadInsideTail() {
				fmt.Printf("%s headinsidetail\n", player_state.Player.ID)
				player_state.Player.Alive = false
				player_state.Player.DiedFrom.CollisionWithPlayerID = player_id
				player_state.Alive = false
			}
		}()
	}

	for player1_id, player1_state := range gs.Plays {
		for player2_id, player2_state := range gs.Plays {
			if player1_id == player2_id {
				continue
			}

			func() {
				if player1_state.CurrentSnake.IsHeadInsideSnake(player2_state.CurrentSnake) {
					fmt.Printf("%s head inside %s\n", player1_id, player2_id)
					player1_state.Player.Alive = false
					player1_state.Player.DiedFrom.CollisionWithPlayerID = player2_id
					player1_state.Alive = false
				}
			}()
		}
	}
}

func (gs *GameState) HasLivingPlayers() bool {
	has_living_players := false
	for _, player_state := range gs.Plays {
		has_living_players = has_living_players || player_state.Alive
	}
	return has_living_players
}

func (gs GameState) MarshalJSON() ([]byte, error) {
	gs_for_json := struct {
		ID           uint
		RenderedGrid Grid
		Plays        map[string]PlayerState
	}{
		ID:           uint(gs.ID),
		RenderedGrid: gs.RenderedGrid,
		Plays:        make(map[string]PlayerState),
	}

	for player_id, player_state := range gs.Plays {
		gs_for_json.Plays[player_id.String()] = *player_state
	}

	return json.Marshal(gs_for_json)
}
