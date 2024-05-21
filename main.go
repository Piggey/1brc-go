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
	// dataFile = "data/example.txt"
	dataFile = "data/measurements.txt"
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
	resultChan := make(chan map[uint64]stationData, cpuCores)

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

func reduceResults(resultChan <-chan map[uint64]stationData) map[uint64]stationData {
	resultMap := map[uint64]stationData{}

	for result := range resultChan {
		for stationHash, station := range result {
			resStation, found := resultMap[stationHash]
			if !found {
				resultMap[stationHash] = station
				continue
			}

			resultMap[stationHash] = stationData{
				name: resStation.name,
				min:  min(station.min, resStation.min),
				max:  max(station.max, resStation.max),
				sum:  station.sum + resStation.sum,
				cnt:  station.cnt + resStation.cnt,
			}
		}
	}

	return resultMap
}

func generateOutput(resultMap map[uint64]stationData) string {
	// sort station stationNames
	type hashName struct {
		hash uint64
		name []byte
	}

	hashNames := make([]hashName, 0, len(resultMap))
	for stationHash, station := range resultMap {
		hashNames = append(hashNames, hashName{stationHash, station.name})
	}
	slices.SortFunc(hashNames, func(a, b hashName) int {
		return slices.Compare(a.name, b.name)
	})

	output := "{"

	for _, hn := range hashNames {
		station := resultMap[hn.hash]

		mini := float64(station.min) / 10
		maxi := float64(station.max) / 10
		mean := float64(station.sum) / float64(station.cnt*10)

		output += fmt.Sprintf("%s=%.1f/%.1f/%.1f, ", hn.name, mini, mean, maxi)
	}

	output, _ = strings.CutSuffix(output, ", ")
	output += "}"
	return output
}
