package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"golang.org/x/crypto/sha3"
	"hash"
	"io"
	"os"
)

func newHash(hash_algo string) hash.Hash {
	switch hash_algo {
	case "md5":
		return md5.New()
	case "sha1":
		return sha1.New()
	case "sha224":
		return sha256.New224()
	case "sha256":
		return sha256.New()
	case "sha384":
		return sha512.New384()
	case "sha512":
		return sha512.New()
	case "sha512_224":
		return sha512.New512_224()
	case "sha512_256":
		return sha512.New512_256()
	case "sha3_224":
		return sha3.New224()
	case "sha3_256":
		return sha3.New256()
	case "sha3_384":
		return sha3.New384()
	case "sha3_512":
		return sha3.New512()
	default:
		return nil
	}
}

func getFileHash(f string, hash_algo string) ([]byte, error) {
	fp, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	info, err := fp.Stat()
	if err != nil {
		return nil, err
	}

	written, b, err := getHash(fp, hash_algo)
	assert(written == info.Size())

	return b, err
}

func getHash(r io.Reader, hash_algo string) (int64, []byte, error) {
	h := newHash(hash_algo)
	if h == nil {
		return 0, nil, fmt.Errorf("invalid hash algorithm %s", hash_algo)
	}

	written, err := io.Copy(h, r)
	if err != nil {
		return 0, nil, err
	}

	return written, h.Sum(nil), nil
}

func getHexSum(sum []byte) string {
	return fmt.Sprintf("%x", sum)
}
