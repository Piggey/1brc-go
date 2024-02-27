package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	lineSeparator      = ";"
	chunksize          = 128 * 1024 * 1024
	exampleDataFile    = "example.txt"
	dataFile           = "measurements.txt"
	stationsFile       = "stations.txt"
	minimumTemperature = -100
	maximumTemperature = 100
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

	minMap, maxMap := loadCachedStations(stationsFile, minimumTemperature, maximumTemperature)
	sumMap := map[string]float64{}
	cntMap := map[string]float64{}

	buf := make([]byte, chunksize)
	line := make([]byte, 0, 256)
	for {
		read, err := f.Read(buf)
		if read == 0 && err == io.EOF {
			break
		}

		for i := 0; i < read; i++ {
			if buf[i] != '\n' {
				line = append(line, buf[i])
				continue
			}

			temp := strings.Split(string(line), lineSeparator)
			stationName := temp[0]
			temperature, _ := strconv.ParseFloat(temp[1], 64)

			minMap[stationName] = min(temperature, minMap[stationName])
			maxMap[stationName] = max(temperature, maxMap[stationName])
			sumMap[stationName] += temperature
			cntMap[stationName] += 1
			line = line[:0]
		}
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

func loadCachedStations(stationsFile string, minimumTemperature, maximumTemperature float64) (minMap, maxMap map[string]float64) {
	minMap = map[string]float64{}
	maxMap = map[string]float64{}

	f, err := os.Open(stationsFile)
	if err != nil {
		log.Panic(err)
	}

	scn := bufio.NewScanner(f)
	for scn.Scan() {
		station := scn.Text()

		minMap[station] = maximumTemperature
		maxMap[station] = minimumTemperature
	}

	return minMap, maxMap
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
