package main

import (
	"io/fs"
	"os"
	"strings"
)

const (
	DIR = iota
	REG
	SYMLINK
	UNSUPPORTED
	INVALID
)

func getRawFileType(f string) (int, error) {
	info, err := os.Lstat(f)
	if err != nil {
		return INVALID, err
	}

	return getModeType(info.Mode())
}

func getFileType(f string) (int, error) {
	info, err := os.Stat(f)
	if err != nil {
		return INVALID, err
	}

	return getModeType(info.Mode())
}

func getModeType(m fs.FileMode) (int, error) {
	if m.IsDir() {
		return DIR, nil
	} else if m.IsRegular() {
		return REG, nil
	} else if m&fs.ModeSymlink != 0 {
		return SYMLINK, nil
	}

	return UNSUPPORTED, nil
}

func pathExists(f string) (bool, error) {
	_, err := os.Stat(f)
	if err == nil {
		return true, nil
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
	kassert(c, "Assert failed")
}

func kassert(c bool, err interface{}) {
	if !c {
		panic(err)
	}
}
