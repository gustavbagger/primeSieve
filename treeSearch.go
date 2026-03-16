package main

import (
	"fmt"
	"time"
)

func treeSearch(
	position int,
	omega int,
	currentLog, optSieveBound float64,
	indexes, allValues []int,
	logs []float64,
	exponents []int,
	buffer buffer,
) {
	if currentLog > optSieveBound {
		return
	}
	_, valid := validExponentSet192(indexes, exponents, allValues)
	if valid {
		count++
		if count%100000 == 0 {
			fmt.Printf("%.2e values found - expect 10^8 (for o=33).\n", float64(count))
			fmt.Println("Total time: ", time.Now().Sub(buffer.start))
		}

		WriteToBin(indexes, exponents, buffer.w, buffer.buf, omega)
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
			buffer,
		)
	}
}
