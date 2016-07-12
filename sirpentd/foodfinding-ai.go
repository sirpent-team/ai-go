package main

import (
	crypto_rand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirpent-team/sirpent-go"
	"math/big"
	"net"
	"os"
)

func main() {
	port := os.Args[1] //"8901"

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(fmt.Sprintf("Could not establish TCP server on port %s.", port))
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(fmt.Sprintf("Could not listen on port %s.", port))
		}

		go handleConnection(conn)
	}
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

func handleConnection(conn net.Conn) {
	pc := NewPlayerClient(conn)

	var player_id sirpent.UUID
	err := pc.Decoder.Decode(&player_id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("player ID = %s\n", player_id)

	var game sirpent.Game
	err = pc.Decoder.Decode(&game)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		var gs sirpent.GameState
		err = pc.Decoder.Decode(&gs)
		if err != nil {
			fmt.Println(err)
			return
		}
		for player_id, player_state := range gs.Plays {
			fmt.Printf("( player_id=%s snake=%+v )\n", player_id, player_state.Snake)
		}

		snake := gs.Plays[player_id].Snake
		head := snake[0]

		var path []sirpent.Direction
		var direction sirpent.Direction
		path, err = pathfind(game.Grid, snake, head, gs.Food, gs.Food)

		if err == nil && len(path) > 0 {
			direction = path[len(path)-1]
		} else {
			directions := game.Grid.Directions()
			for i := range directions {
				direction = directions[i]
				neighbour := game.Grid.CellNeighbour(head, direction)
				grow_extra_segment := gs.Food.Eq(neighbour)
				neighbour_snake := snake.Move(game.Grid, direction)
				if grow_extra_segment {
					neighbour_snake = neighbour_snake.Grow(snake[len(snake)-1])
				}
				if game.Grid.IsCellWithinBounds(neighbour) { //&& !neighbour_snake.HeadIntersectsSelf() {
					break
				}
			}
		}

		/*var direction sirpent.Direction
		directions := sirpent.Directions
		for i := range directions {
			candidate_direction := directions[i]
			neighbour := head.Neighbour(candidate_direction)
			if (len(path) == 0 || neighbour == path[len(path)-1]) && !snake.Move(direction, false).HeadIntersectsSelf() {
				direction = candidate_direction
			}
		}*/

		fmt.Println(direction)
		err = pc.Encoder.Encode(sirpent.PlayerAction{Move: direction}) //sirpent.SouthEast)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func pathfind(grid sirpent.Grid, snake sirpent.Snake, start sirpent.Vector, end sirpent.Vector, food sirpent.Vector) ([]sirpent.Direction, error) {
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
				neighbour_snake = neighbour_snake.Grow(snake_at[current][len(snake_at[current])-1])
			}
			neighbour_cost := cost_to[current] + 1
			if (!already_reached || cost_to[neighbour] > neighbour_cost) && grid.IsCellWithinBounds(neighbour) /*&& !neighbour_snake.HeadIntersectsSelf()*/ {
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
