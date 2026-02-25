package main

import (
	"fmt"
	"math/big"
)

func deltaSum(list []int) float64 {
	var sum float64
	for _, p := range list {
		sum -= 1.0 / float64(p)
	}
	return sum + 1
}

func Prod(list []int) *big.Int {
	prod := big.NewInt(1)
	for _, p := range list {
		prod.Mul(prod, big.NewInt(int64(p)))
	}
	return prod
}

func indexesToValues(indexes, allValues []int) []int {
	output := make([]int, len(indexes))
	for i, index := range indexes {
		output[i] = allValues[index]
	}
	return output
}

func treeSearch(
	position int,
	omega int,
	currentLog, boundLog float64,
	indexes, allValues []int,
	logs []float64,
	exponents []int,
	modState *ModState,
) {
	if position == omega {

		if !modState.isInvalid() && validExponentSet(indexes, exponents, allValues) {
			fmt.Println("--------------------------------------------")
			fmt.Printf("%v\n", indexesToValues(indexes, allValues))
		}
		return
	}

	optSieveBound := optSieveBoundLog(omega, indexes, allValues, boundLog)

	if currentLog > optSieveBound {
		return
	}

	index := indexes[position]
	//p := allValues[index]
	logp := logs[index]

	e := 1
	logAcc := currentLog + logp
	for logAcc <= optSieveBound {
		exponents[position] = e

		newMod := &ModState{mods: append([]int{}, modState.mods...)}
		newMod.pushPrimeExp(allValues[index], e)

		//bail early if we know this path can't yield a valid tuple
		if newMod.isInvalid() {
			e++
			logAcc += logp
			continue
		}

		treeSearch(
			position+1,
			omega,
			logAcc,
			boundLog,
			indexes,
			allValues,
			logs,
			exponents,
			newMod,
		)
		e++
		logAcc += logp
	}
}

func canComplete(
	boundLog float64,
	currentLog float64,
	nextIndex int,
	remaining int,
	logs []float64,
) bool {
	// Not enough primes left to complete the tuple
	if nextIndex+remaining > len(logs) {
		return false
	}

	// sum logs of the next `remaining` smallest primes starting at nextIndex
	var future float64
	for i := 0; i < remaining; i++ {
		future += logs[nextIndex+i]
	}
	return currentLog+future <= boundLog
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
	currentLog float64,
) Status {

	if currentDepth == maxDepth {
		exponents := make([]int, maxDepth)
		modState := newModState()
		treeSearch(
			0,
			maxDepth,
			0,
			boundLog,
			indexes,
			primeList,
			logs,
			exponents,
			modState,
		)
		return Continue
	}
	startIndex := 0
	if currentDepth > 0 {
		startIndex = indexes[currentDepth-1] + 1
	}

	limit := maxIndex - (maxDepth - currentDepth) + 1
	for i := startIndex; i < limit; i++ {

		if currentDepth == 0 {
			doneTopLevel++
			percent := 100 * doneTopLevel / totalTopLevel
			if percent >= nextPercent {
				fmt.Printf("Progress: %v%%\n", percent)
				nextPercent += 1
			}

		}

		indexes[currentDepth] = i

		newLog := currentLog + logs[i]

		remainingDepth := maxDepth - (currentDepth + 1)
		nextIndex := i + 1

		if !canComplete(boundLog, newLog, nextIndex, remainingDepth, logs) {
			continue
		}

		status := recursiveLoop(
			currentDepth+1,
			maxDepth,
			maxIndex,
			boundLog,
			indexes,
			primeList,
			logs,
			newLog,
		)

		switch status {
		case Continue, Backtrack:
			continue
		case Stop:
			return Stop
		}
	}
	return Backtrack
}
