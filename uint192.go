// uint256: Fixed size 256-bit math library
// Copyright 2018-2020 uint256 Authors
// SPDX-License-Identifier: BSD-3-Clause

// Package math provides integer math utilities.

package main

import (
	"math/bits"
)

// Int is represented as an array of 3 uint64, in little-endian order,
// so that Int[2] is the most significant, and Int[0] is the least significant
type Int [3]uint64

// NewInt returns a new initialized Int.
func NewInt(val uint64) *Int {
	z := &Int{}
	z.SetUint64(val)
	return z
}

// Clone creates a new Int identical to z
func (z *Int) Clone() *Int {
	return &Int{z[0], z[1], z[2]}
}

// AddOverflow sets z to the sum x+y, and returns z and whether overflow occurred
func (z *Int) AddOverflow(x, y *Int) (*Int, bool) {
	var carry uint64
	z[0], carry = bits.Add64(x[0], y[0], 0)
	z[1], carry = bits.Add64(x[1], y[1], carry)
	z[2], carry = bits.Add64(x[2], y[2], carry)
	return z, carry != 0
}

// AddMod sets z to the sum ( x+y ) mod m, and returns z.
// If m == 0, z is set to 0 (OBS: differs from the big.Int)
func (z *Int) AddMod(x, y, m *Int) *Int {

	// Fast path for m >= 2^128, with x and y at most slightly bigger than m.
	// This is always the case when x and y are already reduced modulo such m.

	if (m[2] != 0) && (x[2] <= m[2]) && (y[2] <= m[2]) {
		var (
			gteC1 uint64
			gteC2 uint64
			tmpX  Int
			tmpY  Int
			res   Int
		)

		// reduce x/y modulo m if they are gte m
		tmpX[0], gteC1 = bits.Sub64(x[0], m[0], gteC1)
		tmpX[1], gteC1 = bits.Sub64(x[1], m[1], gteC1)
		tmpX[2], gteC1 = bits.Sub64(x[2], m[2], gteC1)

		tmpY[0], gteC2 = bits.Sub64(y[0], m[0], gteC2)
		tmpY[1], gteC2 = bits.Sub64(y[1], m[1], gteC2)
		tmpY[2], gteC2 = bits.Sub64(y[2], m[2], gteC2)

		if gteC1 == 0 {
			x = &tmpX
		}
		if gteC2 == 0 {
			y = &tmpY
		}
		var (
			c1  uint64
			c2  uint64
			tmp Int
		)

		res[0], c1 = bits.Add64(x[0], y[0], c1)
		res[1], c1 = bits.Add64(x[1], y[1], c1)
		res[2], c1 = bits.Add64(x[2], y[2], c1)

		tmp[0], c2 = bits.Sub64(res[0], m[0], c2)
		tmp[1], c2 = bits.Sub64(res[1], m[1], c2)
		tmp[2], c2 = bits.Sub64(res[2], m[2], c2)

		// final sub was unnecessary
		if c1 == 0 && c2 != 0 {
			return z.Set(&res)
		}

		return z.Set(&tmp)
	}

	if m.IsZero() {
		return z.Clear()
	}
	if z == m { // z is an alias for m and will be overwritten by AddOverflow before m is read
		m = m.Clone()
	}
	if _, overflow := z.AddOverflow(x, y); overflow {
		sum := [4]uint64{z[0], z[1], z[2], 1}
		var quot [4]uint64
		var rem Int
		udivrem(quot[:], sum[:], m, &rem)
		return z.Set(&rem)
	}
	return z.Mod(z, m)
}

// SubUint64 set z to the difference x - y, where y is a uint64, and returns z
func (z *Int) SubUint64(x *Int, y uint64) *Int {
	var carry uint64
	z[0], carry = bits.Sub64(x[0], y, carry)
	z[1], carry = bits.Sub64(x[1], 0, carry)
	z[2], _ = bits.Sub64(x[2], 0, carry)
	return z
}

// SubOverflow sets z to the difference x-y and returns z and true if the operation underflowed
func (z *Int) SubOverflow(x, y *Int) (*Int, bool) {
	var carry uint64
	z[0], carry = bits.Sub64(x[0], y[0], 0)
	z[1], carry = bits.Sub64(x[1], y[1], carry)
	z[2], carry = bits.Sub64(x[2], y[2], carry)
	return z, carry != 0
}

