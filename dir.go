package main

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

var (
	inputPrefix string
)

func printInput(f string) error {
	// keep symlink input as is
	if t, err := getRawFileType(f); err != nil {
		return err
	} else if t != SYMLINK {
		if x, err := canonicalizePath(f); err != nil {
			return err
		} else if len(x) == 0 {
			return nil
		} else {
			f = x
		}
		// assert exists
		if _, err := pathExists(f); err != nil {
			return err
		}
	}

	// convert input to abs first
	f, err := filepath.Abs(f)
	if err != nil {
		return err
	}
	assertFilePath(f)

	// keep input prefix based on raw type
	t, err := getRawFileType(f)
	if err != nil {
		return err
	}
	switch t {
	case DIR:
		inputPrefix = f
	case REG:
		fallthrough
	case DEVICE:
		fallthrough
	case SYMLINK:
		inputPrefix = filepath.Dir(f)
	default:
		return fmt.Errorf("%s has unsupported type %d", f, t)
	}

	// prefix is a directory
	t, _ = getFileType(inputPrefix)
	assert(t == DIR)

	// initialize global resource
	initStat()
	initSquashBuffer()

	// start directory walk
	if err := walkDirectory(f); err != nil {
		return err
	}

	// print various stats
	if optVerbose {
		printVerboseStat()
	}
	printStatUnsupported()
	printStatInvalid()

	// print squash hash if specified
	if optSquash {
		b := getSquashBuffer()
		if optVerbose {
			printNumFormatString(uint(len(b)), "squashed byte")
		}
		if err := printByte(f, b); err != nil {
			return err
		}
	}

	return nil
}

func walkDirectory(f string) error {
	var l []string
	if err := filepath.WalkDir(f,
		func(f string, d fs.DirEntry, err error) error {
			assertFilePath(f)
			if err != nil {
				return err
			}
			if optSort {
				l = append(l, f)
				return nil
			} else {
				return walkDirectoryImpl(f)
			}
		}); err != nil {
		return err
	}

	if optSort {
		sort.Strings(l)
		for _, f := range l {
			if err := walkDirectoryImpl(f); err != nil {
				return err
			}
		}
	}

	return nil
}

func walkDirectoryImpl(f string) error {
	t, err := getRawFileType(f)
	if err != nil {
		return err
	}

	if testIgnoreEntry(f, t) {
		appendStatIgnored(f)
		return nil
	}

	// find target if symlink
	var x, l string // l is symlink itself, not its target
	switch t {
	case SYMLINK:
		if optIgnoreSymlink {
			appendStatIgnored(f)
			return nil
		}
		if !optFollowSymlink {
			return printSymlink(f)
		}
		x, err = canonicalizePath(f)
		if err != nil {
			return err
		} else if len(x) == 0 {
			return printInvalid(f)
		}
		assert(filepath.IsAbs(x))
		t, err = getFileType(x) // update type
		if err != nil {
			return err
		}
		assert(t != SYMLINK) // symlink chains resolved
		l = f
	default:
		x = f
		l = ""
	}

	switch t {
	case DIR:
		return handleDirectory(x, l)
	case REG:
		fallthrough
	case DEVICE:
		return printFile(x, l, t)
	case UNSUPPORTED:
		return printUnsupported(x)
	case INVALID:
		return printInvalid(x)
	default:
		panicFileType(x, "unknown", t)
	}

	assert(false)
	return nil
}

func testIgnoreEntry(f string, t fileType) bool {
	assert(filepath.IsAbs(f))

	// only non directory types count
	if t == DIR {
		return false
	}

	base_starts_with_dot := strings.HasPrefix(path.Base(f), ".")
	path_contains_slash_dot := strings.Contains(f, "/.")

	// ignore . directories if specified
	if optIgnoreDotDir {
		if !base_starts_with_dot && path_contains_slash_dot {
			return true
		}
	}

	// ignore . regular files if specified
	if optIgnoreDotFile {
		// XXX limit to REG ?
		if base_starts_with_dot {
			return true
		}
	}

	// ignore . entries if specified
	if optIgnoreDot {
		if base_starts_with_dot || path_contains_slash_dot {
			return true
		}
	}

	return false
}

