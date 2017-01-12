package sirpent

import "encoding/json"

type MsgKind int

const (
	Version MsgKind = iota
	Register
	Welcome
	NewGame
	Turn
	Move
	Died
	Won
	GameOver
)

// http://eagain.net/articles/go-json-kind/
var MsgKindHandlers = map[MsgKind]func() interface{}{
	Version:  func() interface{} { return &VersionMsg{} },
	Register: func() interface{} { return &RegisterMsg{} },
	Welcome:  func() interface{} { return &WelcomeMsg{} },
	NewGame:  func() interface{} { return &NewGameMsg{} },
	Turn:     func() interface{} { return &TurnMsg{} },
	Move:     func() interface{} { return &MoveMsg{} },
	Died:     func() interface{} { return &DiedMsg{} },
	Won:      func() interface{} { return &WonMsg{} },
	GameOver: func() interface{} { return &GameOverMsg{} },
}

type Msg struct {
	Msg  MsgKind     `json:"msg"`
	Data interface{} `json:"data"`
}

type VersionMsg struct {
	Sirpent  string `json:"sirpent"`
	Protocol string `json:"protocol"`
}

type RegisterMsg struct {
	DesiredName string `json:"desired_name"`
	Kind        string `json:"kind"`
}

type Timeout struct {
	Nanos int `json:"nanos"`
	Secs  int `json:"secs"`
}

type WelcomeMsg struct {
	Name    string        `json:"name"`
	Grid    HexagonalGrid `json:"grid"`
	Timeout Timeout       `json:"timeout"`
}

type NewGameMsg struct {
	Game GameState `json:"game"`
}

type TurnMsg struct {
	Turn TurnState `json:"turn"`
}

type MoveMsg struct {
	Direction Direction `json:"direction"`
}

type DiedMsg struct {
	CauseOfDeath json.RawMessage `json:"cause_of_death"`
}

type WonMsg struct{}

type GameOverMsg struct {
	Turn TurnState `json:"turn"`
}
