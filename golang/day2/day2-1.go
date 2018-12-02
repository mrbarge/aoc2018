package main

import (
	"bufio"
	"fmt"
	"os"
)

func process(datafile string) int {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	digraphs := 0
	trigraphs := 0

	for s.Scan() {
		line := s.Text()
		if findFrequencyOfLength(line,2) {
			digraphs += 1
		}
		if findFrequencyOfLength(line, 3) {
			trigraphs += 1
		}
	}

	return digraphs * trigraphs
}

func findFrequencyOfLength(s string, l int) bool {
	freq := make(map[uint8]int)
	for i := 0; i < len(s); i++ {
		freq[s[i]] += 1
	}
	for _, v := range(freq) {
		if v == l {
			return true
		}
	}
	return false
}


func main() {
	ans := process(os.Args[1])
	fmt.Println(ans)
}
