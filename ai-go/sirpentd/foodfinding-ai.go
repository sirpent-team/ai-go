package main

import (
	crypto_rand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirpent-team/sirpent-ai-go"
	"math/big"
	"net"
	"os"
	"log"
)

func main() {
	port := os.Args[1]
	conn, err := net.Dial("tcp", ":" + port)
	if err != nil {
		panic(fmt.Sprintf("Could not connect to a TCP server on port %s.", port))
	}
	handleConnection(conn)
}

type PlayerClient struct {
	//ID sirpent.UUID
	Socket  net.Conn
	Encoder *json.Encoder
	Decoder *json.Decoder
}

func NewPlayerClient(conn net.Conn) *PlayerClient {
	pc := PlayerClient{
		Socket: conn,
	}
	pc.Encoder = json.NewEncoder(pc.Socket)
	pc.Decoder = json.NewDecoder(pc.Socket)
	return &pc
}

func crypto_int(lower int, upper int) int {
	n_big, err := crypto_rand.Int(crypto_rand.Reader, big.NewInt(int64(upper-lower)))
	if err != nil {
		panic(err)
	}
	n := int(n_big.Int64())
	return n + lower
}

func decodeMsg(json_decoder *json.Decoder) (json.RawMessage, sirpent.Msg) {
	var raw json.RawMessage
	env := sirpent.Msg{
		Data: &raw,
	}
	if err := json_decoder.Decode(&env); err != nil {
		log.Fatal(err)
	}
	return raw, env
}

func handleConnection(conn net.Conn) {
	pc := NewPlayerClient(conn)

	// var version_msg sirpent.VersionMsg
	// err := pc.Decoder.Decode(&msg)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// version_msg = msg.Data.(sirpent.VersionMsg)
	// fmt.Println(version_msg)

	raw, env := decodeMsg(pc.Decoder)
	msg := sirpent.MsgKindHandlers[env.Msg]()
	if err := json.Unmarshal(raw, msg); err != nil {
		log.Fatal(err)
	}
	version_msg := msg.(*sirpent.VersionMsg)
	fmt.Printf("%+v\n", version_msg)

	register_msg := sirpent.RegisterMsg{DesiredName: "foodfinding-ai", Kind: "player"}
	msg = sirpent.Msg{Msg: sirpent.Register, Data: register_msg}
	err := pc.Encoder.Encode(msg)
	if err != nil {
		log.Fatal(err)
		return
	}

	raw, env = decodeMsg(pc.Decoder)
	msg = sirpent.MsgKindHandlers[env.Msg]()
	if err := json.Unmarshal(raw, msg); err != nil {
		log.Fatal(err)
	}
	welcome_msg := msg.(*sirpent.WelcomeMsg)
	fmt.Printf("%+v\n", welcome_msg)

	name := welcome_msg.Name;

	raw, env = decodeMsg(pc.Decoder)
	msg = sirpent.MsgKindHandlers[env.Msg]()
	if err := json.Unmarshal(raw, msg); err != nil {
		log.Fatal(err)
	}
	new_game_msg := msg.(*sirpent.NewGameMsg)
	fmt.Printf("%+v\n", new_game_msg)
	game := new_game_msg.Game

	for {
		raw, env = decodeMsg(pc.Decoder)
		msg = sirpent.MsgKindHandlers[env.Msg]()
		if err := json.Unmarshal(raw, msg); err != nil {
			log.Fatal(err)
		}
		turn_msg := msg.(*sirpent.TurnMsg)
		fmt.Printf("%+v\n", turn_msg)
		turn := turn_msg.Turn

		snake := turn.Snakes[name]
		head := snake.Segments[0]

		var path []sirpent.Direction
		var direction sirpent.Direction
		path, err = pathfind(game.Grid, snake, head, turn.Food[0], turn.Food[0])

		if err == nil && len(path) > 0 {
			direction = path[len(path)-1]
		} else {
			fmt.Println(err)
			directions := game.Grid.Directions()
			for i := range directions {
				direction = directions[i]
				neighbour := game.Grid.CellNeighbour(head, direction)
				grow_extra_segment := turn.Food[0].Eq(neighbour)
				neighbour_snake := snake.Move(game.Grid, direction)
				if grow_extra_segment {
					neighbour_snake = neighbour_snake.Grow(snake.Segments[len(snake.Segments)-1])
				}
				if game.Grid.IsCellWithinBounds(neighbour) { //&& !neighbour_snake.HeadIntersectsSelf() {
					break
				}
			}
		}

		move_msg := sirpent.MoveMsg{Direction: direction}
		fmt.Printf("%+v\n", move_msg)
		msg = sirpent.Msg{Msg: sirpent.Move, Data: move_msg}
		err := pc.Encoder.Encode(msg)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func pathfind(grid sirpent.HexagonalGrid, snake sirpent.Snake, start sirpent.Vector, end sirpent.Vector, food sirpent.Vector) ([]sirpent.Direction, error) {
	var frontier []sirpent.Vector
	frontier = append(frontier, start)
	came_from := make(map[sirpent.Vector]sirpent.Vector)
	came_from[start] = start
	cost_to := make(map[sirpent.Vector]int)
	cost_to[start] = 0
	direction_to := make(map[sirpent.Vector]sirpent.Direction)
	snake_at := make(map[sirpent.Vector]sirpent.Snake)

	var current sirpent.Vector
	for len(frontier) > 0 {
		lowest_expected_cost := 0
		lowest_expected_cost_index := 0
		for i := 0; i < len(frontier); i++ {
			current := frontier[i]
			expected_cost := cost_to[current] + grid.DistanceBetweenCells(current, end)
			if expected_cost < lowest_expected_cost {
				lowest_expected_cost = expected_cost
				lowest_expected_cost_index = i
			}
		}
		current = frontier[lowest_expected_cost_index]
		frontier = append(frontier[:lowest_expected_cost_index], frontier[lowest_expected_cost_index+1:]...)

		directions := grid.Directions()
		for i := range directions {
			direction := directions[i]
			neighbour := grid.CellNeighbour(current, direction)
			grow_extra_segment := food == neighbour
			_, already_reached := came_from[neighbour]
			if current == start {
				direction_to[start] = direction
				snake_at[start] = snake
			}
			neighbour_snake := snake_at[current].Move(grid, direction)
			if grow_extra_segment {
				neighbour_snake = neighbour_snake.Grow(snake_at[current].Segments[len(snake_at[current].Segments)-1])
			}
			neighbour_cost := cost_to[current] + 1
			if (!already_reached || cost_to[neighbour] > neighbour_cost) && grid.IsCellWithinBounds(neighbour) { //&& !neighbour_snake.HeadIntersectsSelf() {
				frontier = append(frontier, neighbour)
				came_from[neighbour] = current
				cost_to[neighbour] = neighbour_cost
				direction_to[neighbour] = direction
				snake_at[neighbour] = neighbour_snake
			}
		}
	}

	var path []sirpent.Direction
	_, ended := came_from[end]
	current = end
	if !ended {
		return path, errors.New(fmt.Sprintf("Could not path from %s to %s.", start, end))
	}
	for current != start {
		path = append(path, direction_to[current])
		current = came_from[current]
	}
	return path, nil
}
