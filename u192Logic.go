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

func mul192(a, b uint192) [6]uint64 {
	var t [6]uint64

	// a.Lo * b.Lo
	lo, hi := bits.Mul64(a.Lo, b.Lo)
	t[0], t[1] = lo, hi

	// a.Lo * b.Mid
	lo, hi = bits.Mul64(a.Lo, b.Mid)
	t[1], hi = bits.Add64(t[1], lo, 0)
	t[2] += hi

	// a.Lo * b.Hi
	lo, hi = bits.Mul64(a.Lo, b.Hi)
	t[2], hi = bits.Add64(t[2], lo, 0)
	t[3] += hi

	// a.Mid * b.Lo
	lo, hi = bits.Mul64(a.Mid, b.Lo)
	t[1], hi = bits.Add64(t[1], lo, 0)
	t[2], hi = bits.Add64(t[2], hi, 0)
	t[3] += hi

	// a.Mid * b.Mid
	lo, hi = bits.Mul64(a.Mid, b.Mid)
	t[2], hi = bits.Add64(t[2], lo, 0)
	t[3], hi = bits.Add64(t[3], hi, 0)
	t[4] += hi

	// a.Mid * b.Hi
	lo, hi = bits.Mul64(a.Mid, b.Hi)
	t[3], hi = bits.Add64(t[3], lo, 0)
	t[4], hi = bits.Add64(t[4], hi, 0)
	t[5] += hi

	// a.Hi * b.Lo
	lo, hi = bits.Mul64(a.Hi, b.Lo)
	t[2], hi = bits.Add64(t[2], lo, 0)
	t[3], hi = bits.Add64(t[3], hi, 0)
	t[4] += hi

	// a.Hi * b.Mid
	lo, hi = bits.Mul64(a.Hi, b.Mid)
	t[3], hi = bits.Add64(t[3], lo, 0)
	t[4], hi = bits.Add64(t[4], hi, 0)
	t[5] += hi

	// a.Hi * b.Hi
	lo, hi = bits.Mul64(a.Hi, b.Hi)
	t[4], hi = bits.Add64(t[4], lo, 0)
	t[5] += hi

	return t
}

// MulRedc192 computes (x * y * R^-1) mod n for 192‑bit values,
// where R = 2^192 and npi = -n^{-1} mod 2^64 (using n.Lo). npi is mont_one
func MulRedc192(x, y, n uint192, npi uint64) uint192 {
	// Step 1: t = x * y (384‑bit)
	t := mul192(x, y) // t[0..5]

	// k = 3 limbs
	for i := 0; i < 3; i++ {
		// m_i = t_i * npi mod 2^64
		m, _ := bits.Mul64(t[i], npi)

		if m == 0 {
			continue
		}

		// t += m * n << (64*i)
		var carry uint64

		// add to limb i
		lo, hi := bits.Mul64(m, n.Lo)
		t[i], carry = bits.Add64(t[i], lo, 0)
		carry += hi

		// limb i+1
		lo, hi = bits.Mul64(m, n.Mid)
		t[i+1], carry = bits.Add64(t[i+1], lo, carry)
		carry += hi

		// limb i+2
		lo, hi = bits.Mul64(m, n.Hi)
		t[i+2], carry = bits.Add64(t[i+2], lo, carry)
		carry += hi

		// propagate carry into higher limbs
		j := i + 3
		for carry != 0 && j < 6 {
			t[j], carry = bits.Add64(t[j], 0, carry)
			j++
		}
	}

	// Step 3: u = t / R = t >> 192 = t[3..5]
	u := uint192{
		Lo:  t[3],
		Mid: t[4],
		Hi:  t[5],
	}

	// Step 4: conditional subtraction
	if cmp192(u, n) >= 0 {
		u = sub192(u, n)
	}
	return u
}

/* Still half-baked
// assumes N > 1 is odd
func is_prp(N uint192) int {
	n := uint64(N); 		//convert input to unsigned 64-bit int
	if N != n return 1; 	// declare everything >= 2^64 to be probably prime
	q := n>>1; // since N odd, this right bit-shift is q = (N-1)/2
	k := __builtin_ctzl(q); //int, k = 2-adic valuation of q
	q <<= __builtin_clzl(q); // shift q left by number of leading zeros in q
	nbar := inv192(n); // inverse modulo 2**192 of n
	one := mont_one(n); // 2**192 modulo n
	minusone := n-one; // mont -1
	uint64_t x = addmod64(one,one,n);
	for (q<<=1;q;q<<=1) {
		x = mulredc64(x,x,n,nbar);
		if ((int64_t)q < 0) x = addmod64(x,x,n);
	}
	if (x == one || x == minusone) return 1;
	while (--k >= 0) {
		x = mulredc64(x,x,n,nbar);
		if (x == minusone) return 1;
	}
	return 0;
}
*/
