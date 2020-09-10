package helpers

import (
	"bytes"
	"runtime"
	"strconv"
)

// GetGID - get Gorouting ID from go enviroment
func GetGID() uint64 {
	// Scott Mansfield
	// Goroutine IDs
	// https://blog.sgmansfield.com/2015/12/goroutine-ids/
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
