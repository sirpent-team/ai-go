package sirpent

type PlayerState struct {
	Name   string        `json:"player_id"`
	Action *PlayerAction `json:"action"`
	// Player alive after this tick?
	Alive bool `json:"alive"`
	// The current state of the snake.
	Snake Snake `json:"snake"`
}

func NewPlayerState(p *Player, s Snake) *PlayerState {
	return &PlayerState{
		Name:  p.Name,
		Alive: p.Alive,
		Snake: s,
	}
}
