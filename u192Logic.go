package main

import (
	"math/bits"
)

// unsigned 192-bit int
type uint192 struct {
	Lo  uint64
	Mid uint64
	Hi  uint64
}

func cmp192(a, b uint192) int {
	if a.Hi != b.Hi {
		if a.Hi < b.Hi {
			return -1
		} else {
			return 1
		}

	} else if a.Mid != b.Mid {
		if a.Mid < b.Mid {
			return -1
		} else {
			return 1
		}
	} else if a.Lo != b.Lo {
		if a.Lo < b.Lo {
			return -1
		} else {
			return 1
		}
	}
	return 0
}

func sub192(a, b uint192) uint192 {
	lo, carry := bits.Sub64(a.Lo, b.Lo, 0)
	mid, carry := bits.Sub64(a.Mid, b.Mid, carry)
	hi, _ := bits.Sub64(a.Hi, b.Hi, carry)
	return uint192{Lo: lo, Mid: mid, Hi: hi}
}

func add192(a, b uint192) uint192 {
	var out uint192
	var carry uint64
	out.Lo, carry = bits.Add64(a.Lo, b.Lo, 0)
	out.Mid, carry = bits.Add64(a.Mid, b.Mid, carry)
	out.Hi, _ = bits.Add64(a.Hi, b.Hi, carry)
	return out
}

func mul192(a, b uint192) [6]uint64 {
	var t [6]uint64

	// helper to add a 128-bit value (hi:lo) into t[k..]
	add128 := func(k int, hi, lo uint64) {
		var c uint64
		t[k], c = bits.Add64(t[k], lo, 0)
		t[k+1], c = bits.Add64(t[k+1], hi, c)
		i := k + 2
		for c != 0 && i < 6 {
			t[i], c = bits.Add64(t[i], 0, c)
			i++
		}
	}

	// a.Lo * b.Lo
	hi, lo := bits.Mul64(a.Lo, b.Lo)
	add128(0, hi, lo)

	// a.Lo * b.Mid
	hi, lo = bits.Mul64(a.Lo, b.Mid)
	add128(1, hi, lo)

	// a.Lo * b.Hi
	hi, lo = bits.Mul64(a.Lo, b.Hi)
	add128(2, hi, lo)

	// a.Mid * b.Lo
	hi, lo = bits.Mul64(a.Mid, b.Lo)
	add128(1, hi, lo)

	// a.Mid * b.Mid
	hi, lo = bits.Mul64(a.Mid, b.Mid)
	add128(2, hi, lo)

	// a.Mid * b.Hi
	hi, lo = bits.Mul64(a.Mid, b.Hi)
	add128(3, hi, lo)

	// a.Hi * b.Lo
	hi, lo = bits.Mul64(a.Hi, b.Lo)
	add128(2, hi, lo)

	// a.Hi * b.Mid
	hi, lo = bits.Mul64(a.Hi, b.Mid)
	add128(3, hi, lo)

	// a.Hi * b.Hi
	hi, lo = bits.Mul64(a.Hi, b.Hi)
	add128(4, hi, lo)

	return t
}

func mulMod192(a, b uint192) uint192 {
	t := mul192(a, b)
	return uint192{Hi: t[2], Mid: t[1], Lo: t[0]}
}

// assuming 192>=n>=0
func lSH192(x uint192, n int) uint192 {
	var out uint192
	switch {
	case n < 64:
		out.Hi = (x.Hi << n) | (x.Mid >> (64 - n))
		out.Mid = (x.Mid << n) | (x.Lo >> (64 - n))
		out.Lo = x.Lo << n
	case n < 128:
		s := n - 64
		out.Hi = (x.Mid << s) | (x.Lo >> (64 - s))
		out.Mid = x.Lo << s
	default:
		s := n - 128
		out.Hi = x.Lo << s
	}
	return out
}

// assuming 192>=n>=0
func rSH192(x uint192, n int) uint192 {
	var out uint192
	switch {
	case n < 64:
		out.Hi = x.Hi >> n
		out.Mid = (x.Hi << (64 - n)) | (x.Mid >> n)
		out.Lo = (x.Mid << (64 - n)) | (x.Lo >> n)
	case n < 128:
		s := n - 64
		out.Mid = x.Hi >> s
		out.Lo = (x.Hi << (64 - s)) | (x.Mid >> s)
	default:
		s := n - 128
		out.Lo = x.Hi >> s
	}
	return out
}

