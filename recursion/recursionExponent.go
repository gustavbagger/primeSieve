package recursion

import (
	"fmt"
	"time"

	"github.com/gustavbagger/primeSieve/primality"
)

func (cfg *Config) handleSuccess(indexes, exponents []int) {
	cfg.Count++
	if cfg.Count%1000000 == 0 {
		fmt.Printf("vals: %.2e, time: %v.\n", float64(cfg.Count), time.Since(cfg.Start).Round(time.Second))
	}
	cfg.WriteToBin(indexes, exponents)
}

func (cfg *Config) recursionExponent(
	position int,
	currentLog, optSieveBound float64,
	indexes, allValues []int,
	logs []float64,
	exponents []int,
) {
	if currentLog > optSieveBound {
		return
	}
	_, valid := primality.ValidExponentSet192(indexes, exponents, allValues)
	if valid {
		cfg.handleSuccess(indexes, exponents)

		cfg.WriteToBin(indexes, exponents)
	}

	index := indexes[position]
	logp := logs[index]

	remainingLog := optSieveBound - currentLog

	maxE := int(remainingLog / logp)
	if maxE < 1 {
		return
	}

	for e := 1; e <= maxE; e++ {
		exponents[position] = e

		newLog := currentLog + float64(e)*logp

		cfg.recursionExponent(
			position+1,
			newLog,
			optSieveBound,
			indexes,
			allValues,
			logs,
			exponents,
		)
	}
}
