package sirpent

/*
{
  "segments" : [
     {"y" : 22, "x" : -17 }
  ]
}
*/

type Snake struct {
	Segments []Vector `json:"segments"`
}

func NewSnake(start_position Vector) Snake {
	return Snake{Segments: []Vector{start_position}}
}

func (snake Snake) Copy() Snake {
	copied_snake := Snake{Segments: make([]Vector, len(snake.Segments))}
	// Prevent updating the slice of the previous snake.
	copy(copied_snake.Segments, snake.Segments)
	return copied_snake
}

func (snake Snake) Move(grid HexagonalGrid, direction Direction) Snake {
	moved_snake := snake.Copy()

	// Discard final segment.
	moved_snake.Segments = moved_snake.Segments[:len(snake.Segments)-1]

	// Prepend the new, moved head.
	moved_head := grid.CellNeighbour(snake.Segments[0], direction)
	moved_snake.Segments = append([]Vector{moved_head}, moved_snake.Segments...)

	return moved_snake
}

func (snake Snake) Grow(v Vector) Snake {
	grown_snake := snake.Copy()
	grown_snake.Segments = append(grown_snake.Segments, v)
	return grown_snake
}
