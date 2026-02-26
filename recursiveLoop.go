package main

func recursiveLoop(
	currentDepth, maxDepth, maxIndex int,
	boundLog float64,
	indexes, primeList []int,
	logs []float64,
	currentLog float64,
	exponents []int,
) Status {

	if currentDepth == maxDepth {

		optSieveBound := optSieveBoundLog(maxDepth, indexes, primeList, boundLog)

		treeSearch(
			0,
			maxDepth,
			currentLog,
			optSieveBound,
			indexes,
			primeList,
			logs,
			exponents,
		)
		return Continue
	}
	startIndex := 0
	if currentDepth > 0 {
		startIndex = indexes[currentDepth-1] + 1
	}

	limit := maxIndex - (maxDepth - currentDepth) + 1
	for i := startIndex; i < limit; i++ {

		indexes[currentDepth] = i

		newLog := currentLog + logs[i]

		remainingDepth := maxDepth - (currentDepth + 1)
		nextIndex := i + 1

		if !canComplete(boundLog, newLog, nextIndex, remainingDepth, logs) {
			break
		}

		status := recursiveLoop(
			currentDepth+1,
			maxDepth,
			maxIndex,
			boundLog,
			indexes,
			primeList,
			logs,
			newLog,
			exponents,
		)

		switch status {
		case Continue, Backtrack:
			continue
		case Stop:
			return Stop
		}
	}
	return Backtrack
}
