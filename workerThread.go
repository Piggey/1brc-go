package main

import "github.com/zeebo/xxh3"

const lineSeparator = ';'

func workerThread(fdata []byte, chunkChan <-chan chunk, resultChan chan<- map[uint64]stationData) {
	resultMap := map[uint64]stationData{}

	for chunk := range chunkChan {
		chunkData := fdata[chunk.start:chunk.end]

		linestart := 0
		for i := 0; i < len(chunkData); i++ {
			if chunkData[i] != '\n' {
				continue
			}

			line := chunkData[linestart:i]
			splitIndex := findSplitIndex(line, lineSeparator)
			linestart = i + 1

			stationNameBytes := line[:splitIndex]
			stationHash := xxh3.Hash(line[:splitIndex])
			temperature := parseFloatAsInt(line[splitIndex+1:])

			station, found := resultMap[stationHash]
			if !found {
				resultMap[stationHash] = NewStation(stationNameBytes, temperature)
				continue
			}

			resultMap[stationHash] = stationData{
				name: stationNameBytes,
				min:  min(temperature, station.min),
				max:  max(temperature, station.max),
				sum:  station.sum + temperature,
				cnt:  station.cnt + 1,
			}
		}
	}

	resultChan <- resultMap
}
