package sirpent

import (
	crypto_rand "crypto/rand"
	"github.com/satori/go.uuid"
	"math/big"
)

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func crypto_perm(n int) []int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		j_big, _ := crypto_rand.Int(crypto_rand.Reader, big.NewInt(int64(n)))
		j := int(j_big.Int64())
		m[i] = m[j]
		m[j] = i
	}
	return m
}

func crypto_int(lower int, upper int) int {
	n_big, err := crypto_rand.Int(crypto_rand.Reader, big.NewInt(int64(upper-lower)))
	if err != nil {
		panic(err)
	}
	n := int(n_big.Int64())
	return n + lower
}

func NewUUID() uuid.UUID {
	return uuid.NewV4()
}
