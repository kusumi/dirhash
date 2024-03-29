package main

import (
	"fmt"
	"testing"
)

func Test_canonicalizePath(t *testing.T) {
	path_list := []struct {
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
	for _, x := range path_list {
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
	dir_list = []string{
		".",
		"..",
		"/",
		"/dev"}
	invalid_list = []string{
		"",
		"516e7cb4-6ecf-11d6-8ff8-00022d09712b"}
)

func Test_getRawFileType(t *testing.T) {
	for _, f := range dir_list {
		if ret, err := getRawFileType(f); ret != DIR || err != nil {
			t.Error(f)
		}
	}
	for _, f := range invalid_list {
		if ret, _ := getRawFileType(f); ret != INVALID {
			t.Error(f)
		}
	}
}

func Test_getFileType(t *testing.T) {
	for _, f := range dir_list {
		if ret, err := getFileType(f); ret != DIR || err != nil {
			t.Error(f)
		}
	}
	for _, f := range invalid_list {
		if ret, _ := getFileType(f); ret != INVALID {
			t.Error(f)
		}
	}
}

func Test_getFileTypeString(t *testing.T) {
	file_type_list := []struct {
		typ fileType
		str string
	}{
		{DIR, "directory"},
		{REG, "regular file"},
		{DEVICE, "device"},
		{SYMLINK, "symlink"},
		{UNSUPPORTED, "unsupported file"},
		{INVALID, "invalid file"},
	}
	for _, x := range file_type_list {
		if getFileTypeString(x.typ) != x.str {
			t.Error(x)
		}
	}
}

func Test_pathExists(t *testing.T) {
	for _, f := range dir_list {
		if exists, err := pathExists(f); !exists || err != nil {
			t.Error(f)
		}
	}
	for _, f := range invalid_list {
		if exists, err := pathExists(f); exists || err == nil {
			t.Error(f)
		}
	}
}

func Test_isValidHexSum(t *testing.T) {
	valid_list := []string{
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
	for _, s := range valid_list {
		if _, valid := isValidHexSum(s); !valid {
			t.Error(s)
		}
	}

	invalid_list := []string{
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
	for _, s := range invalid_list {
		if _, valid := isValidHexSum(s); valid {
			t.Error(s)
		}
	}
}

func Test_getNumFormatString(t *testing.T) {
	num_format_list := []struct {
		n      uint64
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
	for _, x := range num_format_list {
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
