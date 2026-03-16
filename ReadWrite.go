package main

import (
	"io"
	"os"
)

func (cfg *Config) WriteToBin(indexes, exponents []int) {
	// Write indexes
	for i := 0; i < cfg.omega; i++ {
		v := uint16(indexes[i])
		cfg.buf[2*i] = byte(v)
		cfg.buf[2*i+1] = byte(v >> 8)
	}
	cfg.w.Write(cfg.buf)

	// Write exponents
	for i := 0; i < cfg.omega; i++ {
		v := uint16(exponents[i])
		cfg.buf[2*i] = byte(v)
		cfg.buf[2*i+1] = byte(v >> 8)
	}
	cfg.w.Write(cfg.buf)
}

func (cfg *Config) ReadRange(path string, a, b int) ([][]uint16, error) {
	sliceLen := cfg.omega
	recordSize := sliceLen * 2 * 2 // indexes + exponents = 4 * omega bytes

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Compute byte offsets
	start := int64(a * recordSize)
	end := int64(b * recordSize)
	length := end - start

	// Seek to the start of the range
	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Read the entire block
	buf := make([]byte, length)
	_, err = io.ReadFull(file, buf)
	if err != nil {
		return nil, err
	}

	// Decode into slices
	out := make([][]uint16, b-a)
	pos := 0

	for i := 0; i < b-a; i++ {
		entry := make([]uint16, sliceLen*2) // indexes + exponents

		for j := 0; j < sliceLen*2; j++ {
			lo := uint16(buf[pos])
			hi := uint16(buf[pos+1])
			entry[j] = lo | (hi << 8)
			pos += 2
		}

		out[i] = entry
	}

	return out, nil
}
