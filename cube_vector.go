package sirpent

import (
	"fmt"
)

type CubeVector struct {
	X, Y, Z int
}

func (cv CubeVector) AsAxialVector() AxialVector {
	return AxialVector{Q: cv.X, R: cv.Z}
}

func (cv CubeVector) Neighbour(d Direction) CubeVector {
	cv2 := cv
	switch d {
	case NorthEast:
		cv2.X++
		cv2.Z--
	case East:
		cv2.X++
		cv2.Y--
	case SouthEast:
		cv2.Y--
		cv2.Z++
	case SouthWest:
		cv2.X--
		cv2.Z++
	case West:
		cv2.X--
		cv2.Y++
	case NorthWest:
		cv2.Y++
		cv2.Z--
	}
	return cv2
}

func (cv CubeVector) Neighbours() []CubeVector {
	directions := Directions()
	ns := make([]CubeVector, len(directions))
	for i := range directions {
		ns[i] = cv.Neighbour(directions[i])
	}
	return ns
}

func (cv CubeVector) Distance(cv2 CubeVector) int {
	return max(max(abs(cv2.X-cv.X), abs(cv2.Y-cv.Y)), abs(cv2.Z-cv.Z))
}

func (cv CubeVector) String() string {
	return fmt.Sprintf("CubeVector(X=%d, Y=%d, Z=%d)", cv.X, cv.Y, cv.Z)
}
