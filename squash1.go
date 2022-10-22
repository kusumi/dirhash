//go:build squash1

package main

import (
	"sort"
	"strings"
)

var (
	squashBuffer [][]byte
)

func initSquashBuffer() {
	squashBuffer = make([][]byte, 0)
	for i := range squashBuffer {
		squashBuffer[i] = make([]byte, 0)
	}
}

func updateSquashBuffer(b []byte) {
	// get hash to minimize total string size
	_, tmp, err := getByteHash(b, MD5)
	if err != nil {
		panic(err)
	}
	squashBuffer = append(squashBuffer, tmp)
}

func getSquashBuffer() []byte {
	// XXX directly sort [][]byte
	s := make([]string, 0)
	for _, b := range squashBuffer {
		s = append(s, getHexSum(b))
	}

	sort.Strings(s)
	return []byte(strings.Join(s, ""))
}
