package main

import (
	"sync"
)

const (
	stationsFile       = "stations.txt"
	minimumTemperature = -1000
	maximumTemperature = 1000
	lineSeparator      = ';'
)

func workerThread(dataCh <-chan []byte, resultCh chan<- map[string]stationData, wg *sync.WaitGroup) {
	defer wg.Done()

	resultMap := initResultMap(minimumTemperature, maximumTemperature)

	for data := range dataCh {
		linestart := 0
		for i := 0; i < len(data); i++ {
			if data[i] != '\n' {
				continue
			}

			stationName, temperatureStr := splitLine(string(data[linestart:i]), lineSeparator)
			linestart = i + 1
			temperature := parseFloatAsInt(temperatureStr)

			data := resultMap[stationName]
			resultMap[stationName] = stationData{
				min: min(temperature, data.min),
				max: max(temperature, data.max),
				sum: data.sum + temperature,
				cnt: data.cnt + 1,
			}
		}
	}

	resultCh <- resultMap
}
