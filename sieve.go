package main

import "math"

var bestS []int

func initBestS(omegaMax int, primeList []int) {
	bestS = make([]int, omegaMax+1)
	for omega := 1; omega <= omegaMax; omega++ {
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
		bestS[omega] = sBest
	}
}

func optSieveBoundLog(omega int, indexes, primeList []int, boundLog float64) float64 {
	s := bestS[omega]
	if s == 0 || len(indexes) < s {
		return boundLog
	}
	return math.Min(boundLog, pSieveLog(omega, s, indexes, primeList))
}

func pSieveLog(omega, s int, indexes, primeList []int) float64 {
	last := make([]int, s)
	for i := 0; i < s; i++ {
		last[i] = primeList[indexes[omega-s+i]]
	}
	delta := deltaSum(last)
	if delta <= 0.0 {
		return 0.0
	}
	return 16 * (math.Log(2*delta+float64(s-1)) + float64(omega+1-s)*math.Log(2.0) - math.Log(delta))
}
