package main

import (
	"testing"
)

func Test_getFileType(t *testing.T) {
	dir_list := []string{
		".",
		"..",
		"/",
		"/dev"}
	for _, f := range dir_list {
		ret, err := getFileType(f)
		if ret != DIR || err != nil {
			t.Error(f)
		}
	}

	invalid_list := []string{
		"",
		"516e7cb4-6ecf-11d6-8ff8-00022d09712b"}
	for _, f := range invalid_list {
		ret, _ := getFileType(f)
		if ret != INVALID {
			t.Error(f)
		}
	}
}

func Test_pathExists(t *testing.T) {
	dir_list := []string{
		".",
		"..",
		"/",
		"/dev"}
	for _, f := range dir_list {
		exists, err := pathExists(f)
		if !exists || err != nil {
			t.Error(f)
		}
	}

	invalid_list := []string{
		"",
		"516e7cb4-6ecf-11d6-8ff8-00022d09712b"}
	for _, f := range invalid_list {
		exists, _ := pathExists(f)
		if exists {
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
		_, valid := isValidHexSum(s)
		if !valid {
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
		_, valid := isValidHexSum(s)
		if valid {
			t.Error(s)
		}
	}
}
