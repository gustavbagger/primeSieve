package main

import (
	"fmt"

	pr "github.com/fxtlabs/primes"
)

func main() {

	primeList := pr.Sieve(2000)

	fmt.Printf("%5.2e\n", Prod(append(primeList[:33], primeList[218])))
	pritnIntervals()

}
