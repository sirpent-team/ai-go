package sirpent_test

/*import (
	"github.com/sirpent-team/sirpent-go"
	"testing"
)*/

/*func TestDistance(t *testing.T) {
	v1 := sirpent.HexagonalVector{X: 11, Y: -33}
	v2 := sirpent.HexagonalVector{X: -11, Y: 33}
	v3 := sirpent.HexagonalVector{X: -3, Y: 99}

	d11 := v1.Distance(v1) // Distance to self.
	expected_d11 := sirpent.HexagonalVector{X: 0, Y: 0}
	if d11 != expected_d11 {
		t.Error(v1, "had a non-zero distance", d11, "to itself.")
	}
	d12 := v1.Distance(v2)
	expected_d12 := sirpent.HexagonalVector{X: -22, Y: 66}
	if d12 != expected_d12 {
		t.Error(v1, "to", v2, "had an incorrect distance expected/actual", expected_d12, d12)
	}
	d13 := v1.Distance(v3)
	expected_d13 := sirpent.HexagonalVector{X: -14, Y: 132}
	if d13 != expected_d13 {
		t.Error(v1, "to", v3, "had an incorrect distance expected/actual", expected_d13, d13)
	}

	d31 := v3.Distance(v1)
	expected_d31 := sirpent.HexagonalVector{X: -expected_d13.X, Y: -expected_d13.Y}
	if d31 != expected_d31 {
		t.Error(v1, "to", v3, "did not have consistent distances. d13/d31:", d13, d31)
	}
}

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
			expected_within_bounds := (v.X >= 0 && v.X < 10 && v.Y >= 0 && v.Y < 10)
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
*/
