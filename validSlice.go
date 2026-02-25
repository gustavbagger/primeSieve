package main

import "math/big"

var smallPrimes = []int{3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47}

type ModState struct {
	mods []int
}

func newModState() *ModState {
	m := &ModState{
		mods: make([]int, len(smallPrimes)),
	}
	for i := range m.mods {
		m.mods[i] = 1
	}
	return m
}

func (m *ModState) pushPrimeExp(p, e int) {
	for i, sp := range smallPrimes {
		pe := 1
		base := p % sp
		for k := 0; k < e; k++ {
			pe = (pe * base) % sp
		}
		m.mods[i] = (m.mods[i] * pe) % sp
	}
}

func (m *ModState) isInvalid() bool {
	for i, q := range smallPrimes {
		if m.mods[i] == q-1 {
			return true
		}
	}
	return false
}

func validExponentSet(indexes, exponents, allValues []int) bool {
	prod := big.NewInt(1)
	for i, index := range indexes {
		p := big.NewInt(int64(allValues[index]))
		prod.Mul(prod, new(big.Int).Exp(p, big.NewInt(int64(exponents[i])), nil))
	}
	prod.Add(prod, big.NewInt(1))
	return prod.ProbablyPrime(32)
}
