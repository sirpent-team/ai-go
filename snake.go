package sirpent

type SnakeSegment struct {
	Position Vector
}

type Snake struct {
	Length   int
	Dead     bool
	Segments []SnakeSegment
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
