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
	stationsNumber     = 413
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
	f, err := os.Open(dataFile)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("info.Size(): %v\n", info.Size())

	resultMap := initResultMap(stationsFile, stationsNumber, minimumTemperature, maximumTemperature)

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

			stationData := resultMap[stationName]
			resultMap[stationName] = stationMapData{
				min: min(temperature, stationData.min),
				max: max(temperature, stationData.max),
				sum: stationData.sum + temperature,
				cnt: stationData.cnt + 1,
			}
			line = line[:0]
		}
	}

	stationsData := convertToArray(resultMap)
	slices.SortFunc(stationsData, func(a, b stationSliceData) int {
		if a.name < b.name {
			return -1
		}
		return 1
	})

	output := generateOutput(stationsData)
	fmt.Println(output)
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
