package main

import "math/big"

func validExponentSet(indexes, exponents, allValues []int) (*big.Int, bool) {
	prod := big.NewInt(1)
	for i, index := range indexes {
		prod.Mul(
			prod,
			new(big.Int).Exp(
				big.NewInt(int64(allValues[index])),
				big.NewInt(int64(exponents[i])),
				nil),
		)
	}
	prod.Add(prod, big.NewInt(1))
	return prod, prod.ProbablyPrime(32)
}

/* Depricated
// assuming indexes
func expSetStrongPRP2(indexes, exponents, primeList []int) (*big.Int, bool) {

	N := big.NewInt(1)
	tmp := new(big.Int)

	for i := range indexes {
		p := big.NewInt(int64(primeList[indexes[i]]))
		tmp.Exp(p, big.NewInt(int64(exponents[i])), nil)
		N.Mul(N, tmp)
	}
	s := 0
	if indexes[0] == 0 { //if 2 is in the product, set s to be its exponent
		s = exponents[0]
	}

	d := new(big.Int).Rsh(N, uint(s)) //odd part of N

	return montStrongPRP2(s, d, N)

}

// assuming indexes
func expSetStrongPRP2u192(indexes, exponents, primeList []int) (*Int, bool) {

	N := NewInt(1)
	tmp := new(Int)

	for i := range indexes {
		p := NewInt(uint64(primeList[indexes[i]]))
		tmp.Exp(p, NewInt(uint64(exponents[i])))
		N.Mul(N, tmp)
	}
	s := 0
	if indexes[0] == 0 { //if 2 is in the product, set s to be its exponent
		s = exponents[0]
	}

	d := new(Int).Rsh(N, uint(s)) //odd part of N

	return montStrongPRP2u192(s, d, N)

}

*/
