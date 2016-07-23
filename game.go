package sirpent

type Game struct {
	ID    UUID  `json:"id"`
	World World `json:"world"`
	// All players in this game.
	Players map[string]*Player `json:"players"`
	// Most recent tick number.
	TickCount TickID `json:"tick_count"`
	// Game state for all (loaded) ticks.
	Ticks map[string]*GameState `json:"ticks"`
}

func NewGame(world World, players map[string]*Player) *Game {
	return &Game{
		ID:        NewUUID(),
		World:     world,
		Players:   players,
		TickCount: TickID(0),
		Ticks:     make(map[string]*GameState),
	}
}

// @TODO: While NewGame creates an initial state, no need for a zero-state error.
func (game *Game) LatestTick() *GameState {
	return game.Ticks[string(game.TickCount-1)]
}
