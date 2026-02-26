package main

import (
	"fmt"
)

func treeSearch(
	position int,
	omega int,
	currentLog, optSieveBound float64,
	indexes, allValues []int,
	logs []float64,
	exponents []int,
) {
	if currentLog > optSieveBound {
		return
	}
	n, valid := validExponentSet(indexes, exponents, allValues)
	if valid {
		count++
		fmt.Println("--------------------------------------------")
		fmt.Printf("%v\n", n)
		fmt.Printf("%v\n", indexes)
		fmt.Printf("%v\n", exponents)
	}

	index := indexes[position]
	logp := logs[index]

	remainingLog := optSieveBound - currentLog

	maxE := int(remainingLog / logp)
	if maxE < 1 {
		return
	}

	for e := 1; e <= maxE; e++ {
		exponents[position] = e

		newLog := currentLog + float64(e)*logp

		treeSearch(
			position+1,
			omega,
			newLog,
			optSieveBound,
			indexes,
			allValues,
			logs,
			exponents,
		)
	}
}
