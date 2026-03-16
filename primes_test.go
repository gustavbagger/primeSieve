package main

import (
	"math/rand"
	"testing"

	pr "github.com/fxtlabs/primes"
)

func BenchmarkValidExponentSets(b *testing.B) {
	// Choose a realistic size for your exponent sets
	const size = 33

	// Generate deterministic pseudo-random test data
	randSrc := rand.New(rand.NewSource(123456))

	exponents := make([]int, size)
	allValues := pr.Sieve(1000)

	indexes := make([]int, size)
	for i := 0; i < size; i++ {
		indexes[i] = i
		exponents[i] = randSrc.Intn(5) + 1
	}

	b.Run("uint192 version", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			validExponentSet192(indexes, exponents, allValues)

		}
	})

	b.Run("big.Int version", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			validExponentSet(indexes, exponents, allValues)
		}
	})
}
