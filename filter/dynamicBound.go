package filter

import (
	"math"
)

// 1.7 works, 1.4193747081548531 works if g(p)<p^0.621526 (since h = 0.82 * p^1/4 is admissible)
var C float64 = 1.4193747081548531

var PrimeListUpperBound int = 1 << 20

func DeltaSum(list []int) float64 {
	var sum float64
	for _, p := range list {
		sum -= 1.0 / float64(p)
	}
	return sum + 1
}

func InitBestS(omegaMax int, primeList []int) []int {
	bestS := make([]int, omegaMax+1)
	for omega := 1; omega <= omegaMax; omega++ {
		sBest := 0
		currentBest := float64(int(1) << (omega + 1))
		for s := 1; s <= omega; s++ {
			delta := DeltaSum(primeList[omega-s : omega])
			if delta <= 0.0 {
				break
			}
			currentTry := (2.0 + float64(s-1)/delta) * float64(int(1)<<(omega-s)) * math.Sqrt(2*C)
			if currentTry < currentBest {
				currentBest = currentTry

				sBest = s

			}
		}
		bestS[omega] = sBest
	}
	return bestS
}

func OptSieveBoundLog(omega, s int, indexes, primeList []int, boundLog float64) float64 {
	if s == 0 || len(indexes) < s {
		return boundLog
	}
	return math.Min(boundLog, PSieveLog(omega, s, indexes, primeList))
}

func PSieveLog(omega, s int, indexes, primeList []int) float64 {
	last := make([]int, s)
	for i := 0; i < s; i++ {
		last[i] = primeList[indexes[omega-s+i]]
	}
	delta := DeltaSum(last)
	if delta <= 0.0 {
		return 0.0
	}
	return (16.0)*(math.Log(2*delta+float64(s-1))+float64(omega-s)*math.Log(2.0)-math.Log(delta)) + 8.0*math.Log(2.0*C)
}
