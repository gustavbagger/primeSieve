package depricated

import (
	u256 "github.com/holiman/uint256"
)

type u256L struct {
	u256.Int
}

type mont256 struct {
	n  uint      // m.BitLen()
	m  *u256.Int //modulus, must be odd
	r2 *u256.Int // (1<<2n) mod m
}

func isEven(n *u256.Int) bool {
	return n[0]&1 == 0
}

// mont constructor
func newMont256(m *u256.Int) *mont256 {
	//check m odd
	if !isEven(m) {
		return nil
	}

	n := uint(m.BitLen())
	x := u256.NewInt(1)
	x.Sub(x.Lsh(x, n), m)
	return &mont256{n, new(u256.Int).Set(m), x.Mod(x.Mul(x, x), m)}
}

func (m mont256) reduce256(t *u256.Int) *u256.Int {
	a := new(u256.Int).Set(t)
	//loop over bitlength
	for i := uint(0); i < m.n; i++ {
		if !isEven(a) {
			a.Add(a, m.m)
		}
		a.Rsh(a, 1)
	}
	if a.Cmp(m.m) >= 0 {
		a.Sub(a, m.m)
	}
	return a
}

// reduce in place
func (m mont256) reduceInPlace256(t *u256.Int) *u256.Int {
	//loop over bitlength
	for i := uint(0); i < m.n; i++ {
		if !isEven(t) {
			t.Add(t, m.m)
		}
		t.Rsh(t, 1)
	}
	if t.Cmp(m.m) >= 0 {
		t.Sub(t, m.m)
	}
	return t
}

// 4-word REDC non-interleaved, see Brent-Zimmerman Alg 2.6
func MontMul(A, B *u256L) *u256L {
	return nil
}
