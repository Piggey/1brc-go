package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	lineSeparator      = ";"
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
	f, err := os.Open("measurements.txt")
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	filescanner := bufio.NewScanner(f)
	filescanner.Split(bufio.ScanLines)

	minMap := map[string]float64{}
	maxMap := map[string]float64{}
	sumMap := map[string]float64{}
	cntMap := map[string]float64{}

	for filescanner.Scan() {
		line := strings.Split(filescanner.Text(), lineSeparator)
		stationName := line[0]
		temperature, _ := strconv.ParseFloat(line[1], 64)

		if _, ok := minMap[stationName]; !ok {
			minMap[stationName] = maximumTemperature
		}

		if _, ok := maxMap[stationName]; !ok {
			maxMap[stationName] = minimumTemperature
		}

		if temperature < minMap[stationName] {
			minMap[stationName] = temperature
		}

		if temperature > maxMap[stationName] {
			maxMap[stationName] = temperature
		}

		sumMap[stationName] += temperature
		cntMap[stationName] += 1
	}

	stationsData := convertToArray(minMap, maxMap, sumMap, cntMap)

	slices.SortFunc(stationsData, func(a, b stationData) int {
		if a.name < b.name {
			return -1
		}
		return 1
	})

	// print output
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
