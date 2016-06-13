package sirpent

type SnakeSegment struct {
	Position Vector
}

type Snake struct {
	Length   int
	Dead     bool
	Segments []SnakeSegment
}

func NewSnake(start_position Vector) *Snake {
	s := Snake{Length: 1, Dead: false, Segments: make([]SnakeSegment, 1)}
	s.Segments[0].Position = start_position
	return &s
}

func (s Snake) IsHeadAt(v Vector) bool {
	return (s.Length > 0 && s.Segments[0].Position == v)
}

func (s Snake) HasSegmentAt(v Vector) bool {
	for i := 0; i < s.Length; i++ {
		if s.Segments[i].Position == v {
			return true
		}
	}
	return false
}

func (s Snake) HasCollidedIntoSnake(s2 Snake) bool {
	return (s.Length > 0 && s2.HasSegmentAt(s.Segments[0].Position))
}

func (s *Snake) StepInDirection(direction int) {
	// Move each segments to their parent location. Skip head for obvious reasons.
	for i := s.Length - 1; i > 0; i-- {
		s.Segments[i].Position = s.Segments[i-1].Position
	}

	// Update head position.
	s.Segments[0].Position = s.Segments[0].Position.Neighbour(direction)
}
