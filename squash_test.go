package main

import (
	"strings"
	"testing"
)

func Test_initSquashBuffer(t *testing.T) {
	initSquashBuffer()

	s := getSquashBuffer()
	if s != "" {
		t.Error(s)
	}
}

func Test_updateSquashBuffer(t *testing.T) {
	initSquashBuffer()

	updateSquashBuffer("")
	s1 := getSquashBuffer()
	if s1 == "" {
		t.Error(s1)
	}

	updateSquashBuffer("")
	s2 := getSquashBuffer()
	if s2 == "" {
		t.Error(s2)
	}

	updateSquashBuffer("xxx")
	s3 := getSquashBuffer()
	if s3 == "" {
		t.Error(s3)
	}

	updateSquashBuffer(strings.Repeat("x", 123456))
	s4 := getSquashBuffer()
	if s4 == "" {
		t.Error(s4)
	}
}
