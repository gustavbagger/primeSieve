package recursion

import (
	"bufio"
	"fmt"
	"time"
)

type Config struct {
	w     *bufio.Writer
	buf   []byte
	Start time.Time
	Count uint64
	omega int
	s     int
}

func NewConfig(
	w *bufio.Writer,
	buf []byte,
	omega int,
	s int,
) Config {
	return Config{buf: buf, w: w, Start: time.Now(), Count: 0, omega: omega, s: s}
}

func (cfg *Config) handleSuccess(indexes, exponents []int) {
	cfg.Count++
	if cfg.Count%100000 == 0 {
		var pivot int
		for i := 0; i < cfg.omega; i++ {
			if i != indexes[i] {
				pivot = i
				break
			}
		}
		fmt.Printf("vals: %.2e, pivot: %v, time: %v.\n", float64(cfg.Count), pivot, time.Since(cfg.Start).Round(time.Second))
	}
	cfg.WriteToBin(indexes, exponents)
}
