package main

import (
  "fmt"
  "time"
  "sync"
  "github.com/sirpent-team/sirpent-go"
  // "github.com/davecgh/go-spew/spew"
  // "golang.org/x/net/websocket"
  // "net/http"
)

func main() {
  fmt.Println("main-2")

  // Configure all players.
  players := make(map[sirpent.UUID]*sirpent.Player)
  player_addresses := []string{"localhost:8901", "localhost:8902"}
  for i := range(player_addresses) {
    player := sirpent.NewPlayer(player_addresses[i])
    players[player.ID] = player
  }

  // Create game datastructures.
  grid := sirpent.HexGridHexagonal{Rings: 20}
  game := sirpent.NewGame(grid, players)
  //fmt.Printf("game = ")
  //spew.Dump(game)

  bs, _ := game.MarshalJSON()
  fmt.Printf("game json = %s\n", string(bs))

  // Connect to all players.
  var ended_wg sync.WaitGroup
  for _, player := range(game.Players) {
    ended_wg.Add(1)
    go player.Connect(game, 5 * time.Second, &ended_wg)
  }
  ended_wg.Wait()

  // 1. Syncronously a game state is created with a Snake for each connected Player.
  new_game_state := sirpent.NewGameState(game)
  for player_id, player := range game.Players {
    v, _ := game.Grid.CryptoRandomCell()
    snake := sirpent.NewSnake(v)
    snake = append(snake, game.Grid.CellNeighbour(snake[0], "SOUTHWEST"))
    snake = append(snake, game.Grid.CellNeighbour(snake[1], "SOUTHWEST"))
    new_game_state.Plays[player_id] = sirpent.NewPlayerState(player, snake)
  }
  game.Ticks[game.TickCount] = new_game_state
  game.TickCount++

  play(game)
}

func play(game *sirpent.Game) {
  current_state := game.LatestTick()
  for current_state.HasLivingPlayers() {
    cs_json, _ := current_state.MarshalJSON()
    fmt.Printf("current_state json = %s\n", string(cs_json))

    next_state := &sirpent.GameState{
      ID:    game.TickCount,
      Plays: make(map[sirpent.UUID]*sirpent.PlayerState),
      Food:  current_state.Food,
    }

    err_chans := make(map[sirpent.UUID]chan error)
    action_chans := make(map[sirpent.UUID]chan sirpent.PlayerAction)
    for player_id, current_player_state := range current_state.Plays {
      if current_player_state.Alive {
        err_chans[player_id] = make(chan error)
        action_chans[player_id] = make(chan sirpent.PlayerAction)
        go current_player_state.Player.PlayTurn(game, action_chans[player_id], err_chans[player_id])
      }
    }

    for player_id, current_player_state := range current_state.Plays {
      if current_player_state.Alive {
        select {
        case err := <- err_chans[player_id]:
          fmt.Printf("Error %s %s\n", player_id, err.Error())
          //current_state.Plays[i].ErrorKillPlayer(err)
        case action := <- action_chans[player_id]:
          // @TODO: Update player.
          fmt.Printf("Action %s %s\n", player_id, action)
        }
      }
    }

    game.Ticks[game.TickCount] = next_state
    game.TickCount++
    current_state = game.LatestTick()
  }
}

/*
  // Run all players gameloop.
  latest_game_state := game.LatestTick()
  for latest_game_state.HasLivingPlayers() {
    new_game_state = &sirpent.GameState{
      ID:    game.TickCount,
      Plays: make(map[sirpent.UUID]*sirpent.PlayerState),
      Food:  latest_game_state.Food,
    }
    game.Ticks[game.TickCount] = new_game_state
    game.TickCount++

    for _, player_state := range latest_game_state.Plays {
      if player_state.Player.Alive {
        ended_wg.Add(1)
        go player_state.Player.PlayTurn(game, &ended_wg)
      }
    }
    ended_wg.Wait()
    latest_game_state = game.LatestTick()

    gs_json, _ := latest_game_state.MarshalJSON()
    fmt.Printf("game state json = %s\n", string(gs_json))
  }
*/
