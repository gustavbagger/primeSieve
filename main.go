package main

var testVals = []uint192{
	// small primes
	{Lo: 3, Mid: 0, Hi: 0},
	{Lo: 5, Mid: 0, Hi: 0},
	{Lo: 17, Mid: 0, Hi: 0},
	{Lo: 97, Mid: 0, Hi: 0},
	{Lo: 65537, Mid: 0, Hi: 0},

	// small composites
	{Lo: 4, Mid: 0, Hi: 0},   // 2^2
	{Lo: 9, Mid: 0, Hi: 0},   // 3^2
	{Lo: 21, Mid: 0, Hi: 0},  // 3 * 7
	{Lo: 221, Mid: 0, Hi: 0}, // 13 * 17
	{Lo: 341, Mid: 0, Hi: 0}, // 11 * 31 (also a Fermat pseudoprime to base 2)
}

func main() {
	for _, val := range testVals {
		testPRP(val)
	}
	testMulRedc(uint192{Lo: 9})
	testMulRedc(uint192{Lo: 21})
	testMulRedc(uint192{Lo: 221})
}

/*
import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	switch len(os.Args) {
	case 2:
		a, err := strconv.Atoi(os.Args[1])
		if err != nil {
			return
		}
		fmt.Println("omega intervals:")
		printIntervals(a, a)
	case 3:
		a, err := strconv.Atoi(os.Args[1])
		if err != nil {
			return
		}
		b, err := strconv.Atoi(os.Args[2])
		if err != nil {
			return
		}
		fmt.Println("omega intervals:")
		printIntervals(a, b)
	case 4:
		a, err := strconv.Atoi(os.Args[1])
		if err != nil {
			return
		}
		b, err := strconv.Atoi(os.Args[2])
		if err != nil {
			return
		}
		c, err := strconv.Atoi(os.Args[3])
		if err != nil {
			return
		}
		fmt.Printf("searching for omega = %v below %v*10^%v:\n", a, b, c)
		search(a, b, c)
	default:
		fmt.Println("requires 1,2 or 3 args. Example: \n",
			"'primeSieve <n> <a> <b>'\n", "computes omega = n with p-1< a*10^b")
	}
	end := time.Now()
	fmt.Println("Time elapsed: ", end.Sub(start))
}
*/
