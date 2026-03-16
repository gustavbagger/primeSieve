package main

import (
	"bufio"
	"time"
)

type Config struct {
	w     *bufio.Writer
	buf   []byte
	start time.Time
	count uint64
	omega int
}

// unsigned 192-bit int
type uint192 struct {
	Lo  uint64
	Mid uint64
	Hi  uint64
}
