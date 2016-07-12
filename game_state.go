package sirpent

import (
	"encoding/json"
)

type TickID uint

type GameState struct {
	// Tick number this is.
	ID TickID
	// Player (and thus Snake) states computed this tick.
	Plays map[UUID]*PlayerState
	// Food states computed this tick.
	Food Vector
	//Foods map[FoodID]FoodTick
}

func NewGameState(game *Game) *GameState {
	food, _ := game.Grid.CryptoRandomCell()
	return &GameState{
		ID:    game.TickCount,
		Plays: make(map[UUID]*PlayerState),
		Food:  food,
	}
}

func (gs *GameState) HasLivingPlayers() bool {
	has_living_players := false
	for _, player_state := range gs.Plays {
		has_living_players = has_living_players || player_state.Alive
	}
	return has_living_players
}

func (gs *GameState) MarshalJSON() ([]byte, error) {
	gs_for_json := struct {
		ID    uint                   `json:"id"`
		Plays map[string]PlayerState `json:"plays"`
		Food  Vector                 `json:"food"`
	}{
		ID:    uint(gs.ID),
		Plays: make(map[string]PlayerState),
		Food:  gs.Food,
	}

	for player_id, player_state := range gs.Plays {
		gs_for_json.Plays[player_id.String()] = *player_state
	}

	return json.Marshal(gs_for_json)
}

func (gs *GameState) UnmarshalJSON(b []byte) error {
	gs_for_json := struct {
		ID    uint                   `json:"id"`
		Plays map[string]PlayerState `json:"plays"`
		Food  Vector                 `json:"food"`
	}{}
	err := json.Unmarshal(b, &gs_for_json)
	if err != nil {
		return err
	}

	gs.ID = TickID(gs_for_json.ID)

	gs.Plays = make(map[UUID]*PlayerState, len(gs_for_json.Plays))
	for player_id_str, player_state := range gs_for_json.Plays {
		var player_id UUID
		player_id, err = UUIDFromString(player_id_str)
		if err != nil {
			return err
		}
		gs.Plays[player_id] = new(PlayerState)
		*gs.Plays[player_id] = player_state
	}

	gs.Food = gs_for_json.Food

	return err
}
