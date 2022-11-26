package hasher

import (
	"crypto/md5"
	"fmt"
	"os"
)

var BUFSIZE = 1048576

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Hash1M(path string) string {
	f, err := os.Open(path)
	check(err)
	defer f.Close()

	// read first BUFSIZE bytes
	buf := make([]byte, BUFSIZE)
	_, err = f.Read(buf)
	check(err)

	// compute hash of (size, first BUFSIZE bytes)
	hash := md5.Sum(buf)

	return string(fmt.Sprintf("%x", hash))
}
