package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	exampleDataFile = "example.txt"
	dataFile        = "measurements.txt"
)

type stationData struct {
	min int
	max int
	sum int
	cnt int
}

func main() {
	cpuCores := runtime.NumCPU()
	fmt.Printf("cpuCores: %v\n", cpuCores)

	f, err := os.Open(dataFile)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	dataCh := readerThread(f, cpuCores)
	resultCh := make(chan map[string]stationData, cpuCores)

	var wg sync.WaitGroup
	for i := 0; i < cpuCores; i++ {
		wg.Add(1)
		go func() {
			workerThread(dataCh, resultCh, &wg)
		}()
	}
	wg.Wait()
	close(resultCh)

	// collect and reduce data
	resultMap := initResultMap(minimumTemperature, maximumTemperature)

	for result := range resultCh {
		for stationName, data := range result {
			currentData := resultMap[stationName]
			resultMap[stationName] = stationData{
				min: min(data.min, currentData.min),
				max: max(data.max, currentData.max),
				sum: data.sum + currentData.sum,
				cnt: data.cnt + currentData.cnt,
			}
		}
	}

	// print result
	output := generateOutput(resultMap)
	fmt.Println(output)
}

func generateOutput(resultMap map[string]stationData) string {
	output := "{"

	for _, stationName := range stationsSorted {
		data := resultMap[stationName]
		minim := float64(data.min) / 10
		maxim := float64(data.max) / 10
		mean := float64(data.sum) / 10 / float64(data.cnt)

		output += fmt.Sprintf("%s=%.1f/%.1f/%.1f, ", stationName, minim, mean, maxim)
	}

	output, _ = strings.CutSuffix(output, ", ")
	output += "}"
	return output
}
