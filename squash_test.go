package main

import (
	"strings"
	"testing"
)

func Test_initSquashBuffer(t *testing.T) {
	initSquashBuffer()

	if b := getSquashBuffer(); b == nil {
		t.Error(b)
	}
}

func Test_updateSquashBuffer(t *testing.T) {
	initSquashBuffer()

	updateSquashBuffer([]byte(""))
	if b := getSquashBuffer(); b == nil {
		t.Error(b)
	}

	updateSquashBuffer([]byte(""))
	if b := getSquashBuffer(); b == nil {
		t.Error(b)
	}

	updateSquashBuffer([]byte("xxx"))
	if b := getSquashBuffer(); b == nil {
		t.Error(b)
	}

	updateSquashBuffer([]byte(strings.Repeat("x", 123456)))
	if b := getSquashBuffer(); b == nil {
		t.Error(b)
	}
}
