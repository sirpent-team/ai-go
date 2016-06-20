package main

import (
	"fmt"
	"github.com/sirpent-team/sirpent-go"
)

func main() {
	grid := &sirpent.AxialGrid{Radius: 30, Origin: sirpent.AxialVector{Q: 0, R: 0}}
	i := 0
	for {
		if i%1000 == 0 {
			fmt.Printf("i=%d\n", i)
		}
		av := grid.CryptoRandomCell()
		if !grid.IsWithinBounds(av) {
			fmt.Printf("%s outside of bounds.\n", av.AsCubeVector())
		}
		i++
	}
}
