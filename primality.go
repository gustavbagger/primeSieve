package main

import (
	"fmt"
	"math/big"
	"time"
)

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

func validExponentSet192(indexes, exponents, allValues []int) (uint192, bool) {
	prod := uint192{Lo: 1}
	for i, index := range indexes {
		for exp := 1; exp <= exponents[i]; exp++ {
			prod = mulMod192(prod, uint192{Lo: uint64(allValues[index])})
		}
	}
	prod = add192(prod, uint192{Lo: 1})
	prp := strongPRP(prod)
	return prod, prp
}

func (cfg *Config) handleSuccess(indexes, exponents []int) {
	cfg.count++
	if cfg.count%100000 == 0 {
		fmt.Printf("%.2e values found - expect 10^8 (for o=33).\n", float64(cfg.count))
		fmt.Println("Total time: ", time.Now().Sub(cfg.start))
	}

	cfg.WriteToBin(indexes, exponents)
}
