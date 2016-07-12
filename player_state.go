package sirpent

type PlayerState struct {
	Player *Player       `json:"player"`
	Action *PlayerAction `json:"action"`
	// Player alive after this tick?
	Alive bool `json:"alive"`
	// The current state of the snake.
	Snake Snake `json:"snake"`
}

func NewPlayerState(p *Player, s Snake) *PlayerState {
	return &PlayerState{
		Player: p,
		Alive:  p.Alive,
		Snake:  s,
	}
}
