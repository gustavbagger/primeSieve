package main

import (
	"fmt"
	"math"

	pr "github.com/fxtlabs/primes"
)

func main() {

	omega := 34
	// bound is a * 10^b
	a := 1
	b := 55
	boundLog := math.Log(float64(a) * math.Pow10(b))
	primeList := pr.Sieve(2000)
	logs := make([]float64, len(primeList))

	for i, p := range primeList {
		logs[i] = math.Log(float64(p))
	}

	maxIndex := len(primeList) - 1
	indexes := make([]int, omega)

	recursiveLoop(
		0,
		omega,
		maxIndex,
		boundLog,
		indexes,
		primeList,
		logs,
	)
	//fmt.Println(logs)
	fmt.Println("--------------------------------------------")
}