// Assuming x non-zero
func twoAdicVal192(x uint192) int {
	if x.Lo != 0 {
		return bits.TrailingZeros64(x.Lo)
	} else if x.Mid != 0 {
		return bits.TrailingZeros64(x.Mid) + 64
	} else {
		return bits.TrailingZeros64(x.Hi) + 128
	}
}

// Assuming x non-zero
func LeadingZeros192(x uint192) int {
	if x.Hi != 0 {
		return bits.LeadingZeros64(x.Hi)
	} else if x.Mid != 0 {
		return bits.LeadingZeros64(x.Mid) + 64
	} else {
		return bits.LeadingZeros64(x.Lo) + 128
	}
}

func isZero192(x uint192) bool {
	return x.Lo == 0 && x.Mid == 0 && x.Hi == 0
}

// this is just add, then reduce modulo n once
func montAddReduce(a, b, n uint192) uint192 {
	s := add192(a, b)
	if cmp192(s, n) >= 0 {
		s = sub192(s, n)
	}
	return s
}

func bit192(x uint192, i int) uint64 {
	switch {
	case i < 64:
		return (x.Lo >> i) & 1
	case i < 128:
		return (x.Mid >> (i - 64)) & 1
	default:
		return (x.Hi >> (i - 128)) & 1
	}
}

//<- Trust --- Dont Trust ->

// REDC reduces C mod N using Montgomery reduction with β = 2^64, n = 3.
// C is a 384-bit value stored as [6]uint64 (little-endian: C[0] least significant).
// N is a 192-bit odd modulus. mu = -N^{-1} mod 2^64.
func REDC(C *[6]uint64, N uint192, mu uint64) uint192 {
	for i := 0; i < 3; i++ {
		// q = (mu * C[i]) mod 2^64
		q := mu * C[i]

		// C += q * N * β^i
		var carry uint64

		// word i: q * N.Lo
		hi, lo := bits.Mul64(q, N.Lo)
		C[i], carry = bits.Add64(C[i], lo, 0)
		carry, _ = bits.Add64(hi, 0, carry)

		// word i+1: q * N.Mid
		hi, lo = bits.Mul64(q, N.Mid)
		C[i+1], carry = bits.Add64(C[i+1], lo, carry)
		carry, _ = bits.Add64(hi, 0, carry)

		// word i+2: q * N.Hi
		hi, lo = bits.Mul64(q, N.Hi)
		C[i+2], carry = bits.Add64(C[i+2], lo, carry)
		carry, _ = bits.Add64(hi, 0, carry)

		// propagate carry upward
		j := i + 3
		for carry != 0 && j < 6 {
			C[j], carry = bits.Add64(C[j], 0, carry)
			j++
		}
	}

	// R = C >> (64*3) = top 3 words
	R := uint192{
		Lo:  C[3],
		Mid: C[4],
		Hi:  C[5],
	}

	// Final conditional subtraction: if R >= N, subtract N
	if cmp192(R, N) >= 0 {
		R = sub192(R, N)
	}

	return R
}

func montMul192(a, b, N uint192, mu uint64) uint192 {
	C := mul192(a, b)
	return REDC(&C, N, mu)
}

func inv64(x uint64) uint64 {
	y := x
	// Newton iteration: y_{k+1} = y_k * (2 - x*y_k) mod 2^64
	y *= 2 - x*y
	y *= 2 - x*y
	y *= 2 - x*y
	y *= 2 - x*y
	y *= 2 - x*y
	y *= 2 - x*y
	return y
}

func strongPRP(N uint192) bool {
	d := rSH192(N, 1) //since N odd, this gives (N-1)/2
	s := twoAdicVal192(d)
	d = rSH192(d, s) // odd part of N-1
	s += 1           // gain another factor from first rshift, now N-1 = 2^s*d, with d odd

	/*
		Idea: Compute the Montgomery forms we need, choose R = 2^192 to make sure we always have R> N
		R is chosen as a power of 2 since it is easy to do division
		Precompute R mod N = mont(1) and mont(-1) = -R modulo N so we can check prime conditions easily without needing to convert back to regular nrs
		Precompute R^{-1} since we want to do (aR mod N)(bR mod N)R^{-1} = (ab)R mod N
		(this is what REDC is doing, we input R=2^192, N our modulus, T = aR * bR and get out abR modulo N = (ab))

	*/

	//compute 2^d modulo N, base case

	//check if base case x is 1 or -1 modulo N, if so its likely prime

	//for each x_i = x^(2^i) for 0<i<s check if x_i = -1 modulo. If so, its likely prime

	//if none of these work, then N is composite for sure

	return true
}
