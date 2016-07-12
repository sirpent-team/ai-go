package sirpent

type PlayerAction struct {
  // Player-chosen Direction of movement this tick.
  Move Direction `json:"move"`
}

var direction Direction
    err = p.connection.receiveOrTimeout(&direction)
    if err == nil {
      err = game.Grid.ValidateDirection(direction)
    }
    if err != nil {
      p.errorKillPlayer(err)
      ended_wg.Done()
      return
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
