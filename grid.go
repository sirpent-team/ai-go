package sirpent

type Vector interface {
	Coord(offset int) int
	Distance(v2 Vector) Vector
	Neighbour(direction int) Vector
}

type Grid interface {
	Dimensions() []int
	IsWithinBounds(v Vector) bool
}