func trimInputPrefix(f string) string {
	if strings.HasPrefix(f, inputPrefix) {
		f = f[len(inputPrefix)+1:]
		assert(!strings.HasPrefix(f, "/"))
	}
	return f
}

func getRealPath(f string) string {
	if optAbs {
		assert(filepath.IsAbs(f))
		return f
	} else if f == inputPrefix {
		return "."
	} else if inputPrefix == "/" {
		return f[1:]
	} else {
		// f is probably symlink target if f unchanged
		return trimInputPrefix(f)
	}
}

func printByte(f string, inb []byte) error {
	assertFilePath(f)

	// get hash value
	_, b, err := getByteHash(inb, optHashAlgo)
	if err != nil {
		return err
	}
	assert(len(b) > 0)
	hex_sum := getHexSum(b)

	// verify hash value if specified
	if len(optHashVerify) != 0 {
		if optHashVerify != hex_sum {
			return nil
		}
	}

	if optHashOnly {
		fmt.Println(hex_sum)
	} else {
		// no space between two
		s := fmt.Sprintf("[%s][v%d]", squashLabel, squashVersion)
		if realf := getRealPath(f); realf == "." {
			fmt.Println(hex_sum + s)
		} else {
			fmt.Println(getXsumFormatString(realf, hex_sum) + s)
		}
	}

	return nil
}

func handleDirectory(f string, l string) error {
	assertFilePath(f)
	if len(l) > 0 {
		assertFilePath(l)
	}

	// nothing to do if input is input prefix
	if f == inputPrefix {
		return nil
	}

	// nothing to do unless squash
	if !optSquash {
		return nil
	}

	// debug print first
	if optDebug {
		if err := printDebug(f, DIR); err != nil {
			return err
		}
	}

	// get hash value
	// path must be relative to input prefix
	s := trimInputPrefix(f)
	written, b, err := getStringHash(s, optHashAlgo)
	if err != nil {
		return err
	}
	assert(len(b) > 0)

	// count this file
	appendStatTotal()
	appendWrittenTotal(written)
	appendStatDirectory(f)
	appendWrittenDirectory(written)

	// squash
	assert(optSquash)
	if optHashOnly {
		updateSquashBuffer(b)
	} else {
		// make link -> target format if symlink
		realf := getRealPath(f)
		if len(l) > 0 {
			assertFilePath(l)
			if !optAbs {
				l = trimInputPrefix(l)
			}
			realf = fmt.Sprintf("%s -> %s", l, realf)
		}
		updateSquashBuffer(append([]byte(realf), b...))
	}

	return nil
}

func printFile(f string, l string, t fileType) error {
	assertFilePath(f)
	if len(l) > 0 {
		assertFilePath(l)
	}

	// debug print first
	if optDebug {
		if err := printDebug(f, t); err != nil {
			return err
		}
	}

	// get hash value
	written, b, err := getFileHash(f, optHashAlgo)
	if err != nil {
		return err
	}
	assert(len(b) > 0)
	hex_sum := getHexSum(b)

	// count this file
	appendStatTotal()
	appendWrittenTotal(written)
	switch t {
	case REG:
		appendStatRegular(f)
		appendWrittenRegular(written)
	case DEVICE:
		appendStatDevice(f)
		appendWrittenDevice(written)
	default:
		panicFileType(f, "invalid", t)
	}

	// verify hash value if specified
	if len(optHashVerify) != 0 {
		if optHashVerify != hex_sum {
			return nil
		}
	}

	// squash or print this file
	if optHashOnly {
		if optSquash {
			updateSquashBuffer(b)
		} else {
			fmt.Println(hex_sum)
		}
	} else {
		// make link -> target format if symlink
		realf := getRealPath(f)
		if len(l) > 0 {
			assertFilePath(l)
			if !optAbs {
				l = trimInputPrefix(l)
			}
			realf = fmt.Sprintf("%s -> %s", l, realf)
		}
		if optSquash {
			updateSquashBuffer(append([]byte(realf), b...))
		} else {
			fmt.Println(getXsumFormatString(realf, hex_sum))
		}
	}

	return nil
}

