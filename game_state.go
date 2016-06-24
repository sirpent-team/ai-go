package sirpent

import (
	"encoding/json"
	"fmt"
	//  "errors"
)

type TickID uint

type GameState struct {
	// Tick number this is.
	ID TickID
	// Player (and thus Snake) states computed this tick.
	Plays map[UUID]*PlayerState
	// Food states computed this tick.
	Food HexVector
	//Foods map[FoodID]FoodTick
}

func NewGameState(game *Game) *GameState {
	return &GameState{
		ID:    game.TickCount,
		Plays: make(map[UUID]*PlayerState),
		Food:  game.Grid.CryptoRandomCell(),
	}
}

func (gs *GameState) Successor(game *Game) *GameState {
	gs2 := &GameState{
		ID:    game.TickCount,
		Plays: make(map[UUID]*PlayerState),
		Food:  gs.Food,
	}

	for player_id, player_state := range gs.Plays {
		if !player_state.Alive {
			continue
		}
		player := player_state.Player
		next_player_state, err := gs.Plays[player_id].Successor(gs)
		if err != nil {
			fmt.Printf("ERROR: player id %s, player %+v, play %+v, err %s\n", player_id, player, next_player_state, err)
		}
		if next_player_state.Snake[0] == gs2.Food {
			previous_snake := gs.Plays[player_id].Snake
			next_player_state.Snake = append(next_player_state.Snake, previous_snake[len(previous_snake)-1])
			gs2.Food = game.Grid.CryptoRandomCell()
		}
		gs2.Plays[player_id] = next_player_state
	}

	gs2.CollisionDetection(game)

	for player_id, _ := range gs.Plays {
		gs2.Plays[player_id].Alive = game.Players[player_id].Alive
	}

	return gs2
}

func (gs *GameState) CollisionDetection(game *Game) {
	// @TODO: Goroutine this, considering time to compare vs overhead.
	for player1_id, player1_state := range gs.Plays {
		if !game.Grid.IsWithinBounds(player1_state.Snake[0]) {
			fmt.Printf("%s outside bounds\n", player1_id)
			player1_state.Player.Alive = false
			player1_state.Player.DiedFrom.HandleError(CollidedWithBoundsError{})
		}

		for player2_id, player2_state := range gs.Plays {
			// Loop covers self-intersection.
			if player1_state.Snake.HeadIntersects(player2_state.Snake) {
				fmt.Printf("%s head inside %s\n", player1_id, player2_id)
				player1_state.Player.Alive = false
				player1_state.Player.DiedFrom.HandleError(CollidedWithPlayerError{CollidedWithPlayerID: player2_id})
			}
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
		ID    uint
		Plays map[string]PlayerState
		Food  HexVector
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
		ID    uint
		Plays map[string]PlayerState
		Food  HexVector
	}{}

	err := json.Unmarshal(b, &gs_for_json)
	if err == nil {
		gs.ID = TickID(gs_for_json.ID)
		gs.Plays = make(map[UUID]*PlayerState, len(gs_for_json.Plays))
		for player_id_str, player_state := range gs_for_json.Plays {
			var player_id UUID
			player_id, err = UUIDFromString(player_id_str)
			// Copy memory out of
			gs.Plays[player_id] = new(PlayerState)
			*gs.Plays[player_id] = player_state
		}
		gs.Food = gs_for_json.Food
	}

	return err
}
