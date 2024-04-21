package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"slices"
	"strings"
	"sync"
	"syscall"
)

const (
	dataFile = "data/example.txt"
	// dataFile = "data/measurements.txt"
)

func main() {
	cpuCores := runtime.NumCPU()
	fmt.Printf("cpuCores: %v\n", cpuCores)

	f, err := os.Open(dataFile)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	fstat, err := f.Stat()
	if err != nil {
		log.Panic(err)
	}

	fsize := int(fstat.Size())
	fdata, err := syscall.Mmap(int(f.Fd()), 0, fsize, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		log.Panic(err)
	}

	chunkChan := createChunks(fdata, fsize, cpuCores)
	resultChan := make(chan map[string]stationData, cpuCores)

	go func() {
		var wg sync.WaitGroup
		for i := 0; i < cpuCores; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				workerThread(fdata, chunkChan, resultChan)
			}()
		}
		wg.Wait()
		close(resultChan)
	}()

	// collect and reduce data
	resultMap := reduceResults(resultChan)

	// print result
	output := generateOutput(resultMap)
	fmt.Println(output)
}

func reduceResults(resultChan <-chan map[string]stationData) map[string]stationData {
	resultMap := map[string]stationData{}

	for result := range resultChan {
		for stationName, station := range result {
			resStation, found := resultMap[stationName]
			if !found {
				resultMap[stationName] = station
				continue
			}

			resultMap[stationName] = stationData{
				min: min(station.min, resStation.min),
				max: max(station.max, resStation.max),
				sum: station.sum + resStation.sum,
				cnt: station.cnt + resStation.cnt,
			}
		}
	}

	return resultMap
}

func generateOutput(resultMap map[string]stationData) string {
	// sort station stationNames
	stationNames := make([]string, 0, len(resultMap))
	for stationName := range resultMap {
		stationNames = append(stationNames, stationName)
	}
	slices.Sort(stationNames)

	output := "{"

	for _, stationName := range stationNames {
		station := resultMap[stationName]

		mini := float64(station.min) / 10
		maxi := float64(station.max) / 10
		mean := float64(station.sum) / float64(station.cnt*10)

		output += fmt.Sprintf("%s=%.1f/%.1f/%.1f, ", stationName, mini, mean, maxi)
	}

	output, _ = strings.CutSuffix(output, ", ")
	output += "}"
	return output
}
