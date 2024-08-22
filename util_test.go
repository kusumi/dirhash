package main

import (
	"fmt"
	"testing"
)

func Test_canonicalizePath(t *testing.T) {
	pathList := []struct {
		i string
		o string
	}{
		{"/", "/"},
		{"/////", "/"},
		{"/..", "/"},
		{"/../", "/"},
		{"/root", "/root"},
		{"/root/", "/root"},
		{"/root/..", "/"},
		{"/root/../dev", "/dev"},
	}
	for _, x := range pathList {
		if s, err := canonicalizePath(x.i); err != nil || s != x.o {
			t.Error(x)
		}
	}
}

func Test_isWindows(t *testing.T) {
	if isWindows() {
		t.Error("Windows unsupported")
	}
}

func Test_getPathSeparator(t *testing.T) {
	if isWindows() {
		return
	}
	if s := getPathSeparator(); s != '/' {
		t.Error(s)
	}
}

var (
	dirList = []string{
		".",
		"..",
		"/",
		"/dev"}
	invalidList = []string{
		"",
		"516e7cb4-6ecf-11d6-8ff8-00022d09712b"}
)

func Test_getRawFileType(t *testing.T) {
	for _, f := range dirList {
		if ret, err := getRawFileType(f); ret != typeDir || err != nil {
			t.Error(f)
		}
	}
	for _, f := range invalidList {
		if ret, _ := getRawFileType(f); ret != typeInvalid {
			t.Error(f)
		}
	}
}

func Test_getFileType(t *testing.T) {
	for _, f := range dirList {
		if ret, err := getFileType(f); ret != typeDir || err != nil {
			t.Error(f)
		}
	}
	for _, f := range invalidList {
		if ret, _ := getFileType(f); ret != typeInvalid {
			t.Error(f)
		}
	}
}

func Test_getFileTypeString(t *testing.T) {
	fileTypeList := []struct {
		typ fileType
		str string
	}{
		{typeDir, "directory"},
		{typeReg, "regular file"},
		{typeDevice, "device"},
		{typeSymlink, "symlink"},
		{typeUnsupported, "unsupported file"},
		{typeInvalid, "invalid file"},
	}
	for _, x := range fileTypeList {
		if getFileTypeString(x.typ) != x.str {
			t.Error(x)
		}
	}
}

func Test_pathExists(t *testing.T) {
	for _, f := range dirList {
		if exists, err := pathExists(f); !exists || err != nil {
			t.Error(f)
		}
	}
	for _, f := range invalidList {
		if exists, err := pathExists(f); exists || err == nil {
			t.Error(f)
		}
	}
}

func Test_isValidHexSum(t *testing.T) {
	validList := []string{
		"00000000000000000000000000000000",
		"11111111111111111111111111111111",
		"22222222222222222222222222222222",
		"33333333333333333333333333333333",
		"44444444444444444444444444444444",
		"55555555555555555555555555555555",
		"66666666666666666666666666666666",
		"77777777777777777777777777777777",
		"88888888888888888888888888888888",
		"99999999999999999999999999999999",
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB",
		"CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC",
		"DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD",
		"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccc",
		"dddddddddddddddddddddddddddddddd",
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		"ffffffffffffffffffffffffffffffff",
		"0123456789ABCDEFabcdef0123456789ABCDEFabcdef",
		"0x00000000000000000000000000000000",
		"0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"0x0123456789ABCDEFabcdef0123456789ABCDEFabcdef"}
	for _, s := range validList {
		if _, valid := isValidHexSum(s); !valid {
			t.Error(s)
		}
	}

	invalidList := []string{
		"gggggggggggggggggggggggggggggggg",
		"GGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG",
		"zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
		"ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ",
		"                                ",
		"################################",
		"--------------------------------",
		"................................",
		"@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@",
		"________________________________",
		"0000000000000000000000000000000",
		"0x0000000000000000000000000000000",
		"0x",
		"0",
		""}
	for _, s := range invalidList {
		if _, valid := isValidHexSum(s); valid {
			t.Error(s)
		}
	}
}

func Test_getNumFormatString(t *testing.T) {
	numFormatList := []struct {
		n      uint
		msg    string
		result string
	}{
		{0, "", "???"},
		{1, "", "???"},
		{2, "", "???"},
		{0, "file", "0 file"},
		{1, "file", "1 file"},
		{2, "file", "2 files"},
	}
	for _, x := range numFormatList {
		if getNumFormatString(x.n, x.msg) != x.result {
			t.Error(x)
		}
	}
}

func Test_assert(t *testing.T) {
	assert(true)
	assert(!false)
	assert(true != false)

	assert(0 == 0+0)
	assert(1 == 1+0)
	assert(0 != 1+0)

	assert("" == ""+"")
	assert("xxx" == "xxx"+"")
	assert("xxx" != "yyy")
}

func Test_kassert(t *testing.T) {
	kassert(true, nil)
	kassert(!false, nil)

	kassert(true, "")
	kassert(!false, "")

	kassert(true, fmt.Errorf(""))
	kassert(!false, fmt.Errorf(""))
}
