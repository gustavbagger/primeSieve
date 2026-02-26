package main

import "math/big"

func findPrimitive(
	p *big.Int,
	indexes []int,
	exponents []int,
	primeList []int,
) int {
	for a := 1; a <= 100; a++ {
		bigA := big.NewInt(int64(a))
		for i, index := range indexes {
			e := big.NewInt(1)
			for j := 0; j < len(indexes); i++ {
				dummy := big.NewInt(1)
				if j != i {
					for pExp := 1; pExp <= exponents[j]; pExp++ {
						e.Mul(e, dummy.Exp(bigA, big.NewInt(int64(primeList[index])), p))
						e.Mod(e, p)
					}
				} else {
					for pExp := 1; pExp <= exponents[j]-1; pExp++ {
						e.Mul(e, dummy.Exp(bigA, big.NewInt(int64(primeList[index])), p))
						e.Mod(e, p)
					}
				}
			}
			if e.Cmp(big.NewInt(1)) == 0 {
				break
			}
			if i == len(indexes) {
				return a
			}
		}

	}
	return 101
}
