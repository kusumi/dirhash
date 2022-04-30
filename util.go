package main

import (
	"io/fs"
	"os"
	"strings"
)

type fileType int

const (
	DIR fileType = iota
	REG
	UNSUPPORTED
	INVALID
)

func getFileType(f string) (fileType, error) {
	info, err := os.Stat(f)
	if err != nil {
		return INVALID, err
	}

	return getModeType(info.Mode())
}

func getModeType(m fs.FileMode) (fileType, error) {
	if m.IsDir() {
		return DIR, nil
	} else if m.IsRegular() {
		return REG, nil
	}

	return UNSUPPORTED, nil
}

func pathExists(f string) (bool, error) {
	_, err := os.Stat(f)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func isValidHexSum(s string) (string, bool) {
	orig := s
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}

	if len([]rune(s)) < 32 {
		return orig, false
	}

	for _, r := range s {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
			return orig, false
		}
	}
	return s, true
}

func assert(c bool) {
	if !c {
		panic("Assertion")
	}
}
