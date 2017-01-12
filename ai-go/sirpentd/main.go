package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirpent-team/sirpent-go"
	"io/ioutil"
	"sync"
	"time"
	//"golang.org/x/net/websocket"
	//"net/http"
)

/*{
	"players": {
		"$name": {
			"server_address": ""
		}
	},
	"world": {
		"grid_kind": "hex_grid_hexagonal",
		"grid": {
			"rings": 20
		},
		[...]
	},
	"rounds": 10
}*/

type game_specification struct {
	Players map[string]*sirpent.Player `json:"players"`
	World   sirpent.World              `json:"world"`
	Rounds  uint                       `json:"rounds"`
}

func main() {
	game_specification_path := flag.String("c", "-", "Path to Game Configuration JSON.")
	flag.Parse()
	fmt.Println(*game_specification_path)

	game_specification_json, err := ioutil.ReadFile(*game_specification_path)
	if err != nil {
		panic("Error reading Game Specification JSON: " + err.Error())
	}
	fmt.Println(string(game_specification_json))

	var game_spec game_specification
	err = json.Unmarshal(game_specification_json, &game_spec)
	if err != nil {
		panic("Error decoding Game Specification: " + err.Error())
	}
	fmt.Printf("%+v\n", game_spec)

	for player_name, player_spec := range game_spec.Players {
		player_spec.Alive = true
		player_spec.Name = player_name
	}

	// Create game datastructures.
	game := sirpent.NewGame(game_spec.World, game_spec.Players)

	// Connect to all players.
	err_chans := make(map[string]chan error)
	for player_name, player := range game.Players {
		err_chans[player_name] = make(chan error)
		go player.Connect(game, 5*time.Second, err_chans[player_name])
	}
	for player_name, player := range game.Players {
		err := <-err_chans[player_name]
		if err != nil {
			player.ErrorKillPlayer(err)
		}
	}

	// 1. Syncronously a game state is created with a Snake for each connected Player.
	new_game_state := sirpent.NewGameState(game)
	for player_name, player := range game.Players {
		v, _ := game.World.Grid.CryptoRandomCell()
		snake := sirpent.NewSnake(v)
		snake = append(snake, game.World.Grid.CellNeighbour(snake[0], "SOUTHWEST"))
		snake = append(snake, game.World.Grid.CellNeighbour(snake[1], "SOUTHWEST"))
		new_game_state.Plays[player_name] = sirpent.NewPlayerState(player, snake)
	}
	game.Ticks[string(game.TickCount)] = new_game_state
	game.TickCount++

	var wg sync.WaitGroup

	wg.Add(1) //2)
	var game_state_chan chan *sirpent.GameState
	go run_play(game, &wg, game_state_chan)
	//go run_api(game, &sirpent.API{}, &wg, game_state_chan)
	wg.Wait()
}

func run_play(game *sirpent.Game, wg *sync.WaitGroup, game_state_chan chan *sirpent.GameState) {
	current_state := game.LatestTick()
	//game_state_chan <- current_state

	for current_state.HasLivingPlayers() {
		next_state := &sirpent.GameState{
			ID:    game.TickCount,
			Plays: make(map[string]*sirpent.PlayerState),
			Food:  current_state.Food,
		}

		err_chans := make(map[string]chan error)
		action_chans := make(map[string]chan *sirpent.PlayerAction)
		for player_name, current_player_state := range current_state.Plays {
			player := game.Players[player_name]
			if current_player_state.Alive {
				err_chans[player_name] = make(chan error)
				action_chans[player_name] = make(chan *sirpent.PlayerAction)
				go player.PlayTurn(game, action_chans[player_name], err_chans[player_name])
			}
		}

		for player_name, current_player_state := range current_state.Plays {
			player := game.Players[player_name]

			if current_player_state.Alive {
				next_snake := current_player_state.Snake

				var err error
				var action *sirpent.PlayerAction
				select {
				case err = <-err_chans[player_name]:
					fmt.Printf("Error %s %s\n", player_name, err.Error())
					player.ErrorKillPlayer(err)
				case action = <-action_chans[player_name]:
					fmt.Printf("Action %s %s\n", player_name, action)
					next_snake = next_snake.Move(game.World.Grid, action.Move)
				}

				next_player_state := sirpent.NewPlayerState(player, next_snake)
				next_player_state.Action = action
				next_state.Plays[player_name] = next_player_state

				nps_json, _ := json.Marshal(next_player_state)
				fmt.Println(string(nps_json))
			}
		}

		game.Ticks[string(game.TickCount)] = next_state
		game.TickCount++
		current_state = game.LatestTick()
		//game_state_chan <- current_state

		cs_json, _ := json.Marshal(current_state)
		fmt.Println(string(cs_json))
	}
	wg.Done()
}

/*func run_api(game *sirpent.Game, api *sirpent.API, wg *sync.WaitGroup, game_state_chan chan *sirpent.GameState) {
	http.Handle("/", http.FileServer(http.Dir("webroot")))
	http.Handle("/worlds/live.json", websocket.Handler(func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				fmt.Printf("websocket.Close err=%s\n", err)
			}
		}()

		// @TODO: Error handling.
		_ = websocket.JSON.Send(ws, game.World)
		api.Websockets = append(api.Websockets, ws)

		for {
			time.Sleep(1 * time.Second)
		}
	}))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
	fmt.Println("ddddd")
	for {
		current_state := <- game_state_chan
		if !current_state.HasLivingPlayers() {
			break
		}
		for ws := range(api.Websockets) {
			_ = websocket.JSON.Send(api.Websockets[ws], current_state)
		}
	}
	wg.Done()
}*/
