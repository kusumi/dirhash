package main

import (
	"testing"
)

func Test_numStatRegular(t *testing.T) {
	// 0
	initStat()
	if x := numStatRegular(); x != 0 {
		t.Error(x)
	}

	// 0
	initStat()
	if x := numStatRegular(); x != 0 {
		t.Error(x)
	}
}

func Test_appendStatRegular(t *testing.T) {
	// 1
	initStat()
	appendStatRegular("a")
	if x := numStatRegular(); x != 1 {
		t.Error(x)
	}
	if f := statRegular[0]; f != "a" {
		t.Error(f)
	}

	// 2
	appendStatRegular("b")
	if x := numStatRegular(); x != 2 {
		t.Error(x)
	}
	if f := statRegular[0]; f != "a" {
		t.Error(f)
	}
	if f := statRegular[1]; f != "b" {
		t.Error(f)
	}

	// 3
	appendStatRegular("c")
	if x := numStatRegular(); x != 3 {
		t.Error(x)
	}
	if f := statRegular[0]; f != "a" {
		t.Error(f)
	}
	if f := statRegular[1]; f != "b" {
		t.Error(f)
	}
	if f := statRegular[2]; f != "c" {
		t.Error(f)
	}

	// 1
	initStat()
	appendStatRegular("d")
	if x := numStatRegular(); x != 1 {
		t.Error(x)
	}
	if f := statRegular[0]; f != "d" {
		t.Error(f)
	}
}

func Test_numWrittenRegular(t *testing.T) {
	initStat()
	if x := numWrittenRegular(); x != 0 {
		t.Error(x)
	}
}

func Test_appendWrittenRegular(t *testing.T) {
	initStat()
	appendWrittenRegular(9999999999)
	if x := numWrittenRegular(); x != 9999999999 {
		t.Error(x)
	}

	appendWrittenRegular(1)
	if x := numWrittenRegular(); x != 10000000000 {
		t.Error(x)
	}

	initStat()
	if x := numWrittenRegular(); x != 0 {
		t.Error(x)
	}
}
