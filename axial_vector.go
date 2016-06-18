package sirpent

import (
	"fmt"
)

type AxialVector struct {
	Q, R int
}

func (av AxialVector) AsCubeVector() CubeVector {
	return CubeVector{X: av.Q, Z: av.R, Y: 0 - av.Q - av.R}
}

func (av AxialVector) Neighbour(d Direction) AxialVector {
	// This will compile down to the Axial neighbour code.
	return av.AsCubeVector().Neighbour(d).AsAxialVector()
}

func (av AxialVector) Neighbours() []AxialVector {
	directions := Directions()
	ns := make([]AxialVector, len(directions))
	for i := range directions {
		ns[i] = av.Neighbour(directions[i])
	}
	return ns
}

func (av AxialVector) Distance(av2 AxialVector) int {
	return av.AsCubeVector().Distance(av2.AsCubeVector())
}

func (av AxialVector) String() string {
	return fmt.Sprintf("AxialVector(Q=%d, R=%d)", av.Q, av.R)
}
