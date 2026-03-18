package recursion

import (
	"bufio"
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
