package sirpent

type PlayerAction struct {
	// Player-chosen Direction of movement this tick.
	Move Direction `json:"move"`
}
