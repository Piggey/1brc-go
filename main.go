package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	lineSeparator   = ";"
	blocksize       = 4096
	exampleDataFile = "example.txt"
	dataFile        = "measurements.txt"
)

type stationData struct {
	name string
	min  float64
	mean float64
	max  float64
}

func main() {
	f, err := os.Open(dataFile)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	minMap := map[string]float64{}
	maxMap := map[string]float64{}
	sumMap := map[string]float64{}
	cntMap := map[string]float64{}

	buf := make([]byte, blocksize)
	line := make([]byte, 0, 256)
	for {
		read, err := f.Read(buf)
		if read == 0 && err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}

		for i := 0; i < read; i++ {
			if buf[i] != '\n' {
				line = append(line, buf[i])
				continue
			}

			temp := strings.Split(string(line), lineSeparator)
			stationName := temp[0]
			temperature, _ := strconv.ParseFloat(temp[1], 64)

			if min, ok := minMap[stationName]; !ok {
				minMap[stationName] = temperature
			} else {
				if temperature < min {
					minMap[stationName] = temperature
				}
			}

			if max, ok := maxMap[stationName]; !ok {
				maxMap[stationName] = temperature
			} else {
				if temperature > max {
					maxMap[stationName] = temperature
				}
			}

			sumMap[stationName] += temperature
			cntMap[stationName] += 1
			line = line[:0]
		}
	}

	for station := range cntMap {
		fmt.Println(station)
	}

	stationsData := convertToArray(minMap, maxMap, sumMap, cntMap)
	slices.SortFunc(stationsData, func(a, b stationData) int {
		if a.name < b.name {
			return -1
		}
		return 1
	})

	output := generateOutput(stationsData)
	fmt.Println(output)
}

func generateOutput(stationsData []stationData) string {
	output := "{"
	for _, s := range stationsData {
		output += fmt.Sprintf("%s=%.1f/%.1f/%.1f, ", s.name, s.min, s.mean, s.max)
	}

	output, _ = strings.CutSuffix(output, ", ")
	output += "}"
	return output
}

func convertToArray(minMap, maxMap, sumMap, cntMap map[string]float64) []stationData {
	out := make([]stationData, 0, len(minMap))

	for station := range minMap {
		mean := sumMap[station] / cntMap[station]
		out = append(out, stationData{
			name: station,
			min:  minMap[station],
			mean: mean,
			max:  maxMap[station],
		})
	}

	return out
}
