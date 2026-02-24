package main

import (
	"fmt"
	"math"

	pr "github.com/fxtlabs/primes"
)

func printIntervals(omegaMax, omegaMin int) {
	//careful here, we need len(primeList)>= omega
	primeList := pr.Sieve(173)
	for omega := omegaMax; omega >= omegaMin; omega-- {
		sBest := 0
		currentBest := float64(int(1) << (omega + 1))
		for s := 1; s <= omega; s++ {
			delta := deltaSum(primeList[omega-s : omega])
			if delta <= 0.0 {
				break
			}
			currentTry := (2.0 + float64(s-1)/delta) * float64(int(1)<<(omega+1-s))
			if currentTry < currentBest {
				currentBest = currentTry
				sBest = s
			}
		}
		var sum float64
		for _, p := range primeList[:omega] {
			sum += math.Log(float64(p))
		}

		fmt.Printf(
			"%5.1e>p>%5.1e ---- o,s = %v,%v\n",
			math.Pow(currentBest, 16),
			math.Pow(math.E, sum),
			omega,
			sBest,
		)

	}
}

func search(omega, a, b int) {

	boundLog := math.Log(float64(a) * math.Pow10(b))
	primeList := pr.Sieve(2000)
	logs := make([]float64, len(primeList))

	for i, p := range primeList {
		logs[i] = math.Log(float64(p))
	}

	maxIndex := len(primeList) - 1
	indexes := make([]int, omega)

	recursiveLoop(
		0,
		omega,
		maxIndex,
		boundLog,
		indexes,
		primeList,
		logs,
	)
	//fmt.Println(logs)
	fmt.Println("--------------------------------------------")
}
