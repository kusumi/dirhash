//go:build squash2

package main

var (
	squashBuffer []byte
)

func initSquashBuffer() {
	squashBuffer = make([]byte, 0)
}

func updateSquashBuffer(b []byte) {
	// result depends on append order
	_, tmp, err := getByteHash(append(squashBuffer, b...), SHA1)
	if err != nil {
		panic(err)
	}
	squashBuffer = tmp
}

func getSquashBuffer() []byte {
	return squashBuffer
}
