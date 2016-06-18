package sirpent

type Snake struct {
	Segments []AxialVector
}

func NewSnake(start_position AxialVector) Snake {
	// @TODO: Ideal capacity of a snake?
	s := Snake{Segments: make([]AxialVector, 1)}
	s.Segments[0] = start_position
	return s
}

func (s Snake) IsHeadAt(av AxialVector) bool {
	return len(s.Segments) > 0 && s.Segments[0] == av
}

func (s Snake) Contains(av AxialVector) bool {
	for i := range s.Segments {
		if s.Segments[i] == av {
			return true
		}
	}
	return false
}

func (s Snake) IsHeadInsideTail() bool {
	for i := range s.Segments {
		if i > 0 && s.Segments[0] == s.Segments[i] {
			return true
		}
	}
	return false
}

func (s Snake) IsHeadInsideSnake(s2 Snake) bool {
	return len(s.Segments) > 0 && s2.Contains(s.Segments[0])
}

func (s *Snake) StepInDirection(direction Direction, grow_extra_segment bool) {
	if grow_extra_segment {
		s.Segments = append(s.Segments, AxialVector{})
	}

	// Move each segments to their parent location. Skip head for obvious reasons.
	for i := len(s.Segments) - 1; i > 0; i-- {
		s.Segments[i] = s.Segments[i-1]
	}

	// Update head position.
	s.Segments[0] = s.Segments[0].Neighbour(direction)
}
