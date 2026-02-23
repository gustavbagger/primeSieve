package main

func deltaSum(list []int) float64 {
	var sum float64
	for _, p := range list {
		sum -= 1.0 / float64(p)
	}
	return sum + 1
}

func Prod(list []int) float64 {
	prod := 1.0
	for _, p := range list {
		prod *= float64(p)
	}
	return prod
}
