package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const measurementsFilename = "measurements.txt"

type stationMeasurements struct {
	min   float32
	mean  float32
	max   float32
	count int
}

type result []string

func (r result) GoString() (s string) {
	for i, value := range r {
		if i == 0 {
			s = fmt.Sprintf("{%s, ", value)
		} else if i == len(r)-1 {
			s += fmt.Sprintf("%s}", value)
		} else {
			s += fmt.Sprintf("%s, ", value)
		}
	}
	return
}

func main() {
	file, err := os.Open(measurementsFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stations := make(map[string]stationMeasurements)

	r := bufio.NewScanner(file)
	for r.Scan() {
		line := r.Text()
		processLine(line, &stations)
	}

	result := make(result, 0, len(stations))
	for station, measurements := range stations {
		result = append(result, fmt.Sprintf("%s=%.1f/%.1f/%.1f", station, measurements.min, measurements.mean, measurements.max))
	}
	sort.Strings(result)

	fmt.Printf("%#v\n", result)
}

func processLine(line string, stations *map[string]stationMeasurements) {
	lineInfo := strings.Split(strings.TrimSpace(line), ";")
	stationName := lineInfo[0]
	tempStr, err := strconv.ParseFloat(lineInfo[1], 32)
	if err != nil {
		panic(err)
	}
	temp := float32(tempStr)

	measurement, exists := (*stations)[stationName]
	if !exists {
		measurement = stationMeasurements{
			max:   temp,
			mean:  temp,
			min:   temp,
			count: 0,
		}
	}
	measurement.count++

	if temp > measurement.max {
		measurement.max = temp
	}
	if temp < measurement.min {
		measurement.min = temp
	}
	measurement.mean = (temp + (measurement.mean * float32(measurement.count-1))) / float32(measurement.count)
	(*stations)[stationName] = measurement
}
