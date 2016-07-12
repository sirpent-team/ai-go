package sirpent

type Snake []Vector

func NewSnake(start_position Vector) Snake {
	return Snake{start_position}
}

func (snake Snake) Copy() Snake {
	copied_snake := make(Snake, len(snake))
	// Prevent updating the slice of the previous snake.
	copy(copied_snake, snake)
	return copied_snake
}

func (snake Snake) Move(grid Grid, direction Direction) Snake {
	moved_snake := snake.Copy()

	// Discard final segment.
	moved_snake = moved_snake[:len(snake)-1]

	// Prepend the new, moved head.
	moved_head := grid.CellNeighbour(moved_snake[0], direction)
	moved_snake = append(Snake{moved_head}, moved_snake...)

	return moved_snake
}

func (snake Snake) Grow(v Vector) Snake {
	grown_snake := snake.Copy()
	grown_snake = append(grown_snake, v)
	return grown_snake
}
