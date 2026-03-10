package main

import "math/rand"

func (x *Int) ProbablyPrime(n int) bool {
	// Note regarding the doc comment above:
	// It would be more precise to say that the Baillie-PSW test uses the
	// extra strong Lucas test as its Lucas test, but since no one knows
	// how to tell any of the Lucas tests apart inside a Baillie-PSW test
	// (they all work equally well empirically), that detail need not be
	// documented or implicitly guaranteed.
	// The comment does avoid saying "the" Baillie-PSW test
	// because of this general ambiguity.

	if n < 0 {
		panic("negative n for ProbablyPrime")
	}
	if x.neg || len(x.abs) == 0 {
		return false
	}

	// primeBitMask records the primes < 64.
	const primeBitMask uint64 = 1<<2 | 1<<3 | 1<<5 | 1<<7 |
		1<<11 | 1<<13 | 1<<17 | 1<<19 | 1<<23 | 1<<29 | 1<<31 |
		1<<37 | 1<<41 | 1<<43 | 1<<47 | 1<<53 | 1<<59 | 1<<61

	w := x.abs[0]
	if len(x.abs) == 1 && w < 64 {
		return primeBitMask&(1<<w) != 0
	}

	if w&1 == 0 {
		return false // x is even
	}

	const primesA = 3 * 5 * 7 * 11 * 13 * 17 * 19 * 23 * 37
	const primesB = 29 * 31 * 41 * 43 * 47 * 53

	var rA, rB uint32
	switch _W {
	case 32:
		rA = uint32(x.abs.modW(primesA))
		rB = uint32(x.abs.modW(primesB))
	case 64:
		r := x.abs.modW((primesA * primesB) & _M)
		rA = uint32(r % primesA)
		rB = uint32(r % primesB)
	default:
		panic("math/big: invalid word size")
	}

	if rA%3 == 0 || rA%5 == 0 || rA%7 == 0 || rA%11 == 0 || rA%13 == 0 || rA%17 == 0 || rA%19 == 0 || rA%23 == 0 || rA%37 == 0 ||
		rB%29 == 0 || rB%31 == 0 || rB%41 == 0 || rB%43 == 0 || rB%47 == 0 || rB%53 == 0 {
		return false
	}

	stk := getStack()
	defer stk.free()
	return x.abs.probablyPrimeMillerRabin(stk, n+1, true) && x.abs.probablyPrimeLucas(stk)
}

// probablyPrimeMillerRabin reports whether n passes reps rounds of the
// Miller-Rabin primality test, using pseudo-randomly chosen bases.
// If force2 is true, one of the rounds is forced to use base 2.
// See Handbook of Applied Cryptography, p. 139, Algorithm 4.24.
// The number n is known to be non-zero.
func (n nat) probablyPrimeMillerRabin(stk *stack, reps int, force2 bool) bool {
	nm1 := nat(nil).sub(n, natOne)
	// determine q, k such that nm1 = q << k
	k := nm1.trailingZeroBits()
	q := nat(nil).rsh(nm1, k)

	nm3 := nat(nil).sub(nm1, natTwo)
	rand := rand.New(rand.NewSource(int64(n[0])))

	var x, y, quotient nat
	nm3Len := nm3.bitLen()

NextRandom:
	for i := 0; i < reps; i++ {
		if i == reps-1 && force2 {
			x = x.set(natTwo)
		} else {
			x = x.random(rand, nm3, nm3Len)
			x = x.add(x, natTwo)
		}
		y = y.expNN(stk, x, q, n, false)
		if y.cmp(natOne) == 0 || y.cmp(nm1) == 0 {
			continue
		}
		for j := uint(1); j < k; j++ {
			y = y.sqr(stk, y)
			quotient, y = quotient.div(stk, y, y, n)
			if y.cmp(nm1) == 0 {
				continue NextRandom
			}
			if y.cmp(natOne) == 0 {
				return false
			}
		}
		return false
	}

	return true
}
