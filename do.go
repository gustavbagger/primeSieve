package main

import (
	"fmt"
	"math"
	"sort"

	pr "github.com/fxtlabs/primes"
)

func computePrimeCutoff(omega int, boundLog float64, primeList []int, logs []float64) int {
	s := bestS[omega]

	// log of product of smallest ω−1 primes
	baseLog := 0.0
	for i := 0; i < omega-1; i++ {
		baseLog += logs[i]
	}

	// Try primes from largest to smallest
	for idx := len(primeList) - 1; idx >= omega-1; idx-- {
		p := primeList[idx]
		logMin := baseLog + math.Log(float64(p))

		// Build hypothetical indexes: smallest ω−1 primes + p
		indexes := make([]int, omega)
		for i := 0; i < omega-1; i++ {
			indexes[i] = i
		}
		indexes[omega-1] = idx

		sieveBound := pSieveLog(omega, s, indexes, primeList)
		effectiveBound := math.Min(boundLog, sieveBound)

		if logMin <= effectiveBound {
			return p
		}
	}

	// fallback: only smallest ω primes fit
	return primeList[omega-1]
}

func printIntervals(omegaMax, omegaMin int) {
	//careful here, we need len(primeList)>= omega
	primeList := pr.Sieve(1000)
	for omega := omegaMax; omega >= omegaMin; omega-- {
		sBest := 0
		deltaBest := 1.0
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
				deltaBest = delta
			}
		}
		var sum float64
		for _, p := range primeList[:omega] {
			sum += math.Log(float64(p))
		}

		fmt.Printf(
			"%5.1e>p>%5.1e ---- o,s,d = %v,%v,%2.2f\n",
			math.Pow(currentBest, 16),
			math.Pow(math.E, sum),
			omega,
			sBest,
			deltaBest,
		)

	}
}

var totalTopLevel int
var doneTopLevel int
var nextPercent int

func search(omega, a, b int) {

	boundLog := math.Log(float64(a) * math.Pow10(b))
	fullPrimeList := pr.Sieve(5000)
	logs := make([]float64, len(fullPrimeList))

	for i, p := range fullPrimeList {
		logs[i] = math.Log(float64(p))
	}

	initBestS(omega, fullPrimeList)

	cutoff := computePrimeCutoff(omega, boundLog, fullPrimeList, logs)
	fmt.Println("Exact prime cutoff=", cutoff)

	limit := sort.SearchInts(fullPrimeList, cutoff+1)
	primeList := fullPrimeList[:limit]
	logs = logs[:limit]

	maxIndex := len(primeList) - 1
	indexes := make([]int, omega)

	totalTopLevel = maxIndex - (omega - 1)
	doneTopLevel = 0
	nextPercent = 1

	recursiveLoop(
		0,
		omega,
		maxIndex,
		boundLog,
		indexes,
		primeList,
		logs,
		0.0,
	)
	fmt.Println("--------------------------------------------")
}
