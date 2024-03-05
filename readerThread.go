package main

import (
	"io"
	"os"
)

const (
	chunksize = 32 * 1024 * 1024
)

func readerThread(f *os.File, cores int) <-chan []byte {
	ch := make(chan []byte, cores)

	go func() {
		defer close(ch)

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

			// references bad!!
			bufCopy := make([]byte, len(buf))
			copy(bufCopy, buf)
			ch <- bufCopy[:lastNewlineIdx+1]

			// move bytes that werent sent to front
			offset = datalen - lastNewlineIdx - 1
			copy(buf[:offset], buf[lastNewlineIdx+1:datalen])
		}
	}()

	return ch
}

func lastIndexByte(buf []byte, c byte, start int) int {
	for i := start - 1; i >= 0; i-- {
		if buf[i] == c {
			return i
		}
	}

	return -1
}
