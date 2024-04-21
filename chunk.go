package main

type chunk struct {
	start int
	end   int
}

func createChunks(fdata []byte, fsize int, numCores int) <-chan chunk {
	ch := make(chan chunk, numCores)

	go func() {
		defer close(ch)

		// split to chunks and align to newlines
		chunkSize := fsize / numCores
		chunkStart := 0
		chunkEnd := 0
		for c := 0; c < numCores; c++ {
			chunkEnd = findByte(fdata, fsize, '\n', min(chunkStart+chunkSize, fsize-1))

			ch <- chunk{
				start: chunkStart,
				end:   chunkEnd,
			}

			chunkStart = chunkEnd + 1
		}
	}()

	return ch
}

func findByte(fdata []byte, fsize int, b byte, start int) int {
	for i := start; i < fsize; i++ {
		if fdata[i] == b {
			return i
		}
	}
	panic("err")
}
