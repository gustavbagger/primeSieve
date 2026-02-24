package main

import (
	"fmt"
	"math"
	"math/big"
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
	prod := big.NewInt(1)
	for _, index := range indexes {
		prod.Mul(prod, big.NewInt(int64(allValues[index])))
	}
	prod.Add(prod, big.NewInt(1))
	return prod.ProbablyPrime(32)
}

func treeSearch(
	position int,
	currentLog, boundLog float64,
	currentProd *big.Int,
	indexes, currentGuess, allValues []int,
	logs []float64,
) {
	if position == len(indexes) {
		if validIndexSet(currentGuess, allValues) {
			fmt.Println("--------------------------------------------")
			fmt.Printf("%v\n", indexesToValues(currentGuess, allValues))
		}
		return
	}

	index := indexes[position]
	val := allValues[index]
	valLog := logs[index]

	if currentLog+valLog > boundLog {
		return
	}

	nextGuess := append(append([]int{}, currentGuess...), index)
	nextProd := new(big.Int).Mul(currentProd, big.NewInt(int64(val)))

	treeSearch(
		position+1,
		currentLog+valLog,
		boundLog,
		nextProd,
		indexes,
		nextGuess,
		allValues,
		logs,
	)

	logAcc := currentLog + valLog
	prodAcc := new(big.Int).Set(nextProd)
	guessAcc := append([]int{}, nextGuess...)

	for {
		logAcc += valLog
		if logAcc > boundLog {
			return
		}

		guessAcc = append(guessAcc, index)
		prodAcc.Mul(prodAcc, big.NewInt(int64(val)))

		treeSearch(
			position+1,
			logAcc,
			boundLog,
			new(big.Int).Set(prodAcc),
			indexes,
			guessAcc,
			allValues,
			logs,
		)
	}
}

func suffSmallLog(boundLog float64, indexes []int, logs []float64) bool {
	var sum float64
	for _, index := range indexes {
		sum += logs[index]
	}
	return sum <= boundLog
}

type Status int

const (
	Continue Status = iota
	Backtrack
	Stop
)

func recursiveLoop(
	currentDepth, maxDepth, maxIndex int,
	boundLog float64,
	indexes, primeList []int,
	logs []float64,
) Status {

	if currentDepth == maxDepth {
		fmt.Println(suffSmallLog(boundLog, indexes, logs))
		treeSearch(
			0,
			0,
			boundLog,
			big.NewInt(1),
			indexes,
			[]int{},
			primeList,
			logs)
		return Continue
	}
	startIndex := 0
	if currentDepth > 0 {
		startIndex = indexes[currentDepth-1] + 1
	}

	for i := startIndex; i <= maxIndex-(maxDepth-currentDepth); i++ {
		indexes[currentDepth] = i

		if !suffSmallLog(boundLog, indexes[:currentDepth+1], logs) {
			return Backtrack
		}

		status := recursiveLoop(
			currentDepth+1,
			maxDepth,
			maxIndex,
			boundLog,
			indexes,
			primeList,
			logs)

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
