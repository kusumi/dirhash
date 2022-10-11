//go:build squash2

package main

var (
	squashBuffer string
)

func initSquashBuffer() {
	squashBuffer = ""
}

func updateSquashBuffer(s string) {
	// result depends on append order
	b, err := getStringHash(squashBuffer+s, SHA1)
	if err != nil {
		panic(err)
	}
	squashBuffer = getHexSum(b)
}

func getSquashBuffer() string {
	return squashBuffer
}
