package main

import (
	"bufio"
	"log"
	"os"
)

func splitLine(line string, separator byte) (stationName, temperatureStr string) {
	for i := 0; i < len(line); i++ {
		if line[i] == separator {
			return line[:i], line[i+1:]
		}
	}
	return "very", "wrong"
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

func fastParseFloat(s string) float64 {
	mantissa, exp, neg := fastReadFloat(s)
	return atof64(mantissa, exp, neg)
}

func atof64(mantissa uint64, exp int, neg bool) float64 {
	f := float64(mantissa)
	if neg {
		f = -f
	}

	return f / float64pow10[-exp]
}

func fastReadFloat(s string) (mantissa uint64, exp int, neg bool) {
	i := 0
	if s[i] == '-' {
		neg = true
		i += 1
	}

	base := uint64(10)
	var nd, dp, ndMant int
	for ; i < len(s); i++ {
		switch c := s[i]; {
		case c == '.':
			dp = nd
			continue

		case c >= '0' && c <= '9':
			nd += 1
			mantissa *= base
			mantissa += uint64(c - '0')
			ndMant += 1
			continue
		}
	}

	exp = dp - ndMant
	return mantissa, exp, neg
}

var float64pow10 = []float64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	1e20, 1e21, 1e22,
}
