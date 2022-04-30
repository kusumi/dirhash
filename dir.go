package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

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

func printDir(f string, hash_algo string) error {
	return filepath.WalkDir(f, getWalkDirHandler(hash_algo))
}

func printReg(f string, hash_algo string) error {
	exists, _ := pathExists(f)
	assert(exists)

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

	// convert file path to abs if specified
	if optAbs {
		f, err = filepath.Abs(f)
		if err != nil {
			return err
		}
	}

	// compatible with shaXsum commands
	_, err = fmt.Printf("%s  %s\n", hex_sum, f)
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

func getWalkDirHandler(hash_algo string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		t, err := getModeType(d.Type())
		if err != nil {
			return err
		}

		switch t {
		case DIR:
			assert(false)
		case REG:
			return printReg(path, hash_algo)
		case UNSUPPORTED:
			return printUnsupported(path)
		case INVALID:
			return printInvalid(path)
		default:
			assert(false)
		}
		return nil
	}
}

func printInput(input string, hash_algo string) error {
	t, err := getFileType(input)
	if err != nil {
		return err
	}

	switch t {
	case DIR:
		return printDir(input, hash_algo)
	case REG:
		return printReg(input, hash_algo)
	case UNSUPPORTED:
		return printUnsupported(input)
	case INVALID:
		return printInvalid(input)
	default:
		assert(false)
	}
	return nil
}
