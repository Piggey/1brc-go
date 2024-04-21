package main

const lineSeparator = ';'

func workerThread(dataCh <-chan []byte, resultCh chan<- map[string]stationData) {
	resultMap := map[string]stationData{}

	for data := range dataCh {
		linestart := 0
		for i := 0; i < len(data); i++ {
			if data[i] != '\n' {
				continue
			}

			line := data[linestart:i]
			splitIdx := splitIndex(line, lineSeparator)
			linestart = i + 1

			stationName := string(line[:splitIdx])
			temperature := parseFloatAsInt(line[splitIdx+1:])

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

	resultCh <- resultMap
}
