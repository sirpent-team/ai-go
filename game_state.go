package sirpent

type TickID uint

type GameState struct {
	// Tick number this is.
	ID TickID `json:"id"`
	// Player (and thus Snake) states computed this tick.
	Plays map[string]*PlayerState `json:"plays"`
	// Food states computed this tick.
	Food Vector `json:"food"`
	//Foods map[FoodID]FoodTick
}

func NewGameState(game *Game) *GameState {
	food, _ := game.World.Grid.CryptoRandomCell()
	return &GameState{
		ID:    game.TickCount,
		Plays: make(map[string]*PlayerState),
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
