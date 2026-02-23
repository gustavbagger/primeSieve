package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	pr "github.com/fxtlabs/primes"
)

func deltaSum(list []int) float64 {
	var sum float64
	for _, p := range list {
		sum -= 1.0 / float64(p)
	}
	return sum + 1
}

func Prod(list []int) float64 {
	prod := 1.0
	for _, p := range list {
		prod *= float64(p)
	}
	return prod
}

func pritnIntervals() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("wrong arguments")
		return
	}
	var omegaMax int
	var err error
	omegaMax, err = strconv.Atoi(os.Args[1])
	if err != nil {
		return
	}
	omegaMin := omegaMax
	if len(os.Args) == 3 {
		omegaMin, err = strconv.Atoi(os.Args[2])
		if err != nil {
			return
		}
	}
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
		fmt.Printf(
			"%5.1e>p>%5.1e ---- o,s = %v,%v\n",
			math.Pow(currentBest, 16),
			Prod(primeList[:omega]),
			omega,
			sBest,
		)
	}
}
