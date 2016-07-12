package sirpent

import (
	"strconv"
	//	"errors"
	"encoding/json"
)

type Game struct {
	ID   UUID `json:"id"`
	Grid Grid `json:"grid"`
	// All players in this game.
	Players map[UUID]*Player `json:"players"`
	// Most recent tick number.
	TickCount TickID `json:"tick_count"`
	// Game state for all (loaded) ticks.
	Ticks map[TickID]*GameState `json:"ticks"`
}

func NewGame(grid Grid, players map[UUID]*Player) *Game {
	return &Game{
		ID:        NewUUID(),
		Grid:      grid,
		Players:   players,
		TickCount: TickID(0),
		Ticks:     make(map[TickID]*GameState),
	}
}

// @TODO: While NewGame creates an initial state, no need for a zero-state error.
func (game *Game) LatestTick() *GameState {
	return game.Ticks[game.TickCount-1]
}

func (game *Game) MarshalJSON() ([]byte, error) {
	game_for_json := struct {
		ID        UUID                  `json:"id"`
		Grid      Grid                  `json:"grid"`
		Players   map[string]*Player    `json:"players"`
		TickCount TickID                `json:"tick_count"`
		Ticks     map[string]*GameState `json:"ticks"`
	}{
		ID:        game.ID,
		Grid:      game.Grid,
		Players:   make(map[string]*Player),
		TickCount: game.TickCount,
		Ticks:     make(map[string]*GameState),
	}

	for player_id, player := range game.Players {
		game_for_json.Players[player_id.String()] = player
	}
	for tick_id, tick := range game.Ticks {
		game_for_json.Ticks[strconv.Itoa(int(tick_id))] = tick
	}

	return json.Marshal(game_for_json)
}

func (game *Game) UnmarshalJSON(b []byte) error {
	game_for_json := struct {
		ID        UUID                  `json:"id"`
		Grid      json.RawMessage       `json:"grid"`
		Players   map[string]*Player    `json:"players"`
		TickCount TickID                `json:"tick_count"`
		Ticks     map[string]*GameState `json:"ticks"`
	}{}
	err := json.Unmarshal(b, &game_for_json)
	if err != nil {
		return err
	}

	game.ID = game_for_json.ID

	game.Grid, err = ParseGridJSON(game_for_json.Grid)
	if err != nil {
		return err
	}

	game.Players = make(map[UUID]*Player, len(game_for_json.Players))
	for player_id_str, player := range game_for_json.Players {
		var player_id UUID
		player_id, err = UUIDFromString(player_id_str)
		game.Players[player_id] = player
	}

	game.TickCount = game_for_json.TickCount

	game.Ticks = make(map[TickID]*GameState, len(game_for_json.Ticks))
	for tick_id_str, tick := range game_for_json.Ticks {
		tick_id, err := strconv.ParseUint(tick_id_str, 10, 64)
		if err != nil {
			return err
		}
		game.Ticks[TickID(tick_id)] = tick
	}

	return err
}
