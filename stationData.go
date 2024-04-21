package main

type stationData struct {
	min int
	max int
	sum int
	cnt int
}

func NewStation(temp int) stationData {
	return stationData{
		min: temp,
		max: temp,
		sum: temp,
		cnt: 1,
	}
}
