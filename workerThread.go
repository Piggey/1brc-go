package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	stationsFile       = "stations.txt"
	minimumTemperature = -100
	maximumTemperature = 100
	uniqueStationsNum  = 413
)

func workerThread(threadId int, dataCh <-chan []byte, resultCh chan<- map[string]stationMapData, wg *sync.WaitGroup) {
	defer wg.Done()
	resultMap := initResultMap(stationsFile, uniqueStationsNum, minimumTemperature, maximumTemperature)

	for data := range dataCh {
		linestart := 0
		for i := 0; i < len(data); i++ {
			if data[i] != '\n' {
				continue
			}

			temp := strings.Split(string(data[linestart:i]), lineSeparator)
			linestart = i + 1
			stationName := temp[0]
			temperature, _ := strconv.ParseFloat(temp[1], 64)

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

func initResultMap(stationsFile string, stationsNumber int, minimumTemperature, maximumTemperature float64) map[string]stationMapData {
	out := make(map[string]stationMapData, stationsNumber)

	f, err := os.Open(stationsFile)
	if err != nil {
		log.Panic(err)
	}

	scn := bufio.NewScanner(f)
	for scn.Scan() {
		station := scn.Text()

		out[station] = stationMapData{
			min: maximumTemperature,
			max: minimumTemperature,
		}
	}

	return out
}
