package sirpent_test

import (
	"github.com/Taneb/sirpent"
	"testing"
)

func TestIsHeadAt(t *testing.T) {
	s := sirpent.Snake{Length: 3, Dead: false, Segments: make([]sirpent.SnakeSegment, 3)}
	s.Segments[0] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 0, Y: 0}}
	s.Segments[1] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 1, Y: 0}}
	s.Segments[2] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 1, Y: 1}}

	v := sirpent.HexagonalVector{X: 0, Y: 0} // Head position.
	if !s.IsHeadAt(v) {
		t.Error("Head segment was incorrectly not occupying", v)
	}
	v = sirpent.HexagonalVector{X: 1, Y: 0} // Foremost non-head segment position.
	if s.IsHeadAt(v) {
		t.Error("Head segment was incorrectly occupying", v)
	}
	v = sirpent.HexagonalVector{X: 5, Y: -5} // None of the snake is here.
	if s.IsHeadAt(v) {
		t.Error("Head segment was incorrectly occupying", v)
	}
}

func TestHasSegmentAt(t *testing.T) {
	s := sirpent.Snake{Length: 3, Dead: false, Segments: make([]sirpent.SnakeSegment, 3)}
	s.Segments[0] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 0, Y: 0}}
	s.Segments[1] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 1, Y: 0}}
	s.Segments[2] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 1, Y: 1}}

	v := sirpent.HexagonalVector{X: 0, Y: 0} // Head position.
	if !s.HasSegmentAt(v) {
		t.Error("Head segment was incorrectly not occupying", v)
	}
	v = sirpent.HexagonalVector{X: 1, Y: 0} // Segment 1 position.
	if !s.HasSegmentAt(v) {
		t.Error("Second segment was incorrectly not occupying", v)
	}
	v = sirpent.HexagonalVector{X: 1, Y: 1} // Tail position.
	if !s.HasSegmentAt(v) {
		t.Error("Tail segment was incorrectly not occupying", v)
	}
	v = sirpent.HexagonalVector{X: 3, Y: 3} // None of the snake is here.
	if s.HasSegmentAt(v) {
		t.Error("A segment was incorrectly occupying", v)
	}
}

func TestHasCollidedIntoSnake(t *testing.T) {
	s1 := sirpent.Snake{Length: 3, Dead: false, Segments: make([]sirpent.SnakeSegment, 3)}
	s1.Segments[0] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 0, Y: 0}}
	s1.Segments[1] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 1, Y: 0}}
	s1.Segments[2] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 1, Y: 1}}

	s2 := sirpent.Snake{Length: 3, Dead: false, Segments: make([]sirpent.SnakeSegment, 3)} // No collision with s1
	s2.Segments[0] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 5, Y: 5}}
	s2.Segments[1] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 6, Y: 5}}
	s2.Segments[2] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 6, Y: 6}}

	s3 := sirpent.Snake{Length: 3, Dead: false, Segments: make([]sirpent.SnakeSegment, 3)} // s1 collides with head
	s3.Segments[0] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 0, Y: 0}}
	s3.Segments[1] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: -1, Y: 0}}
	s3.Segments[2] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: -2, Y: 0}}

	s4 := sirpent.Snake{Length: 3, Dead: false, Segments: make([]sirpent.SnakeSegment, 3)} // s1 collides with body
	s4.Segments[0] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 0, Y: 1}}
	s4.Segments[1] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 0, Y: 0}}
	s4.Segments[2] = sirpent.SnakeSegment{Position: sirpent.HexagonalVector{X: 0, Y: -1}}

	if s1.HasCollidedIntoSnake(s2) {
		t.Error("Incorrect collision with snake s2.")
	}
	if !s1.HasCollidedIntoSnake(s3) {
		t.Error("Unnoticed head collision with snake s3.")
	}
	if !s1.HasCollidedIntoSnake(s4) {
		t.Error("Unnoticed body collision with snake s4.")
	}
}
