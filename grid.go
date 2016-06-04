package sirpent

type Vector interface {
  Neighbour(direction int) Vector
}

type Grid interface {
  Dimensions() []int
  IsWithinBounds(v Vector) bool
}
