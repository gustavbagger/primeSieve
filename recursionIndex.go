package main

func (cfg *Config) recursionIndex(
	currentDepth, maxIndex int,
	boundLog float64,
	indexes, primeList []int,
	logs []float64,
	currentLog float64,
	exponents []int,
) Status {

	if currentDepth == cfg.omega {

		optSieveBound := optSieveBoundLog(cfg.omega, indexes, primeList, boundLog)

		cfg.recursionExponent(
			0,
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

	limit := maxIndex - (cfg.omega - currentDepth) + 1
	for i := startIndex; i < limit; i++ {

		indexes[currentDepth] = i

		newLog := currentLog + logs[i]

		remainingDepth := cfg.omega - (currentDepth + 1)
		nextIndex := i + 1

		if !canComplete(boundLog, newLog, nextIndex, remainingDepth, logs) {
			break
		}

		status := cfg.recursionIndex(
			currentDepth+1,
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
