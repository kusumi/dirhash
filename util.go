package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type fileType int

const (
	DIR fileType = iota
	REG
	DEVICE
	SYMLINK
	UNSUPPORTED
	INVALID

	DIR_STR         = "directory"
	REG_STR         = "regular file"
	DEVICE_STR      = "device"
	SYMLINK_STR     = "symlink"
	UNSUPPORTED_STR = "unsupported file"
	INVALID_STR     = "invalid file"
)

func canonicalizePath(l string) (string, error) {
	if s, err := filepath.EvalSymlinks(l); err != nil {
		if info, err := os.Lstat(l); err != nil {
			return "", err
		} else if info.Mode()&fs.ModeSymlink != 0 {
			return "", nil // ignore broken symlink
		} else {
			return "", err
		}
	} else {
		return s, nil
	}
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func getPathSeparator() rune {
	return os.PathSeparator
}

func getRawFileType(f string) (fileType, error) {
	if info, err := os.Lstat(f); err != nil {
		return INVALID, err
	} else {
		return getModeType(info.Mode())
	}
}

func getFileType(f string) (fileType, error) {
	if info, err := os.Stat(f); err != nil {
		return INVALID, err
	} else {
		return getModeType(info.Mode())
	}
}

func getFileTypeString(t fileType) string {
	switch t {
	case DIR:
		return DIR_STR
	case REG:
		return REG_STR
	case DEVICE:
		return DEVICE_STR
	case SYMLINK:
		return SYMLINK_STR
	case UNSUPPORTED:
		return UNSUPPORTED_STR
	case INVALID:
		return INVALID_STR
	default:
		panicFileType("", "unknown", t)
		return "" // not reached
	}
}

func getModeType(m fs.FileMode) (fileType, error) {
	if m.IsDir() {
		return DIR, nil
	} else if m.IsRegular() {
		return REG, nil
	} else if m&fs.ModeDevice != 0 {
		// XXX assuming blk on Linux, chr on *BSD
		return DEVICE, nil
	} else if m&fs.ModeSymlink != 0 {
		return SYMLINK, nil
	} else {
		return UNSUPPORTED, nil
	}
}

func pathExists(f string) (bool, error) {
	if _, err := os.Stat(f); err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func isValidHexSum(s string) (string, bool) {
	orig := s
	s = strings.TrimPrefix(s, "0x")

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

func getXsumFormatString(f string, h string) string {
	if optSwap {
		h, f = f, h
	}
	// compatible with shaXsum commands
	return fmt.Sprintf("%s  %s", h, f)
}

func getNumFormatString(n uint, msg string) string {
	if len(msg) == 0 {
		return "???"
	}

	s := fmt.Sprintf("%d %s", n, msg)
	if n > 1 {
		if msg == DIR_STR {
			s = s[:len(s)-1]
			s += "ies"
			assert(strings.HasSuffix(s, "directories"))
		} else {
			s += "s"
		}
	}
	return s
}

func printNumFormatString(n uint, msg string) {
	fmt.Println(getNumFormatString(n, msg))
}

func assert(c bool) {
	kassert(c, "Assert failed")
}

func kassert(c bool, err interface{}) {
	if !c {
		panic(err)
	}
}

func panicFileType(f string, how string, t fileType) {
	if len(f) != 0 {
		panic(fmt.Sprintf("%s has %s file type %d", f, how, t))
	} else {
		panic(fmt.Sprintf("%s file type %d", how, t))
	}
}
