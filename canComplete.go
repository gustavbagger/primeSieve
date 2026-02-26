package main

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
