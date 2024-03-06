package main

func splitLine(line string, separator byte) (stationName, temperatureStr string) {
	// heurystyka zamiast iterowania
	// ;00.0  -> len - 5
	// ;-0.0  -> len - 5
	// ;0.0   -> len - 4
	// ;-00.0 -> len - 6

	l := len(line)
	switch {
	case line[l-5] == separator:
		return line[:l-5], line[l-4:]

	case line[l-4] == separator:
		return line[:l-4], line[l-3:]

	case line[l-6] == separator:
		return line[:l-6], line[l-5:]
	}

	panic(line)
}

func initResultMap(minTemp, maxTemp int) (resultMap map[string]stationData) {
	resultMap = make(map[string]stationData, len(stationsSorted))

	for _, stationName := range stationsSorted {
		resultMap[stationName] = stationData{
			min: maxTemp,
			max: minTemp,
		}
	}

	return resultMap
}

func parseFloatAsInt(s string) (num int) {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			num *= 10
			num += int(c - '0')
		}
	}

	if s[0] != '-' {
		return num
	}

	return -num
}
