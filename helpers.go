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

func Prod(list []int) int {
	prod := 1
	for _, p := range list {
		prod *= p
	}
	return prod
}

func printIntervals() {
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

func indexesToValues(indexes, allValues []int) []int {
	output := make([]int, len(indexes))
	for i, index := range indexes {
		output[i] = allValues[index]
	}
	return output
}

func validIndexSet(indexes, allValues []int) bool {
	return pr.IsPrime(Prod(indexesToValues(indexes, allValues)) + 1)
}

func treeSearch(position, currentProduct, bound int, indexes, currentGuess, allValues []int) {
	if position == len(indexes) {
		if validIndexSet(currentGuess, allValues) {
			fmt.Println("--------------------------------------------")
			fmt.Printf("%v\n", indexesToValues(currentGuess, allValues))
		}
		return
	}

	index := indexes[position]
	i := allValues[index]

	nextProduct := currentProduct * i
	if nextProduct > bound {
		return
	}

	currentGuess = append(currentGuess, index)
	currentProduct = nextProduct

	treeSearch(position+1, currentProduct, bound, indexes, currentGuess, allValues)

	for {
		nextProduct = currentProduct * i
		if nextProduct > bound {
			return
		}
		currentGuess = append(currentGuess, index)
		currentProduct = nextProduct

		treeSearch(position+1, currentProduct, bound, indexes, currentGuess, allValues)
	}
}

func suffSmall(bound int, indexes, allValues []int) bool {
	return Prod(indexesToValues(indexes, allValues)) <= bound
}

type Status int

const (
	Continue Status = iota
	Backtrack
	Stop
)

func recursiveLoop(currentDepth, maxDepth, maxIndex, upperBound int, indexes, primeList []int) Status {
	if currentDepth == maxDepth {
		fmt.Println(suffSmall(upperBound, indexes, primeList))
		treeSearch(0, 1, upperBound, indexes, []int{}, primeList)
		return Continue
	}
	startIndex := 0
	if currentDepth > 0 {
		startIndex = indexes[currentDepth-1] + 1
	}

	for i := startIndex; i <= maxIndex-(maxDepth-currentDepth); i++ {
		indexes[currentDepth] = i

		if !suffSmall(upperBound, indexes, primeList) {
			return Backtrack
		}

		status := recursiveLoop(currentDepth+1, maxDepth, maxIndex, upperBound, indexes, primeList)

		switch status {
		case Continue:
			continue
		case Backtrack:
			continue
		case Stop:
			return Stop
		}
	}
	return Backtrack
}
