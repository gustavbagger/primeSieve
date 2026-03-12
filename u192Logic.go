package main

import (
	"errors"
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

func sub192(a, b uint192) (uint192, error) {
	lo, carry := bits.Sub64(a.Lo, b.Lo, 0)
	mid, carry := bits.Sub64(a.Mid, b.Mid, carry)
	hi, carry := bits.Sub64(a.Hi, b.Hi, carry)
	if carry != 0 {
		return uint192{}, errors.New("a<b")
	}
	return uint192{Lo: lo, Mid: mid, Hi: hi}, nil
}

// no need to worry about carry here: a,b are always chosen to be <N and 2N < 2^192
func add192(a, b uint192) uint192 {
	var out uint192
	var carry uint64
	out.Lo, carry = bits.Add64(a.Lo, b.Lo, 0)
	out.Mid, carry = bits.Add64(a.Mid, b.Mid, carry)
	out.Hi, _ = bits.Add64(a.Hi, b.Hi, carry)
	return out
}

func mul192(a, b uint192) [7]uint64 {
	var t [7]uint64

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

func neg192(x uint192) uint192 {
	var z uint192
	carry := uint64(1)

	// z = (~x + 1)  (two's complement)
	z.Lo = ^x.Lo + carry
	carry = 0
	if z.Lo == 0 {
		carry = 1
	}

	z.Mid = ^x.Mid + carry
	carry = 0
	if z.Mid == 0 && carry == 1 {
		carry = 1
	} else {
		carry = 0
	}

	z.Hi = ^x.Hi + carry

	return z
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

func montMul192(a, b, N uint192, mu uint64) (uint192, error) {
	return REDC(mul192(a, b), N, mu)
}

// this is just add, then reduce modulo n once. Assumes a,b<n
func montAddReduce(a, b, n uint192) (uint192, error) {
	s := add192(a, b)
	var err error
	if cmp192(s, n) >= 0 {
		s, err = sub192(s, n)
		if err != nil {
			return uint192{}, err
		}
	}
	return s, nil
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

var two = uint192{Lo: 2}

// newtons method (hensel lifting)
func inv192(N uint192) (uint192, error) {
	x := uint192{Lo: 1} //seed

	for i := 0; i < 8; i++ {
		s, err := sub192(two, mulMod192(N, x))
		if err != nil {
			return uint192{}, err
		}
		x = mulMod192(x, s)
	}
	return x, nil
}

// inv64 computes N^{-1} mod 2^64 using Newton iteration.
// Requires N to be odd.
func inv64(N uint64) uint64 {
	// Seed: inverse mod 2 (i.e., 1)
	x := uint64(1)

	// Each iteration doubles the number of correct bits.
	for i := 0; i < 6; i++ {
		// x = x * (2 - N*x) mod 2^64
		x = x * (2 - N*x)
	}

	return x
}

func addMod192(a, b, N uint192) (uint192, error) {
	var out uint192
	var carry uint64
	var err error

	// Add Lo limbs
	out.Lo, carry = bits.Add64(a.Lo, b.Lo, 0)

	// Add Mid limbs + carry
	out.Mid, carry = bits.Add64(a.Mid, b.Mid, carry)

	// Add Hi limbs + carry
	out.Hi, _ = bits.Add64(a.Hi, b.Hi, carry)

	// If result >= N, subtract N
	if cmp192(out, N) >= 0 {
		out, err = sub192(out, N)
		if err != nil {
			return uint192{}, err
		}
	}

	return out, nil
}

func montOne(N uint192) (uint192, error) {
	x := uint192{Lo: 1}
	var err error

	for i := 0; i < 192; i++ {
		x, err = addMod192(x, x, N)
		if err != nil {
			return uint192{}, err
		}
	}
	return x, nil
}

// Consider using the known value of s since N is computed via the exponent set
func strongPRP(N uint192) (bool, error) {
	d := rSH192(N, 1)                 //since N odd, this gives (N-1)/2
	s := twoAdicVal192(d)             // this is the 2-adic valuation of (N-1)/2
	d = lSH192(d, LeadingZeros192(d)) // odd part of N-1 shifted up to the MSB

	nbar := -inv64(N.Lo)       //Pre-compute -N^{-1} modulo 2^192 for REDC
	montOne, err := montOne(N) //Pre-compute 1R modulo N for primality checks
	if err != nil {
		return false, err
	}
	montNegOne, err := sub192(N, montOne) //Pre-compute -1R modulo N for primality checks
	if err != nil {
		return false, err
	}

	x, err := montAddReduce(montOne, montOne, N) //Initialise x = 2R mod N
	if err != nil {
		return false, err
	}
	//compute 2^d modulo N by MSB->LSB exponentiation
	for d = lSH192(d, 1); !isZero192(d); d = lSH192(d, 1) {
		x, err = montMul192(x, x, N, nbar)
		if err != nil {
			return false, err
		}
		if (d.Hi & (1 << 63)) != 0 { //This is checking if the MSB is 1, if so d.Hi is interpreted as negative as a uint192
			x, err = montAddReduce(x, x, N)
			if err != nil {
				return false, err
			}
		}
	}

	//check if 2^d is +-1 modulo N, if so, N is prop prime
	if cmp192(x, montOne) == 0 || cmp192(x, montNegOne) == 0 {
		return true, nil
	}

	// for each x_i = x^(d * 2^i) for 0<i<s check if x_i = -1 modulo. If so, its likely prime
	for s--; s >= 0; s-- {
		x, err = montMul192(x, x, N, nbar)
		if err != nil {
			return false, err
		}
		if cmp192(x, montNegOne) == 0 {
			return true, nil
		}
	}
	return false, nil

	/*
			Idea: Compute the Montgomery forms we need, choose R = 2^192 to make sure we always have R> N
			R is chosen as a power of 2 since it is easy to do division

		Precompute R mod N = mont(1) and mont(-1) = -R modulo N
			so we can check prime conditions easily without needing to convert back to regular nrs

		compute 2^d modulo N, base case
			do this by repeated squaring in montgomery form


		- check if base case x is 1 or -1 modulo N, if so its likely prime

		- for each x_i = x^(2^i) for 0<i<s check if x_i = -1 modulo. If so, its likely prime

		- if none of these work, then N is composite for sure
	*/
}

//<- Trust --- Dont Trust ->

// REDC reduces C mod N using Montgomery reduction with β = 2^64, n = 3.
// C is a 384-bit value stored as [6]uint64 (little-endian: C[0] least significant).
// N is a 192-bit odd modulus. mu = -N^{-1} mod 2^192.
func REDC(C [7]uint64, Nuint uint192, mu uint64) (uint192, error) {
	var err error
	N := [3]uint64{Nuint.Lo, Nuint.Mid, Nuint.Hi}
	// Step 1: main loop
	for i := 0; i < 3; i++ {
		ci := C[i]
		qi := ci * mu // automatically mod 2^64

		// Add qi * N shifted by i limbs
		var carry uint64
		var hi uint64

		// j runs over limbs of N
		for j := 0; j < 3; j++ {
			// Multiply qi * N[j]
			lo, hiMul := bits.Mul64(qi, N[j])

			// Add into C[i+j] with carry
			lo, c1 := bits.Add64(lo, C[i+j], 0)
			lo, c2 := bits.Add64(lo, carry, 0)

			C[i+j] = lo
			carry = hiMul + c1 + c2
		}

		// Propagate carry into higher limbs
		k := i + 3
		for carry != 0 {
			C[k], hi = bits.Add64(C[k], carry, 0)
			carry = hi
			k++
		}
	}

	// Step 2: R = C >> 192 bits (drop first 3 limbs)
	R := uint192{Lo: C[3], Mid: C[4], Hi: C[5]}

	// Step 3: final conditional subtraction
	if cmp192(R, Nuint) >= 0 {
		R, err = sub192(R, Nuint)
		if err != nil {
			return uint192{}, err
		}
	}

	return R, nil
}
