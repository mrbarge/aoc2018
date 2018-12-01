package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func readIn(datafile string) []int64 {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	frequencies := make([]int64, 0)
	for s.Scan() {
		line := s.Text()
		freq, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			fmt.Printf("Invalid frequency read: %s\n", line)
			continue
		}
		frequencies = append(frequencies, freq)
	}
	return frequencies
}

func applyFreq(frequencies []int64, startingValue int64) int64 {
	x := startingValue
	for i := 0; i < len(frequencies); i++ {
		x = x + frequencies[i]
	}
	return x
}

func main() {
	frequencies := readIn(os.Args[1])
	startingValue, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		startingValue = 0
	}
	sum := applyFreq(frequencies, startingValue)
	fmt.Printf("Finishing frequency is: %d\n", sum)
}
