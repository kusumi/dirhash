package main

import (
	"testing"
)

func Test_newHash(t *testing.T) {
	for _, s := range getAvailableHashAlgo() {
		h := newHash(s)
		if h == nil {
			t.Error(s)
		}
	}

	invalid_list := []string{
		"",
		"xxx",
		"SHA256",
		"516e7cb4-6ecf-11d6-8ff8-00022d09712b"}
	for _, s := range invalid_list {
		h := newHash(s)
		if h != nil {
			t.Error(s)
		}
	}
}

var (
	alg_sum_list = []struct {
		hash_algo string // name
		hex_sum   string // value for empty input
	}{
		{MD5, "d41d8cd98f00b204e9800998ecf8427e"},
		{SHA1, "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{SHA224, "d14a028c2a3a2bc9476102bb288234c415a2b01f828ea62ac5b3e42f"},
		{SHA256, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{SHA384, "38b060a751ac96384cd9327eb1b1e36a21fdb71114be07434c0cc7bf63f6e1da274edebfe76f65fbd51ad2f14898b95b"},
		{SHA512, "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
	}
)

func Test_getByteHash(t *testing.T) {
	for name, x := range alg_sum_list {
		sum, err := getByteHash([]byte{}, x.hash_algo)
		if err != nil {
			t.Error(name, x)
		}
		hex_sum := getHexSum(sum)
		if hex_sum != x.hex_sum {
			t.Error(name, x)
		}
	}
}

func Test_getStringHash(t *testing.T) {
	for name, x := range alg_sum_list {
		sum, err := getStringHash("", x.hash_algo)
		if err != nil {
			t.Error(name, x)
		}
		hex_sum := getHexSum(sum)
		if hex_sum != x.hex_sum {
			t.Error(name, x)
		}
	}
}
