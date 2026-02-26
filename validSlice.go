package main

import "math/big"

var smallPrimes = []int{
	101, 103, 107, 109, 113, 127, 131, 137, 139, 149,
	151, 157, 163, 167, 173, 179, 181, 191, 193, 197,
	199, 211, 223, 227, 229, 233, 239, 241, 251, 257,
	263, 269, 271, 277, 281, 283, 293, 307, 311, 313,
	317, 331, 337, 347, 349, 353, 359, 367, 373, 379,
}

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

func validExponentSet(indexes, exponents, allValues []int) (*big.Int, bool) {
	prod := big.NewInt(1)
	for i, index := range indexes {
		pReg := allValues[index]
		eReg := exponents[i]
		p := big.NewInt(int64(pReg))
		prod.Mul(prod, new(big.Int).Exp(p, big.NewInt(int64(eReg)), nil))
	}
	n := big.NewInt(0)
	n.Add(prod, big.NewInt(1))
	return n, n.ProbablyPrime(32)
}
