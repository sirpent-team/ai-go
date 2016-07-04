package sirpent

import (
//	"errors"
)

type Game struct {
	ID   UUID    `json:"id"`
	Grid HexGrid `json:"grid"`
	// All players in this game.
	Players map[UUID]*Player `json:"players"`
	// Most recent tick number.
	TickCount TickID `json:"tick_count"`
	// Game state for all (loaded) ticks.
	Ticks map[TickID]*GameState `json:"ticks"`
}

func NewGame(grid HexGrid, players map[UUID]*Player) *Game {
	game := &Game{
		ID:        NewUUID(),
		Grid:      grid,
		Players:   players,
		TickCount: TickID(0),
		Ticks:     make(map[TickID]*GameState),
	}

	// Create first Tick, for starting positions etc.
	gs := NewGameState(game)
	// Generate player starting positions.
	for player_id, player := range game.Players {
		// @TODO: Ensure player starting positions do not intersect.
		snake := NewSnake(game.Grid.CryptoRandomCell())
		snake = append(snake, snake[0].Neighbour(SouthWest))
		snake = append(snake, snake[1].Neighbour(SouthWest))
		gs.Plays[player_id] = NewPlayerState(player, snake)
	}
	game.Ticks[game.TickCount] = gs
	game.TickCount++

	return game
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
