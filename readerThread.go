package main

import (
	"io"
	"os"
)

const (
	chunksize = 4 * 1024 * 1024
)

func readerThread(f *os.File, cores int) {
	var err error
	var read, lastNewlineIdx, offset int
	buf := make([]byte, chunksize)
	for {
		read, err = f.Read(buf[offset:])
		if read == 0 && err == io.EOF {
			break
		}

		datalen := read + offset
		lastNewlineIdx = lastIndexByte(buf, '\n', datalen)

		// move bytes that werent sent to front
		offset = datalen - lastNewlineIdx - 1
		copy(buf[:offset], buf[lastNewlineIdx+1:datalen])
	}
}

func lastIndexByte(buf []byte, c byte, start int) int {
	for i := start - 1; i >= 0; i-- {
		if buf[i] == c {
			return i
		}
	}

	return -1
}
