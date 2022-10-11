package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	inputPrefix string
	totalFiles  uint64
)

func printInput(f string, hash_algo string) error {
	// assert exists
	_, err := pathExists(f)
	if err != nil {
		return err
	}

	// make input abs, and keep base directory as prefix
	t, err := getFileType(f)
	if err != nil {
		return err
	}
	switch t {
	case DIR:
		f, err = filepath.Abs(f)
		if err != nil {
			return err
		}
		inputPrefix = f
	case REG:
		f, err = filepath.Abs(f)
		if err != nil {
			return err
		}
		inputPrefix = filepath.Dir(f)
	default:
		return fmt.Errorf("%s has unsupported type %d", f, t)
	}

	assertFilePath(f)

	// prefix is a directory
	t, _ = getFileType(inputPrefix)
	assert(t == DIR)

	// initialize global variables
	totalFiles = 0
	initSquashBuffer()

	// start directory walk
	err = filepath.WalkDir(f, getWalkDirHandler(hash_algo))

	// print squash hash if optSquash
	if err == nil && optSquash {
		s := getSquashBuffer()
		if optVerbose {
			fmt.Println(totalFiles, "files")
			fmt.Println(len(s), "bytes")
		}
		err = printString(f, s, hash_algo)
	}

	return err
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

		// find target if symlink
		var l string
		switch t {
		case SYMLINK:
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
			assert(t != SYMLINK)
		default:
			l = ""
		}

		switch t {
		case DIR:
			return nil
		case REG:
			return printReg(f, l, hash_algo)
		case UNSUPPORTED:
			return printUnsupported(f)
		case INVALID:
			return printInvalid(f)
		default:
			panic(fmt.Sprintf("%s has invalid type %d", f, t))
		}
		return nil
	}
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

func getXsumFormat(f string, h string) string {
	// compatible with shaXsum commands
	return fmt.Sprintf("%s  %s", h, f)
}

func printString(f string, s string, hash_algo string) error {
	assertFilePath(f)

	// get hash value
	b, err := getStringHash(s, hash_algo)
	if err != nil {
		return err
	}
	assert(len(b) > 0)

	// verify hash value if specified
	hex_sum := getHexSum(b)
	if optHashVerify != "" {
		if optHashVerify != hex_sum {
			return nil
		}
	}

	realf := getRealPath(f)
	if realf == "." {
		_, err = fmt.Println(hex_sum)
	} else {
		_, err = fmt.Println(getXsumFormat(realf, hex_sum))
	}

	return err
}

func printReg(f string, l string, hash_algo string) error {
	assertFilePath(f)

	// debug print first
	if optDebug {
		err := printDebug(f, "reg")
		if err != nil {
			return err
		}
	}

	// get hash value
	b, err := getFileHash(f, hash_algo)
	if err != nil {
		return err
	}
	assert(len(b) > 0)

	// verify hash value if specified
	hex_sum := getHexSum(b)
	if optHashVerify != "" {
		if optHashVerify != hex_sum {
			return nil
		}
	}

	// count this file
	totalFiles++

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

	// squash or print this file
	if optSquash {
		updateSquashBuffer(fmt.Sprintf("%s:%s", realf, hex_sum))
	} else {
		_, err = fmt.Println(getXsumFormat(realf, hex_sum))
	}

	return err
}

func printUnsupported(f string) error {
	if optDebug {
		err := printDebug(f, "unsupported")
		if err != nil {
			return err
		}
	}

	return nil
}

func printInvalid(f string) error {
	if optDebug {
		err := printDebug(f, "invalid")
		if err != nil {
			return err
		}
	}

	return nil
}

func printDebug(f string, s string) error {
	assert(optDebug)
	if optAbs {
		var err error
		f, err = filepath.Abs(f)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Println("###", f, s)
	return err
}

func assertFilePath(f string) {
	// must always handle file as abs
	assert(filepath.IsAbs(f))

	// file must not end with "/"
	assert(!strings.HasSuffix(f, "/"))

	// inputPrefix must not end with "/"
	assert(!strings.HasSuffix(inputPrefix, "/"))
}
