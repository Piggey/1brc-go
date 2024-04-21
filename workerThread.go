package main

const lineSeparator = ';'

func workerThread(fdata []byte, chunkChan <-chan chunk, resultChan chan<- map[string]stationData) {
	resultMap := map[string]stationData{}

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

			stationName := string(line[:splitIndex])
			temperature := parseFloatAsInt(line[splitIndex+1:])

			station, found := resultMap[stationName]
			if !found {
				resultMap[stationName] = NewStation(temperature)
				continue
			}

			resultMap[stationName] = stationData{
				min: min(temperature, station.min),
				max: max(temperature, station.max),
				sum: station.sum + temperature,
				cnt: station.cnt + 1,
			}
		}
	}

	resultChan <- resultMap
}
