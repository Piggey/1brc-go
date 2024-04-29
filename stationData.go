package main

type stationData struct {
	name []byte
	min  int
	max  int
	sum  int
	cnt  int
}

func NewStation(name []byte, temp int) stationData {
	return stationData{
		name: name,
		min:  temp,
		max:  temp,
		sum:  temp,
		cnt:  1,
	}
}
