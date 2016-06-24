package main

import (
	crypto_rand "crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/sirpent-team/sirpent-go"
	"math/big"
	"net"
	"os"
	//"time"
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

	var hex_grid sirpent.HexHexGrid
	err = pc.Decoder.Decode(&hex_grid)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Hex Grid = %s\n", hex_grid)

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

		var direction sirpent.Direction
		for {
			direction_index := crypto_int(0, len(sirpent.Directions()))
			direction = sirpent.Direction(direction_index)

			snake := gs.Plays[player_id].Snake
			head := snake[0]
			//growing := head == gs.Food
			directed_head := head.Neighbour(direction)
			// @TODO: Somehow exactly 1 snake dies quickly, by moving onto their first tail segment.
			// I don't get how. The move does definitely result in self-intersection but how this
			// gets past these checks I do not know. It almost makes more sense if the AI has incorrect
			// information on its own snake but the player id in these cases *is* correct.
			fmt.Printf("snake=%+v directed_head=%+v\n", snake, directed_head)
			if !snake.TailContains(directed_head) && hex_grid.IsWithinBounds(directed_head) {
				break
			}
		}
		fmt.Println(direction)
		err = pc.Encoder.Encode(direction) //sirpent.SouthEast)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	/*for {
		for i := 0; i < 3; i++ {
			bufio.NewReader(conn).ReadString('\n')
			conn.Write([]byte("SOUTHEAST\n"))
		}

		for i := 0; i < 3; i++ {
			bufio.NewReader(conn).ReadString('\n')
			//time.Sleep(20000 * time.Millisecond)
			conn.Write([]byte("N"))
			//time.Sleep(2000 * time.Millisecond)
			conn.Write([]byte("O"))
			//time.Sleep(2000 * time.Millisecond)
			conn.Write([]byte("RT"))
			//time.Sleep(2000 * time.Millisecond)
			conn.Write([]byte("H\n"))
			//time.Sleep(2000 * time.Millisecond)
		}

		for i := 0; i < 3; i++ {
			bufio.NewReader(conn).ReadString('\n')
			conn.Write([]byte("NORTHWEST\n"))
		}

		for i := 0; i < 3; i++ {
			bufio.NewReader(conn).ReadString('\n')
			conn.Write([]byte("SOUTH\n"))
		}
	}*/
}
