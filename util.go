package main

import (
	"fmt"
	"io/fs"
	"os"
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

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func getPathSeparator() string {
	return string(os.PathSeparator)
}

func getRawFileType(f string) (fileType, error) {
	info, err := os.Lstat(f)
	if err != nil {
		return INVALID, err
	}

	return getModeType(info.Mode())
}

func getFileType(f string) (fileType, error) {
	info, err := os.Stat(f)
	if err != nil {
		return INVALID, err
	}

	return getModeType(info.Mode())
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
	}

	return UNSUPPORTED, nil
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
	// compatible with shaXsum commands
	return fmt.Sprintf("%s  %s", h, f)
}

func getNumFormatString(n uint64, msg string) string {
	if msg == "" {
		return "???"
	}

	s := fmt.Sprintf("%d %s", n, msg)
	if n > 1 {
		s += "s"
	}

	return s
}

func printNumFormatString(n uint64, msg string) {
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
	var s string
	if f != "" {
		s = fmt.Sprintf("%s has %s file type %d", f, how, t)
	} else {
		s = fmt.Sprintf("%s file type %d", how, t)
	}
	panic(s)
}
