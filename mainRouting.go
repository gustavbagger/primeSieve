package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"time"

	pr "github.com/fxtlabs/primes"
)

func (cfg *Config) computePrimeCutoff(boundLog float64, primeList []int, logs []float64) int {
	s := bestS[cfg.omega]

	// log of product of smallest ω−1 primes
	baseLog := 0.0
	for i := 0; i < cfg.omega-1; i++ {
		baseLog += logs[i]
	}

	// Try primes from largest to smallest
	for idx := len(primeList) - 1; idx >= cfg.omega-1; idx-- {
		p := primeList[idx]
		logMin := baseLog + math.Log(float64(p))

		// Build hypothetical indexes: smallest ω−1 primes + p
		indexes := make([]int, cfg.omega)
		for i := 0; i < cfg.omega-1; i++ {
			indexes[i] = i
		}
		indexes[cfg.omega-1] = idx

		sieveBound := pSieveLog(cfg.omega, s, indexes, primeList)
		effectiveBound := math.Min(boundLog, sieveBound)

		if logMin <= effectiveBound {
			return p
		}
	}

	// fallback: only smallest ω primes fit
	return primeList[cfg.omega-1]
}

func printIntervals(omegaMax, omegaMin int) {
	//careful here, we need len(primeList)>= omega
	primeList := pr.Sieve(500)
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

func search(omega, a, b int, path string) {
	file, _ := os.Create(path)

	w := bufio.NewWriterSize(file, 16*1024*1024) // 16MB buffer
	buf := make([]byte, omega*2)                 // Reusable 66‑byte buffer for one 33‑element slice
	cfg := Config{buf: buf, w: w, start: time.Now(), count: 0, omega: omega}

	boundLog := math.Log(float64(a) * math.Pow10(b))
	fullPrimeList := pr.Sieve(1000000)
	logs := make([]float64, len(fullPrimeList))

	for i, p := range fullPrimeList {
		logs[i] = math.Log(float64(p))
	}

	initBestS(cfg.omega, fullPrimeList)

	cutoff := cfg.computePrimeCutoff(boundLog, fullPrimeList, logs)
	fmt.Println("Exact prime cutoff=", cutoff)

	limit := sort.SearchInts(fullPrimeList, cutoff+1)
	primeList := fullPrimeList[:limit]
	logs = logs[:limit]

	maxIndex := len(primeList) - 1
	indexes := make([]int, cfg.omega)

	var exponents []int = make([]int, cfg.omega)
	for i := range exponents {
		exponents[i] = 1
	}
	cfg.recursionIndex(
		0,
		maxIndex,
		boundLog,
		indexes,
		primeList,
		logs,
		0.0,
		exponents,
	)
	fmt.Println("--------------------------------------------")
	fmt.Println("total count: ", cfg.count)
	w.Flush()
	file.Close()

	//fmt.Println(ReadRange("./data.bin", 0, 1000,omega))

	end := time.Now()
	fmt.Println("Time elapsed: ", end.Sub(cfg.start))

}
