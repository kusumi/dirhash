package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	inputPrefix string
)

func printInput(f string, hash_algo string) error {
	// assert exists
	_, err := pathExists(f)
	if err != nil {
		return err
	}

	// convert input to abs first
	f, err = filepath.Abs(f)
	if err != nil {
		return err
	}
	assertFilePath(f)

	// keep input prefix based on raw type
	t, err := getFileType(f)
	if err != nil {
		return err
	}
	switch t {
	case DIR:
		inputPrefix = f
	case REG:
		fallthrough
	case DEVICE:
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
	err = filepath.WalkDir(f, getWalkDirHandler(hash_algo))
	if err != nil {
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
			printNumFormatString(uint64(len(b)), "squashed byte")
		}
		err = printByte(f, b, hash_algo)
		if err != nil {
			return err
		}
	}

	return nil
}

func getWalkDirHandler(hash_algo string) fs.WalkDirFunc {
	return func(f string, d fs.DirEntry, err error) error {
		assertFilePath(f)

		if err != nil {
			return err
		}

		t, err := getRawFileType(f)
		if err != nil {
			return err
		}

		if testIgnoreEntry(f, t) {
			appendStatIgnored(f)
			return nil
		}

		// find target if symlink
		var l string // symlink itself, not its target
		switch t {
		case SYMLINK:
			if optIgnoreSymlink {
				appendStatIgnored(f)
				return nil
			}
			if optLstat {
				return printSymlink(f, hash_algo)
			}
			l = f
			f, err = os.Readlink(f)
			if err != nil {
				return err
			}
			if !filepath.IsAbs(f) {
				f = filepath.Join(filepath.Dir(l), f)
				assert(filepath.IsAbs(f))
			}
			t, err = getFileType(f)
			if err != nil {
				return err
			}
			assert(t != SYMLINK) // symlink chains resolved
		default:
			l = ""
		}

		switch t {
		case DIR:
			// A regular directory isn't considered ignored,
			// then don't count symlink to directory as ignored.
			//if len(l) > 0 {
			//	appendStatIgnored(l)
			//}
			return nil
		case REG:
			fallthrough
		case DEVICE:
			return printFile(f, l, t, hash_algo)
		case UNSUPPORTED:
			return printUnsupported(f)
		case INVALID:
			return printInvalid(f)
		default:
			panicFileType(f, "unknown", t)
		}

		assert(false)
		return nil
	}
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

func getRealPath(f string) string {
	if optAbs {
		assert(filepath.IsAbs(f))
		return f
	} else if f == inputPrefix {
		return "."
	} else if inputPrefix == "/" {
		return f[1:]
	} else if strings.HasPrefix(f, inputPrefix) {
		f = f[len(inputPrefix)+1:]
		assert(!strings.HasPrefix(f, "/"))
		return f
	} else {
		return f // f is probably symlink target
	}
}

func printByte(f string, inb []byte, hash_algo string) error {
	assertFilePath(f)

	// get hash value
	_, b, err := getByteHash(inb, hash_algo)
	if err != nil {
		return err
	}
	assert(len(b) > 0)
	hex_sum := getHexSum(b)

	// verify hash value if specified
	if optHashVerify != "" {
		if optHashVerify != hex_sum {
			return nil
		}
	}

	if optHashOnly {
		fmt.Println(hex_sum)
	} else {
		if realf := getRealPath(f); realf == "." {
			fmt.Println(hex_sum)
		} else {
			fmt.Println(getXsumFormatString(realf, hex_sum))
		}
	}

	return nil
}

func printFile(f string, l string, t fileType, hash_algo string) error {
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
	written, b, err := getFileHash(f, hash_algo)
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
	if optHashVerify != "" {
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
				if strings.HasPrefix(l, inputPrefix) {
					l = l[len(inputPrefix)+1:]
					assert(!strings.HasPrefix(l, "/"))
				}
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

func printSymlink(f string, hash_algo string) error {
	assertFilePath(f)

	// debug print first
	if optDebug {
		if err := printDebug(f, SYMLINK); err != nil {
			return err
		}
	}

	// get a symlink string to get hash value
	s, err := os.Readlink(f)
	if err != nil {
		return err
	}

	// get hash value
	written, b, err := getStringHash(s, hash_algo)
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
	if optHashVerify != "" {
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
		// hash value is from s, but print realf path for clarity
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
	a1 := numStatRegular()
	a2 := numStatDevice()
	a3 := numStatSymlink()
	assert(a1+a2+a3 == numStatTotal())
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
	b1 := numWrittenRegular()
	b2 := numWrittenDevice()
	b3 := numWrittenSymlink()
	assert(b1+b2+b3 == numWrittenTotal())
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
