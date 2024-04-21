package main

func findSplitIndex(line []byte, separator byte) (index int) {
	// heurystyka zamiast iterowania
	// ;00.0  -> len - 5
	// ;-0.0  -> len - 5
	// ;0.0   -> len - 4
	// ;-00.0 -> len - 6

	l := len(line)
	switch {
	case line[l-5] == separator:
		index = l - 5

	case line[l-4] == separator:
		index = l - 4

	case line[l-6] == separator:
		index = l - 6
	}

	return index
}

func parseFloatAsInt(s []byte) (num int) {
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
