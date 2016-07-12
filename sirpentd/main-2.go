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

  // Connect to all players.
  var ended_wg sync.WaitGroup
  for _, player := range(game.Players) {
    ended_wg.Add(1)
    go player.Connect(game, 5 * time.Second, &ended_wg)
  }
  ended_wg.Wait()

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

  // @TODO: Ensure all player connections closed.

  ////////////////////////////////////////////////////////////////////////
  ////////////////////////////////////////////////////////////////////////
  ////////////////////////////////////////////////////////////////////////
  ////////////////////////////////////////////////////////////////////////
  ////////////////////////////////////////////////////////////////////////

  /*
  grid := sirpent.NewHexHexGrid(20)

  players := make(map[sirpent.UUID]*sirpent.Player)
  player0 := sirpent.NewPlayer("localhost:8901")
  players[player0.ID] = player0

  game := sirpent.NewGame(grid, players)

  // Connect to players.
  // @TODO: Prevent waiting for network N times, using sync.WaitGroup and Goroutines.
  for _, player := range game.Players {
    err := player.Connect(game)
    if err != nil {
      // @TODO: Decide how to handle connection unestablished.
      fmt.Println("Player connection failed.")
      panic(err)
    }
  }
  // @TODO: Tell players about game, grid etc!

  // @TODO: Basic setup complete; start API server.
  //        Expand upon api.go
  api := sirpent.API{}

  go func() {
    http.Handle("/", http.FileServer(http.Dir("webroot")))
    http.Handle("/worlds/live.json", websocket.Handler(func(ws *websocket.Conn) {
      defer func() {
        err := ws.Close()
        if err != nil {
          fmt.Printf("websocket.Close err=%s\n", err)
        }
      }()

      // @TODO: Error handling.
      _ = websocket.JSON.Send(ws, game.Grid)
      api.Websockets = append(api.Websockets, ws)

      // @TODO: Keep Websocket alive without an infinite loop.
      // Use channels properly?
      for {
        time.Sleep(1 * time.Second)
      }
    }))
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
      panic("ListenAndServe: " + err.Error())
    }
  }()

  // Begin the game.
  for {
    fmt.Printf("Tick %d\n", game.TickCount)
    latest_state := game.Tick()

    for i := range api.Websockets {
      func(websocket_index int) {
        api.Websockets[websocket_index].SetDeadline(time.Now().Add(5 * time.Second))
        err := websocket.JSON.Send(api.Websockets[websocket_index], latest_state)
        if err != nil {
          fmt.Printf("%+v\n", err)
          // @TODO: Locking needed?
          api.Websockets = append(api.Websockets[:websocket_index], api.Websockets[websocket_index+1:]...)
        }
      }(i)
    }

    if !latest_state.HasLivingPlayers() {
      fmt.Printf("(NoLivingPlayers)\n")
      break
    }

    time.Sleep(75 * time.Millisecond)
  }
  */
}
