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
		fmt.Printf("vals: %.2e, time: %v.\n", float64(cfg.Count), time.Since(cfg.Start).Round(time.Second))
	}
	cfg.WriteToBin(indexes, exponents)
}
