//go:build squash1

package main

import (
	"sort"
	"strings"
)

var (
	squashBuffer []string
)

func initSquashBuffer() {
	squashBuffer = make([]string, 0)
}

func updateSquashBuffer(s string) {
	// get hash to minimize total string size
	b, err := getStringHash(s, SHA1)
	if err != nil {
		panic(err)
	}
	squashBuffer = append(squashBuffer, getHexSum(b))
}

func getSquashBuffer() string {
	sort.Strings(squashBuffer)
	return strings.Join(squashBuffer, ",")
}
