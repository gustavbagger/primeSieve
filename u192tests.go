package main

import (
	"fmt"
	"math/big"
)

// convert uint192 → *big.Int
func uint192ToBig(x uint192) *big.Int {
	b := new(big.Int)
	b.Or(b, new(big.Int).Lsh(big.NewInt(0).SetUint64(x.Hi), 128))
	b.Or(b, new(big.Int).Lsh(big.NewInt(0).SetUint64(x.Mid), 64))
	b.Or(b, new(big.Int).SetUint64(x.Lo))
	return b
}

func testPRP(n uint192) {
	// 1) test using your 192‑bit PRP
	prp192 := isPRP192(n)

	// 2) convert to big.Int
	bn := uint192ToBig(n)

	// 3) test using math/big
	prpBig := bn.ProbablyPrime(20)

	// 4) print results
	fmt.Printf("n = %s\n", bn.String())
	fmt.Printf("isPRP192:   %v\n", prp192)
	fmt.Printf("big.Int PRP: %v\n", prpBig)
}

func testMulRedc(n uint192) {
	inv := inv64(n.Lo)
	npi := ^inv + 1

	for a := uint64(1); a < 32; a++ {
		for b := uint64(1); b < 32; b++ {
			A := uint192{Lo: a}
			B := uint192{Lo: b}

			got := MulRedc192(A, B, n, npi)
			want := uint192{Lo: (a * b) % n.Lo}

			if cmp192(got, want) != 0 {
				println("n =", n.Lo, "a =", a, "b =", b,
					"got =", got.Lo, "want =", want.Lo)
			}
		}
	}
}
