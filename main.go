package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {

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
		search(a, b, c, "data.bin")
	case 5:
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
		search(a, b, c, "./data/"+os.Args[4]+".bin")
	default:
		fmt.Println("requires 1,2 or 3 args. Example: \n",
			"'primeSieve <n> <a> <b>'\n", "computes omega = n with p-1< a*10^b")
	}
}

/*
func main() {

	a := 33
	b := 28
	c := 54
	fmt.Printf("searching for omega = %v below %v*10^%v:\n", a, b, c)
	search(a, b, c, fmt.Sprintf("./data/o%v.bin", a))
}

*/
