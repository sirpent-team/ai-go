package sirpent_test

import (
  "testing"
  "github.com/Taneb/sirpent"
)

func TestNeighbours(t *testing.T) {
  v := sirpent.HexagonalVector{X: 0, Y: 0}
  var directions [6]int
  var expected_direction_vectors [6]sirpent.HexagonalVector

  directions[0] = sirpent.NORTH
  expected_direction_vectors[0] = sirpent.HexagonalVector{X: v.X, Y: v.Y - 1}
  directions[1] = sirpent.NORTHEAST
  expected_direction_vectors[1] = sirpent.HexagonalVector{X: v.X + 1, Y: v.Y - 1}
  directions[2] = sirpent.SOUTHEAST
  expected_direction_vectors[2] = sirpent.HexagonalVector{X: v.X + 1, Y: v.Y}
  directions[3] = sirpent.SOUTH
  expected_direction_vectors[3] = sirpent.HexagonalVector{X: v.X, Y: v.Y + 1}
  directions[4] = sirpent.SOUTHWEST
  expected_direction_vectors[4] = sirpent.HexagonalVector{X: v.X - 1, Y: v.Y + 1}
  directions[5] = sirpent.NORTHWEST
  expected_direction_vectors[5] = sirpent.HexagonalVector{X: v.X - 1, Y: v.Y}

  for i := 0; i < 6; i++ {
    actual_direction_vector := v.Neighbour(directions[i])
    if actual_direction_vector != expected_direction_vectors[i] {
      t.Error("In direction", directions[i], "direction vector was incorrect. Expected/actual:", expected_direction_vectors[i], actual_direction_vector)
    }
  }
}

func TestIsWithinBounds(t *testing.T) {
  g := sirpent.HexagonalGrid{Width: 10, Height: 10}

  for i := -2; i < 12; i++ {
    for j := -2; j < 12; j++ {
      v := sirpent.HexagonalVector{X: i, Y: j}
      expected_within_bounds := (v.X >= 0 && v.X < 10 && v.Y >=0 && v.Y < 10)
      actual_within_bounds := g.IsWithinBounds(v)
      if expected_within_bounds && !actual_within_bounds {
        t.Error(v, "was wrongly out of bounds.")
      }
      if !expected_within_bounds && actual_within_bounds {
        t.Error(v, "was wrongly within bounds.")
      }
    }
  }
}
