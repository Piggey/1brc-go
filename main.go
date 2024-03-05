package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"slices"
	"strings"
	"sync"
)

const (
	exampleDataFile = "example.txt"
	dataFile        = "measurements.txt"
)

type (
	stationMapData struct {
		min float64
		max float64
		sum float64
		cnt int
	}

	stationSliceData struct {
		name string
		min  float64
		max  float64
		mean float64
	}
)

func main() {
	cpuCores := runtime.NumCPU()
	log.Printf("cpuCores: %v\n", cpuCores)

	f, err := os.Open(dataFile)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	dataCh := readerThread(f)
	resultCh := make(chan map[string]stationMapData, cpuCores)

	log.Println("tworzenie worker threadow")
	var wg sync.WaitGroup
	for i := 0; i < cpuCores; i++ {
		wg.Add(1)
		go workerThread(i, dataCh, resultCh, &wg)
	}
	wg.Wait()
	close(resultCh)

	// collect and reduce data
	log.Println("zbieranie wynikow")
	resultMap := make(map[string]stationMapData)
	for result := range resultCh {
		for stationName, data := range result {
			stationData := resultMap[stationName]
			resultMap[stationName] = stationMapData{
				min: min(data.min, stationData.min),
				max: max(data.max, stationData.max),
				sum: data.sum + stationData.sum,
				cnt: data.cnt + stationData.cnt,
			}
		}
	}

	// print result
	log.Println("konwertowanie do listy")
	stationsData := convertToArray(resultMap)
	slices.SortFunc(stationsData, func(a, b stationSliceData) int {
		if a.name < b.name {
			return -1
		}
		return 1
	})

	log.Println("wypisywanie wyniku")
	output := generateOutput(stationsData)
	fmt.Println(output)
}

func generateOutput(stationsData []stationSliceData) string {
	output := "{"
	for _, s := range stationsData {
		output += fmt.Sprintf("%s=%.1f/%.1f/%.1f, ", s.name, s.min, s.mean, s.max)
	}

	output, _ = strings.CutSuffix(output, ", ")
	output += "}"
	return output
}

func convertToArray(m map[string]stationMapData) []stationSliceData {
	out := make([]stationSliceData, 0, len(m))

	for stationName, data := range m {
		out = append(out, stationSliceData{
			name: stationName,
			min:  data.min,
			mean: data.sum / float64(data.cnt),
			max:  data.max,
		})
	}

	return out
}