func printSymlink(f string) error {
	assertFilePath(f)

	// debug print first
	if optDebug {
		if err := printDebug(f, SYMLINK); err != nil {
			return err
		}
	}

	// get hash value of symlink base name
	written, b, err := getStringHash(path.Base(f), optHashAlgo)
	if err != nil {
		return err
	}
	assert(len(b) > 0)
	hex_sum := getHexSum(b)

	// count this file
	appendStatTotal()
	appendWrittenTotal(written)
	appendStatSymlink(f)
	appendWrittenSymlink(written)

	// verify hash value if specified
	if len(optHashVerify) != 0 {
		if optHashVerify != hex_sum {
			return nil
		}
	}

	// squash or print this file
	if optHashOnly {
		if optSquash {
			updateSquashBuffer(b)
		} else {
			fmt.Println(hex_sum)
		}
	} else {
		if realf := getRealPath(f); optSquash {
			updateSquashBuffer(append([]byte(realf), b...))
		} else {
			fmt.Println(getXsumFormatString(realf, hex_sum))
		}
	}

	return nil
}

func printUnsupported(f string) error {
	if optDebug {
		if err := printDebug(f, UNSUPPORTED); err != nil {
			return err
		}
	}

	appendStatUnsupported(f)
	return nil
}

func printInvalid(f string) error {
	if optDebug {
		if err := printDebug(f, INVALID); err != nil {
			return err
		}
	}

	appendStatInvalid(f)
	return nil
}

func printDebug(f string, t fileType) error {
	assert(optDebug)
	if optAbs {
		var err error
		f, err = filepath.Abs(f)
		if err != nil {
			return err
		}
	}

	fmt.Println("###", f, getFileTypeString(t))
	return nil
}

func printVerboseStat() {
	indent := " "

	printNumFormatString(numStatTotal(), "file")
	a0 := numStatDirectory()
	a1 := numStatRegular()
	a2 := numStatDevice()
	a3 := numStatSymlink()
	assert(a0+a1+a2+a3 == numStatTotal())
	if a0 > 0 {
		fmt.Print(indent)
		printNumFormatString(a0, DIR_STR)
	}
	if a1 > 0 {
		fmt.Print(indent)
		printNumFormatString(a1, REG_STR)
	}
	if a2 > 0 {
		fmt.Print(indent)
		printNumFormatString(a2, DEVICE_STR)
	}
	if a3 > 0 {
		fmt.Print(indent)
		printNumFormatString(a3, SYMLINK_STR)
	}

	printNumFormatString(numWrittenTotal(), "byte")
	b0 := numWrittenDirectory()
	b1 := numWrittenRegular()
	b2 := numWrittenDevice()
	b3 := numWrittenSymlink()
	assert(b0+b1+b2+b3 == numWrittenTotal())
	if b0 > 0 {
		fmt.Print(indent)
		printNumFormatString(b0, DIR_STR+" byte")
	}
	if b1 > 0 {
		fmt.Print(indent)
		printNumFormatString(b1, REG_STR+" byte")
	}
	if b2 > 0 {
		fmt.Print(indent)
		printNumFormatString(b2, DEVICE_STR+" byte")
	}
	if b3 > 0 {
		fmt.Print(indent)
		printNumFormatString(b3, SYMLINK_STR+" byte")
	}

	printStatIgnored()
}

func assertFilePath(f string) {
	// must always handle file as abs
	assert(filepath.IsAbs(f))

	// file must not end with "/"
	assert(!strings.HasSuffix(f, "/"))

	// inputPrefix must not end with "/"
	assert(!strings.HasSuffix(inputPrefix, "/"))
}
