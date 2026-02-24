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
		search(a, b, c)
	default:
		fmt.Println("wrong args")
	}
}