// umulStep computes (hi * 2^64 + lo) = z + (x * y) + carry.
func umulStep(z, x, y, carry uint64) (hi, lo uint64) {
	hi, lo = bits.Mul64(x, y)
	lo, carry = bits.Add64(lo, carry, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	lo, carry = bits.Add64(lo, z, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return hi, lo
}

// umulHop computes (hi * 2^64 + lo) = z + (x * y)
func umulHop(z, x, y uint64) (hi, lo uint64) {
	hi, lo = bits.Mul64(x, y)
	lo, carry := bits.Add64(lo, z, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return hi, lo
}

// CHECK THIS
// Mul sets z to the product x*y
func (z *Int) Mul(x, y *Int) *Int {
	var (
		carry0, carry1 uint64
		res1, res2     uint64
		x0, x1, x2     = x[0], x[1], x[2]
		y0, y1, y2     = y[0], y[1], y[2]
	)

	carry0, z[0] = bits.Mul64(x0, y0)
	carry0, res1 = umulHop(carry0, x1, y0)
	carry0, res2 = umulHop(carry0, x2, y0)

	carry1, z[1] = umulHop(res1, x0, y1)
	carry1, res2 = umulStep(res2, x1, y1, carry1)

	_, z[2] = umulHop(res2, x0, y2)

	return z
}

// IsUint64 reports whether z can be represented as a uint64.
func (z *Int) IsUint64() bool {
	return (z[1] | z[2]) == 0
}

// Uint64 returns the lower 64-bits of z
func (z *Int) Uint64() uint64 {
	return z[0]
}

// reciprocal2by1 computes <^d, ^0> / d.
func reciprocal2by1(d uint64) uint64 {
	reciprocal, _ := bits.Div64(^d, ^uint64(0), d)
	return reciprocal
}

// udivrem2by1 divides <uh, ul> / d and produces both quotient and remainder.
// It uses the provided d's reciprocal.
// Implementation ported from https://github.com/chfast/intx and is based on
// "Improved division by invariant integers", Algorithm 4.
func udivrem2by1(uh, ul, d, reciprocal uint64) (quot, rem uint64) {
	qh, ql := bits.Mul64(reciprocal, uh)
	ql, carry := bits.Add64(ql, ul, 0)
	qh, _ = bits.Add64(qh, uh, carry)
	qh++

	r := ul - qh*d

	if r > ql {
		qh--
		r += d
	}

	if r >= d {
		qh++
		r -= d
	}

	return qh, r
}

// udivremBy1 divides u by single normalized word d and produces both quotient and remainder.
// The quotient is stored in provided quot.
func udivremBy1(quot, u []uint64, d uint64) (rem uint64) {
	reciprocal := reciprocal2by1(d)
	rem = u[len(u)-1] // Set the top word as remainder.
	for j := len(u) - 2; j >= 0; j-- {
		quot[j], rem = udivrem2by1(rem, u[j], d, reciprocal)
	}
	return rem
}

// subMulTo computes x -= y * multiplier.
// Requires len(x) >= len(y) > 0.
func subMulTo(x, y []uint64, multiplier uint64) uint64 {
	var borrow uint64
	_ = x[len(y)-1] // bounds check hint to compiler; see golang.org/issue/14808
	for i := 0; i < len(y); i++ {
		s, carry1 := bits.Sub64(x[i], borrow, 0)
		ph, pl := bits.Mul64(y[i], multiplier)
		t, carry2 := bits.Sub64(s, pl, 0)
		x[i] = t
		borrow = ph + carry1 + carry2
	}
	return borrow
}

// addTo computes x += y.
// Requires len(x) >= len(y) > 0.
func addTo(x, y []uint64) uint64 {
	var carry uint64
	_ = x[len(y)-1] // bounds check hint to compiler; see golang.org/issue/14808
	for i := 0; i < len(y); i++ {
		x[i], carry = bits.Add64(x[i], y[i], carry)
	}
	return carry
}

// udivremKnuth implements the division of u by normalized multiple word d from the Knuth's division algorithm.
// The quotient is stored in provided quot - len(u)-len(d) words.
// Updates u to contain the remainder - len(d) words.
func udivremKnuth(quot, u, d []uint64) {
	dh := d[len(d)-1]
	dl := d[len(d)-2]
	reciprocal := reciprocal2by1(dh)

	for j := len(u) - len(d) - 1; j >= 0; j-- {
		u2 := u[j+len(d)]
		u1 := u[j+len(d)-1]
		u0 := u[j+len(d)-2]

		var qhat, rhat uint64
		if u2 >= dh { // Division overflows.
			qhat = ^uint64(0)
			// TODO: Add "qhat one to big" adjustment (not needed for correctness, but helps avoiding "add back" case).
		} else {
			qhat, rhat = udivrem2by1(u2, u1, dh, reciprocal)
			ph, pl := bits.Mul64(qhat, dl)
			if ph > rhat || (ph == rhat && pl > u0) {
				qhat--
				// TODO: Add "qhat one to big" adjustment (not needed for correctness, but helps avoiding "add back" case).
			}
		}

		// Multiply and subtract.
		borrow := subMulTo(u[j:], d, qhat)
		u[j+len(d)] = u2 - borrow
		if u2 < borrow { // Too much subtracted, add back.
			qhat--
			u[j+len(d)] += addTo(u[j:], d)
		}

		quot[j] = qhat // Store quotient digit.
	}
}

// udivrem divides u by d and produces both quotient and remainder.
// The quotient is stored in provided quot - len(u)-len(d)+1 words.
// It loosely follows the Knuth's division algorithm (sometimes referenced as "schoolbook" division) using 64-bit words.
// See Knuth, Volume 2, section 4.3.1, Algorithm D.
func udivrem(quot, u []uint64, d, rem *Int) {
	var dLen int
	for i := len(d) - 1; i >= 0; i-- {
		if d[i] != 0 {
			dLen = i + 1
			break
		}
	}

	shift := uint(bits.LeadingZeros64(d[dLen-1]))

	var dnStorage Int
	dn := dnStorage[:dLen]
	for i := dLen - 1; i > 0; i-- {
		dn[i] = (d[i] << shift) | (d[i-1] >> (64 - shift))
	}
	dn[0] = d[0] << shift

	var uLen int
	for i := len(u) - 1; i >= 0; i-- {
		if u[i] != 0 {
			uLen = i + 1
			break
		}
	}

	if uLen < dLen {
		if rem != nil {
			copy(rem[:], u)
		}
		return
	}

	var unStorage [9]uint64
	un := unStorage[:uLen+1]
	un[uLen] = u[uLen-1] >> (64 - shift)
	for i := uLen - 1; i > 0; i-- {
		un[i] = (u[i] << shift) | (u[i-1] >> (64 - shift))
	}
	un[0] = u[0] << shift

	// TODO: Skip the highest word of numerator if not significant.

	if dLen == 1 {
		r := udivremBy1(quot, un, dn[0])
		if rem != nil {
			rem.SetUint64(r >> shift)
		}
		return
	}

	udivremKnuth(quot, un, dn)

	if rem != nil {
		for i := 0; i < dLen-1; i++ {
			rem[i] = (un[i] >> shift) | (un[i+1] << (64 - shift))
		}
		rem[dLen-1] = un[dLen-1] >> shift
	}
}

// Mod sets z to the modulus x%y for y != 0 and returns z.
// If y == 0, z is set to 0 (OBS: differs from the big.Int)
func (z *Int) Mod(x, y *Int) *Int {
	if y.IsZero() || x.Eq(y) {
		return z.Clear()
	}
	if x.Lt(y) {
		return z.Set(x)
	}
	// At this point:
	// x != 0
	// y != 0
	// x > y

	// Shortcut trivial case
	if x.IsUint64() {
		return z.SetUint64(x.Uint64() % y.Uint64())
	}

	var quot, rem Int
	udivrem(quot[:], x[:], y, &rem)
	return z.Set(&rem)
}

// Neg returns -x mod 2**256.
func (z *Int) Neg(x *Int) *Int {
	z.SubOverflow(new(Int), x)
	return z
}

// BitLen returns the number of bits required to represent z
func (z *Int) BitLen() int {
	switch {
	case z[2] != 0:
		return 128 + bits.Len64(z[2])
	case z[1] != 0:
		return 64 + bits.Len64(z[1])
	default:
		return bits.Len64(z[0])
	}
}

func (z *Int) Bit(i int) bool {
	switch {
	case i >= 128:
		return (((z[2] >> (i - 128)) >> 1) << 1) != (z[2] >> (i - 128))
	case i >= 64:
		return (((z[1] >> (i - 64)) >> 1) << 1) != (z[1] >> (i - 64))
	default:
		return (((z[0] >> i) >> 1) << 1) != (z[0] >> i)
	}
}

func (z *Int) lsh64(x *Int) {
	z[2], z[1], z[0] = x[1], x[0], 0
}

func (z *Int) lsh128(x *Int) {
	z[2], z[1], z[0] = x[0], 0, 0
}

func (z *Int) rsh64(x *Int) {
	z[2], z[1], z[0] = 0, x[2], x[1]
}

func (z *Int) rsh128(x *Int) {
	z[2], z[1], z[0] = 0, 0, x[2]
}

// Gt returns true if z > x
func (z *Int) Gt(x *Int) bool {
	return x.Lt(z)
}

// Lt returns true if z < x
func (z *Int) Lt(x *Int) bool {
	// z < x <=> z - x < 0 i.e. when subtraction overflows.
	_, carry := bits.Sub64(z[0], x[0], 0)
	_, carry = bits.Sub64(z[1], x[1], carry)
	_, carry = bits.Sub64(z[2], x[2], carry)
	return carry != 0
}

// SetUint64 sets z to the value x
func (z *Int) SetUint64(x uint64) *Int {
	z[2], z[1], z[0] = 0, 0, x
	return z
}

// Eq returns true if z == x
func (z *Int) Eq(x *Int) bool {
	return ((z[0] ^ x[0]) | (z[1] ^ x[1]) | (z[2] ^ x[2])) == 0
}

// Cmp compares z and x and returns:
//
//	-1 if z <  x
//	 0 if z == x
//	+1 if z >  x
func (z *Int) Cmp(x *Int) (r int) {
	// z < x <=> z - x < 0 i.e. when subtraction overflows.
	d0, carry := bits.Sub64(z[0], x[0], 0)
	d1, carry := bits.Sub64(z[1], x[1], carry)
	d2, carry := bits.Sub64(z[2], x[2], carry)
	if carry == 1 {
		return -1
	}
	if d0|d1|d2 == 0 {
		return 0
	}
	return 1
}

// IsZero returns true if z == 0
func (z *Int) IsZero() bool {
	return (z[0] | z[1] | z[2]) == 0
}

// Clear sets z to 0
func (z *Int) Clear() *Int {
	z[2], z[1], z[0] = 0, 0, 0
	return z
}

// SetOne sets z to 1
func (z *Int) SetOne() *Int {
	z[2], z[1], z[0] = 0, 0, 1
	return z
}

// Lsh sets z = x << n and returns z.
func (z *Int) Lsh(x *Int, n uint) *Int {
	switch {
	case n == 0:
		return z.Set(x)
	case n >= 128:
		z.lsh128(x)
		n -= 128
		z[2] <<= n
		return z
	case n >= 64:
		z.lsh64(x)
		n -= 64
		z[2] = (z[2] << n) | (z[1] >> (64 - n))
		z[1] <<= n
		return z
	default:
		z.Set(x)
		z[2] = (z[2] << n) | (z[1] >> (64 - n))
		z[1] = (z[1] << n) | (z[0] >> (64 - n))
		z[0] <<= n
		return z
	}
}

// Rsh sets z = x >> n and returns z.
func (z *Int) Rsh(x *Int, n uint) *Int {
	switch {
	case n == 0:
		return z.Set(x)
	case n >= 128:
		z.rsh128(x)
		n -= 128
		z[0] >>= n

		return z
	case n >= 64:
		z.rsh64(x)
		n -= 64
		z[0] = (z[0] >> n) | (z[1] << (64 - n))
		z[1] >>= n
		return z
	default:
		z.Set(x)
		z[0] = (z[0] >> n) | (z[1] << (64 - n))
		z[1] = (z[1] >> n) | (z[2] << (64 - n))
		z[2] >>= n
		return z
	}
}

// Set sets z to x and returns z.
func (z *Int) Set(x *Int) *Int {
	z[0], z[1], z[2] = x[0], x[1], x[2]
	return z
}

// Or sets z = x | y and returns z.
func (z *Int) Or(x, y *Int) *Int {
	z[0] = x[0] | y[0]
	z[1] = x[1] | y[1]
	z[2] = x[2] | y[2]

	return z
}

// And sets z = x & y and returns z.
func (z *Int) And(x, y *Int) *Int {
	z[0] = x[0] & y[0]
	z[1] = x[1] & y[1]
	z[2] = x[2] & y[2]

	return z
}

// Xor sets z = x ^ y and returns z.
func (z *Int) Xor(x, y *Int) *Int {
	z[0] = x[0] ^ y[0]
	z[1] = x[1] ^ y[1]
	z[2] = x[2] ^ y[2]

	return z
}

func (z *Int) squared() {
	var (
		carry0, carry1   uint64
		res0, res1, res2 uint64
	)

	carry0, res0 = bits.Mul64(z[0], z[0])
	carry0, res1 = umulHop(carry0, z[0], z[1])
	carry0, res2 = umulHop(carry0, z[0], z[2])

	carry1, res1 = umulHop(res1, z[0], z[1])
	carry1, res2 = umulStep(res2, z[1], z[1], carry1)

	_, res2 = umulHop(res2, z[0], z[2])

	z[0], z[1], z[2] = res0, res1, res2
}

// Exp sets z = base**exponent mod 2**256, and returns z.
func (z *Int) Exp(base, exponent *Int) *Int {
	var (
		res        = Int{1, 0, 0}
		multiplier = *base
		expBitLen  = exponent.BitLen()
		curBit     = 0
		word       = exponent[0]
		even       = base[0]&1 == 0
	)
	if even && expBitLen > 8 {
		return z.Clear()
	}

	for ; curBit < expBitLen && curBit < 64; curBit++ {
		if word&1 == 1 {
			res.Mul(&res, &multiplier)
		}
		multiplier.squared()
		word >>= 1
	}
	if even { // If the base was even, we are finished now
		return z.Set(&res)
	}

	word = exponent[1]
	for ; curBit < expBitLen && curBit < 128; curBit++ {
		if word&1 == 1 {
			res.Mul(&res, &multiplier)
		}
		multiplier.squared()
		word >>= 1
	}

	word = exponent[2]
	for ; curBit < expBitLen && curBit < 192; curBit++ {
		if word&1 == 1 {
			res.Mul(&res, &multiplier)
		}
		multiplier.squared()
		word >>= 1
	}
	return z.Set(&res)
}

var (
	// pows64 contains 10^0 ... 10^19
	pows64 = [20]uint64{
		1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	}
	// pows contain 10 ** 20 ... 10 ** 57
	pows = [38]Int{
		{7766279631452241920, 5, 0}, {3875820019684212736, 54, 0}, {1864712049423024128, 542, 0}, {200376420520689664, 5421, 0}, {2003764205206896640, 54210, 0}, {1590897978359414784, 542101, 0}, {15908979783594147840, 5421010, 0}, {11515845246265065472, 54210108, 0}, {4477988020393345024, 542101086, 0}, {7886392056514347008, 5421010862, 0}, {5076944270305263616, 54210108624, 0}, {13875954555633532928, 542101086242, 0}, {9632337040368467968, 5421010862427, 0},
		{4089650035136921600, 54210108624275, 0}, {4003012203950112768, 542101086242752, 0}, {3136633892082024448, 5421010862427522, 0}, {12919594847110692864, 54210108624275221, 0}, {68739955140067328, 542101086242752217, 0}, {687399551400673280, 5421010862427522170, 0}, {6873995514006732800, 17316620476856118468, 2}, {13399722918938673152, 7145508105175220139, 29}, {4870020673419870208, 16114848830623546549, 293}, {11806718586779598848, 13574535716559052564, 2938},
		{7386721425538678784, 6618148649623664334, 29387}, {80237960548581376, 10841254275107988496, 293873}, {802379605485813760, 16178822382532126880, 2938735}, {8023796054858137600, 14214271235644855872, 29387358}, {6450984253743169536, 13015503840481697412, 293873587}, {9169610316303040512, 1027829888850112811, 2938735877}, {17909126868192198656, 10278298888501128114, 29387358770}, {13070572018536022016, 10549268516463523069, 293873587705}, {1578511669393358848, 13258964796087472617, 2938735877055}, {15785116693933588480, 3462439444907864858, 29387358770557},
		{10277214349659471872, 16177650375369096972, 293873587705571}, {10538423128046960640, 14202551164014556797, 2938735877055718}, {13150510911921848320, 12898303124178706663, 29387358770557187}, {2377900603251621888, 18302566799529756941, 293873587705571876}, {5332261958806667264, 17004971331911604867, 2938735877055718769},
	}
)
