package main

import (
	"bytes"
	"strings"
	"testing"
)

func Test_newHash(t *testing.T) {
	for _, s := range getAvailableHashAlgo() {
		if h := newHash(s); h == nil {
			t.Error(s)
		}
	}

	invalidList := []string{
		"",
		"xxx",
		"SHA256",
		"516e7cb4-6ecf-11d6-8ff8-00022d09712b"}
	for _, s := range invalidList {
		if h := newHash(s); h != nil {
			t.Error(s)
		}
	}
}

var (
	algSumList1 = []struct {
		hashAlgo string
		hexSum   string // value for empty input
	}{
		{MD5, "d41d8cd98f00b204e9800998ecf8427e"},
		{SHA1, "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{SHA224, "d14a028c2a3a2bc9476102bb288234c415a2b01f828ea62ac5b3e42f"},
		{SHA256, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{SHA384, "38b060a751ac96384cd9327eb1b1e36a21fdb71114be07434c0cc7bf63f6e1da274edebfe76f65fbd51ad2f14898b95b"},
		{SHA512, "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
	}

	algSumList2 = []struct {
		hashAlgo string
		hexSum   string
	}{
		{MD5, "48fcdb8b87ce8ef779774199a856091d"},
		{SHA1, "065e431442d313aa4c4345f1c7f3d3a84a9b201f"},
		{SHA224, "62f2929306a761f06a3b055aac36ec38df8e275a8b66e68c52f030d3"},
		{SHA256, "e23c0cda5bcdecddec446b54439995c7260c8cdcf2953eec9f5cdb6948e5898d"},
		{SHA384, "3a52aaed14b5b6f9f7208914e5c34f0e16e70a285c37fd964ab918980a40acb52be0a71d43cdabb702aa2d025ce9ab7b"},
		{SHA512, "990fed5cd10a549977ef6c9e58019a467f6c7aadffb9a6d22b2d060e6989a06d5beb473ebc217f3d553e16bf482efdc4dd91870e7943723fdc387c2e9fa3a4b8"},
	}
	algSumList2Str    = "A"
	algSumList2Repeat = 1000000
)

func Test_getByteHash(t *testing.T) {
	for _, x := range algSumList1 {
		written, sum, err := getByteHash([]byte{}, x.hashAlgo)
		if err != nil {
			t.Error(x)
		}
		if written != 0 {
			t.Error(written)
		}
		if hexSum := getHexSum(sum); hexSum != x.hexSum {
			t.Error(x)
		}
	}

	s := bytes.Repeat([]byte(algSumList2Str), algSumList2Repeat)
	for _, x := range algSumList2 {
		written, sum, err := getByteHash(s, x.hashAlgo)
		if err != nil {
			t.Error(x)
		}
		if written != uint64(algSumList2Repeat) {
			t.Error(written)
		}
		if hexSum := getHexSum(sum); hexSum != x.hexSum {
			t.Error(x)
		}
	}
}

func Test_getStringHash(t *testing.T) {
	for _, x := range algSumList1 {
		written, sum, err := getStringHash("", x.hashAlgo)
		if err != nil {
			t.Error(x)
		}
		if written != 0 {
			t.Error(written)
		}
		if hexSum := getHexSum(sum); hexSum != x.hexSum {
			t.Error(x)
		}
	}

	s := strings.Repeat(algSumList2Str, algSumList2Repeat)
	for _, x := range algSumList2 {
		written, sum, err := getStringHash(s, x.hashAlgo)
		if err != nil {
			t.Error(x)
		}
		if written != uint64(algSumList2Repeat) {
			t.Error(written)
		}
		if hexSum := getHexSum(sum); hexSum != x.hexSum {
			t.Error(x)
		}
	}
}
