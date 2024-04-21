package main

const (
	stationsFile       = "stations.txt"
	minimumTemperature = -1000
	maximumTemperature = 1000
	lineSeparator      = ';'
)

func workerThread(dataCh <-chan []byte, resultCh chan<- map[string]stationData) {
	resultMap := initResultMap(minimumTemperature, maximumTemperature)

	for data := range dataCh {
		linestart := 0
		for i := 0; i < len(data); i++ {
			if data[i] != '\n' {
				continue
			}

			line := data[linestart:i]
			splitIdx := splitIndex(line, lineSeparator)
			linestart = i + 1

			stationName := line[:splitIdx]
			temperature := parseFloatAsInt(line[splitIdx+1:])

			data := resultMap[string(stationName)]
			resultMap[string(stationName)] = stationData{
				min: min(temperature, data.min),
				max: max(temperature, data.max),
				sum: data.sum + temperature,
				cnt: data.cnt + 1,
			}
		}
	}

	resultCh <- resultMap
}
