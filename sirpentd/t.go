package main

import (
  "fmt"
  "github.com/sirpent-team/sirpent-go"
  // "github.com/davecgh/go-spew/spew"
)

func main() {
  players := make(map[sirpent.UUID]*sirpent.Player)
  grid := sirpent.HexGridHexagonal{Rings: 20}
  game := sirpent.NewGame(grid, players)

  fmt.Printf("grid = %#v\n", grid)
  fmt.Printf("grid = %#v\n", game.Grid)
  fmt.Printf("CellNeighbour = %#v\n", game.Grid.CellNeighbour)

  game.Grid.CellNeighbour(sirpent.Vector{1, 2, 3}, "SOUTHWEST")

  v, err := game.Grid.CryptoRandomCell()
  if err != nil {
    panic(err)
  }

  snake := sirpent.NewSnake(v)
  fmt.Printf("snake = %#v\n", snake)

  fmt.Printf("snake[0] = %#v\n", snake[0])
  game.Grid.CellNeighbour(snake[0], "SOUTHWEST")

  snake = snake.Grow(v)
  //snake = append(snake, v)
  fmt.Printf("snake = %#v\n", snake)

  v = game.Grid.CellNeighbour(snake[1], "SOUTHWEST")
  snake = snake.Grow(v)
  //snake = append(snake, v)
  fmt.Printf("snake = %#v\n", snake)
}
