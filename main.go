package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

const (
	exampleDataFile = "data/example.txt"
	dataFile        = "data/measurements.txt"
)

type stationData struct {
	min int
	max int
	sum int
	cnt int
}

func main() {
	cpuCores := runtime.NumCPU()
	fmt.Printf("cpuCores: %v\n", cpuCores)

	f, err := os.Open(dataFile)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	readerThread(f, cpuCores)
}
