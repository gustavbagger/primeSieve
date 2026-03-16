package main

import (
	"math/rand"
	"testing"

	pr "github.com/fxtlabs/primes"
)

func BenchmarkValidExponentSets(b *testing.B) {
	const size = 34

	// Generate deterministic pseudo-random test data
	randSrc := rand.New(rand.NewSource(123456))

	exponents := make([]int, size)
	allValues := pr.Sieve(10000)

	indexes := make([]int, size)

	max := 3
	cases := make([][]int, max)
	for i := 1; i < size; i++ {
		indexes[i] = i
		exponents[i] = 1 + randSrc.Intn(5)
	}
	cases[0] = indexes
	cases[1] = []int{0, 1, 2, 3, 4, 5, 6, 7, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 27, 28, 29, 30, 31, 32, 33, 34}
	cases[2] = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 28, 29, 30, 31, 35, 38, 48}

	b.Run("uint192 version", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < max; j++ {
				validExponentSet192(cases[j], exponents, allValues)
			}
		}
	})

	b.Run("big.Int version", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < max; j++ {
				validExponentSet(cases[j], exponents, allValues)
			}
		}
	})
}
