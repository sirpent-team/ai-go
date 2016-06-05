package sirpent

type Vector interface {
	Distance(v2 Vector) Vector
	Neighbour(direction int) Vector
}

type Grid interface {
	Dimensions() []int
	IsWithinBounds(v Vector) bool
}
