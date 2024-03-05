package main

import (
	"sync"
)

const (
	stationsFile       = "stations.txt"
	minimumTemperature = -100
	maximumTemperature = 100
	uniqueStationsNum  = 413
	lineSeparator      = ';'
)

func workerThread(dataCh <-chan []byte, resultMap map[string]stationMapData, resultCh chan<- map[string]stationMapData, wg *sync.WaitGroup) {
	defer wg.Done()

	for data := range dataCh {
		linestart := 0
		for i := 0; i < len(data); i++ {
			if data[i] != '\n' {
				continue
			}

			stationName, temperatureStr := splitLine(string(data[linestart:i]), lineSeparator)
			linestart = i + 1
			temperature := fastParseFloat(temperatureStr)

			stationData := resultMap[stationName]
			resultMap[stationName] = stationMapData{
				min: min(temperature, stationData.min),
				max: max(temperature, stationData.max),
				sum: stationData.sum + temperature,
				cnt: stationData.cnt + 1,
			}
		}
	}

	resultCh <- resultMap
}
