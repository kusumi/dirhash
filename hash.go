package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"hash"
	"io"
	"os"
	"strings"
)

const (
	MD5        = "md5"
	SHA1       = "sha1"
	SHA224     = "sha224"
	SHA256     = "sha256"
	SHA384     = "sha384"
	SHA512     = "sha512"
	SHA512_224 = "sha512_224"
	SHA512_256 = "sha512_256"
	SHA3_224   = "sha3_224"
	SHA3_256   = "sha3_256"
	SHA3_384   = "sha3_384"
	SHA3_512   = "sha3_512"
)

func getAvailableHashAlgo() []string {
	return []string{
		MD5,
		SHA1,
		SHA224,
		SHA256,
		SHA384,
		SHA512,
		SHA512_224,
		SHA512_256,
		SHA3_224,
		SHA3_256,
		SHA3_384,
		SHA3_512,
	}
}

func newHash(hash_algo string) hash.Hash {
	switch hash_algo {
	case MD5:
		return md5.New()
	case SHA1:
		return sha1.New()
	case SHA224:
		return sha256.New224()
	case SHA256:
		return sha256.New()
	case SHA384:
		return sha512.New384()
	case SHA512:
		return sha512.New()
	case SHA512_224:
		return sha512.New512_224()
	case SHA512_256:
		return sha512.New512_256()
	case SHA3_224:
		return sha3.New224()
	case SHA3_256:
		return sha3.New256()
	case SHA3_384:
		return sha3.New384()
	case SHA3_512:
		return sha3.New512()
	default:
		return nil
	}
}

func getFileHash(f string, hash_algo string) (uint64, []byte, error) {
	fp, err := os.Open(f)
	if err != nil {
		return 0, nil, err
	}
	defer fp.Close()

	info, err := fp.Stat()
	if err != nil {
		return 0, nil, err
	}
	stat_size := uint64(info.Size())

	written, b, err := getHash(fp, hash_algo)
	assert(written == stat_size || stat_size == 0)

	return written, b, err
}

func getByteHash(s []byte, hash_algo string) (uint64, []byte, error) {
	r := bytes.NewReader(s)

	written, b, err := getHash(r, hash_algo)
	assert(written == uint64(len(s)))

	return written, b, err
}

func getStringHash(s string, hash_algo string) (uint64, []byte, error) {
	r := strings.NewReader(s)

	written, b, err := getHash(r, hash_algo)
	assert(written == uint64(len(s)))

	return written, b, err
}

func getHash(r io.Reader, hash_algo string) (uint64, []byte, error) {
	h := newHash(hash_algo)
	if h == nil {
		return 0, nil, fmt.Errorf("invalid hash algorithm %s", hash_algo)
	}

	written, err := io.Copy(h, r)
	if err != nil {
		return 0, nil, err
	}

	return uint64(written), h.Sum(nil), nil
}

func getHexSum(sum []byte) string {
	return hex.EncodeToString(sum)
}
