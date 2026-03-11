package depricated

import (
	"math/big"
)

type mont struct {
	n  uint     // m.BitLen()
	m  *big.Int // modulus, must be odd
	r2 *big.Int // (1<<2n) mod m
}

func newMont(m *big.Int) *mont {
	if m.Bit(0) != 1 { //if m even, return nil
		return nil
	}
	n := uint(m.BitLen())
	x := big.NewInt(1)
	x.Sub(x.Lsh(x, n), m)                                       //set x to x<<n, then subtract m
	return &mont{n, new(big.Int).Set(m), x.Mod(x.Mul(x, x), m)} // r2 = square x, then reduce modulo m
}

func (m mont) reduce(t *big.Int) *big.Int {
	a := new(big.Int).Set(t)
	for i := uint(0); i < m.n; i++ {
		if a.Bit(0) == 1 {
			a.Add(a, m.m)
		}
		a.Rsh(a, 1)
	}
	if a.Cmp(m.m) >= 0 {
		a.Sub(a, m.m)
	}
	return a
}

func (m *mont) montMul(x, y *big.Int) *big.Int {
	t := new(big.Int).Mul(x, y)
	return m.reduce(t)
}

func (m *mont) toMont(x *big.Int) *big.Int {
	t := new(big.Int).Mul(x, m.r2)
	return m.reduce(t)
}

func (m *mont) fromMont(x *big.Int) *big.Int {
	one := big.NewInt(1)
	return m.montMul(x, one)
}

func (m *mont) montPow(base, exp *big.Int) *big.Int {
	one := big.NewInt(1)
	result := m.toMont(one)
	x := m.toMont(base)
	for i := exp.BitLen() - 1; i >= 0; i-- {
		result = m.montMul(result, result)
		if exp.Bit(i) == 1 {
			result = m.montMul(result, x)
		}
	}
	return m.fromMont(result)
}

// Assumed even N>2
func montStrongPRP2(diaticValuation int, oddPart, N *big.Int) (*big.Int, bool) {
	p := new(big.Int).Add(N, big.NewInt(1))

	m := newMont(p) //setting up a new Montgomery struct

	x := m.montPow(big.NewInt(2), oddPart) // x = 2^d modulo p

	if x.Cmp(big.NewInt(1)) == 0 || x.Cmp(N) == 0 { // check if x == +- 1 modulo p
		return p, true
	}

	for r := 1; r < diaticValuation; r++ {
		x = m.montMul(x, x)
		if x.Cmp(N) == 0 {
			return p, true
		}
	}
	return nil, false
}

type mont192 struct {
	n  uint // m.BitLen()
	m  *Int // modulus, must be odd
	r2 *Int // (1<<2n) mod m
}

func newMont192(m *Int) *mont192 {
	if (m[0]>>1)<<1 == m[0] { //if m even, return nil
		return nil
	}
	n := uint(m.BitLen())
	x := NewInt(1)
	x.SubOverflow(x.Lsh(x, n), m)                              //set x to x<<n, then subtract m
	return &mont192{n, new(Int).Set(m), x.Mod(x.Mul(x, x), m)} // r2 = square x, then reduce modulo m
}

func (m mont192) reduce192(t *Int) *Int {
	a := new(Int).Set(t)
	for i := uint(0); i < m.n; i++ {
		if a.Bit(0) {
			a.AddOverflow(a, m.m)
		}
		a.Rsh(a, 1)
	}
	if a.Cmp(m.m) >= 0 {
		a.SubOverflow(a, m.m)
	}
	return a
}

func (m *mont192) montMul192(x, y *Int) *Int {
	t := new(Int).Mul(x, y)
	return m.reduce192(t)
}

func (m *mont192) toMont192(x *Int) *Int {
	t := new(Int).Mul(x, m.r2)
	return m.reduce192(t)
}

func (m *mont192) fromMont192(x *Int) *Int {
	one := NewInt(1)
	return m.montMul192(x, one)
}

func (m *mont192) montPow192(base, exp *Int) *Int {
	one := NewInt(1)
	result := m.toMont192(one)
	x := m.toMont192(base)
	for i := exp.BitLen() - 1; i >= 0; i-- {
		result = m.montMul192(result, result)
		if exp.Bit(i) {
			result = m.montMul192(result, x)
		}
	}
	return m.fromMont192(result)
}

// Assumed even N>2
func montStrongPRP2u192(diaticValuation int, oddPart, N *Int) (*Int, bool) {
	p, _ := new(Int).AddOverflow(N, NewInt(1))

	m := newMont192(p) //setting up a new Montgomery struct

	x := m.montPow192(NewInt(2), oddPart) // x = 2^d modulo p

	if x.Eq(NewInt(1)) || x.Eq(N) { // check if x == +- 1 modulo p
		return p, true
	}

	for r := 1; r < diaticValuation; r++ {
		x = m.montMul192(x, x)
		if x.Eq(N) {
			return p, true
		}
	}
	return nil, false
}
