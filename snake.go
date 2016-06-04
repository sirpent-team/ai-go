package sirpent

type SnakeSegment struct {
  Position Vector
}

type Snake struct {
  Length int
  Dead bool
  Segments []SnakeSegment
}

func (s Snake) IsHeadAt(v Vector) bool {
  panic("Snake.IsHeadAt not yet implemented.")
}

func (s Snake) IsASegmentAt(v Vector) bool {
  panic("Snake.IsASegmentAt not yet implemented.")
}

func (s Snake) IsCollidingWithSnake(s2 Snake) bool {
  panic("Snake.IsCollidingWithSnake not yet implemented.")
}
