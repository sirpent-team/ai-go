package main

import (
	"fmt"
	"time"
	// "sync"
	"github.com/Taneb/sirpent"
	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"net/http"
)

func main() {
	grid := &sirpent.Grid{Radius: 30, Origin: sirpent.AxialVector{Q: 0, R: 0}}

	players := make(map[uuid.UUID]*sirpent.Player)
	player0 := sirpent.NewPlayer("localhost:8901")
	players[player0.ID] = player0
	player1 := sirpent.NewPlayer("localhost:8902")
	players[player1.ID] = player1

	game := sirpent.NewGame(grid, players)

	// Connect to players.
	// @TODO: Prevent waiting for network N times, using sync.WaitGroup and Goroutines.
	for _, player := range game.Players {
		err := player.ConnectToPlayer()
		if err != nil {
			// @TODO: Decide how to handle connection unestablished.
			panic("Player connection failed.")
		}
	}
	// @TODO: Tell players about game, grid etc!

	// @TODO: Basic setup complete; start API server.
	//        Expand upon api.go
	api := sirpent.API{}

	go func() {
		http.Handle("/worlds/live", websocket.Handler(func(ws *websocket.Conn) {
			api.Websockets = append(api.Websockets, ws)
			// @TODO: Keep Websocket alive without an infinite loop.
			// Use channels properly?
			for {
			}
		}))
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()

	// Begin the game.
	for {
		//if game.TickCount%1000 == 0 {
		fmt.Printf("Tick %d\n", game.TickCount)
		//}

		latest_state := game.Tick()
		// @TODO: Dead players should stop being ticked, and thus need to iterate over plays not players here.
		/*for player_id, player := range game.Players {
			fmt.Printf("player id %s, current snake %+v\n", player_id, latest_state.Plays[player_id].CurrentSnake)
			fmt.Printf("player %+v, player state %+v\n", player, latest_state.Plays[player_id])
		}
		fmt.Printf("%+v\n", latest_state)*/

		for i := range api.Websockets {
			err := websocket.JSON.Send(api.Websockets[i], latest_state)
			if err != nil {
				fmt.Printf("%+v\n", err)
			}
			/*for _, player_state := range latest_state.Plays {
				fmt.Printf("%d\n", i)
				//fmt.Fprintf(w.Websockets[i], w.Players[j].S, "")
				websocket.JSON.Send(api.Websockets[i], player_state.CurrentSnake)
			}*/
			//fmt.Fprintf(w.Websockets[i], "Test %s\r\n", "abc")
		}

		if !latest_state.HasLivingPlayers() {
			break
		}

		time.Sleep(300 * time.Millisecond)
	}
}
