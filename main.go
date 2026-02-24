package main

import (
	"fmt"

	pr "github.com/fxtlabs/primes"
)

func main() {
	omega := 34
	upperBound := 1 << 62
	primeList := pr.Sieve(2000)

	maxIndex := len(primeList) - 1
	indexes := make([]int, omega)
	prod := Prod([]int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 53, 53, 67, 71, 73, 83, 89, 89, 97, 101, 103, 107, 107, 107, 109, 109, 137, 137, 139, 139, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 149, 151, 151, 151, 157, 157, 157, 163, 163, 173, 179})
	fmt.Println(prod + 1)
	fmt.Println(pr.IsPrime(prod + 1))
	recursiveLoop(0, omega, maxIndex, upperBound, indexes, primeList)
	fmt.Println("--------------------------------------------")
}
