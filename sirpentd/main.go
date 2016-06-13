package main

import (
  "fmt"
  "time"
//  "sync"
  "net/http"
  "golang.org/x/net/websocket"
  "github.com/Taneb/sirpent"
)

func main() {
  //var a sirpent.Vector
  //a = sirpent.HexagonalVector{X: 5, Y: 2}
  //w := sirpent.World{G: sirpent.HexagonalGrid{Width: 30, Height: 30}, Players: make([]*sirpent.Player, 0)}

  w := new(sirpent.World)
  w.G = sirpent.HexagonalGrid{Width: 30, Height: 30}
  w.Players = make([]*sirpent.Player, 0)

  p := sirpent.NewPlayer("localhost:8901")
  p.S = sirpent.NewSnake(sirpent.HexagonalVector{X: 5, Y: 5})
  //p.S.Segments = append(p.S.Segments, sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 5, Y: 4}})
  err := p.ConnectToPlayer()
  if err == nil {
    w.Players = append(w.Players, p)
  }

  /*w.Snakes[0] = *sirpent.NewSnake(sirpent.HexagonalVector{X: 1, Y: 2})
  fmt.Println(w)
  w.Snakes[0].StepInDirection(sirpent.NORTHEAST)
  fmt.Println(w)
  w.Snakes[0].StepInDirection(sirpent.NORTHEAST)
  fmt.Println(w)
  w.Snakes[0].StepInDirection(sirpent.SOUTH)
  fmt.Println(w)
  w.Snakes[0].StepInDirection(sirpent.SOUTHEAST)
  fmt.Println(w)*/
  go ws(w)

  for {
    tick(w)
    time.Sleep(1000 * time.Millisecond)
  }
}

func ws(w *sirpent.World) {
  http.Handle("/worlds/live", websocket.Handler(func(ws *websocket.Conn) {
    w.Websockets = append(w.Websockets, ws)
    for {}
  }))
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    panic("ListenAndServe: " + err.Error())
  }
}

func tick(w *sirpent.World) {
  fmt.Println("Tick")

  //var wg sync.WaitGroup
  for i := range w.Players {
    player_tick(w, w.Players[i])
  }

  // Detect collisions.
  // @TODO: Wait for all players to have moved.
  /*for i := range w.Players {
    for j := range w.Players {
      //has_collided := w.Players[i].S.HasCollidedIntoSnake(*w.Players[j].S)
      //head_on_collided := w.Players[i].S.IsHeadAt(w.Players[j].S.Segments[0].Position)
    }
  }*/

  // Notify dead players. (Close connections.)

  // Broadcast new world on websocket.
  for i := range w.Websockets {
    //fmt.Fprintf(w.Websockets[i], "Test %s\r\n", "abc")
    for j := range w.Players {
      //fmt.Fprintf(w.Websockets[i], w.Players[j].S, "")
      websocket.JSON.Send(w.Websockets[i], w.Players[j].S)
    }
    //fmt.Fprintf(w.Websockets[i], "Test %s\r\n", "abc")
  }
}

func player_tick(w *sirpent.World, p *sirpent.Player) {
  // Send World to player as JSON.
  p.SendWorld(w)

  // As N goroutines, get move from all live players.
  direction, err := p.ReceiveMove()
  if err != nil {
    // Kill the player.
    fmt.Printf("ERROR: %s\n", err)
  } else {
    // Apply move.
    fmt.Printf("%s -> ", p.S.Segments[0])
    p.S.StepInDirection(direction)
    fmt.Printf("%s\n", p.S.Segments[0])
  }
}

/*
messages := make(chan int)
    var wg sync.WaitGroup

    // you can also add these one at
    // a time if you need to

    wg.Add(3)
    go func() {
        defer wg.Done()
        time.Sleep(time.Second * 3)
        messages <- 1
    }()
    go func() {
        defer wg.Done()
        time.Sleep(time.Second * 2)
        messages <- 2
    }()
    go func() {
        defer wg.Done()
        time.Sleep(time.Second * 1)
        messages <- 3
    }()
    go func() {
        for i := range messages {
            fmt.Println(i)
        }
    }()

    wg.Wait()
*/
