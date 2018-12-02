package main

import (
	"bufio"
	"fmt"
	"os"
)

func process(datafile string) int {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	lines := make([]string, 0)
	for s.Scan() {
		line := s.Text()
		lines = append(lines, line)
	}

	for _, i := range lines {
		for _, j := range lines {
			isCommon, commonLetters := commonIds(i, j)
			if isCommon {
				fmt.Printf("Common: %s\n",commonLetters)
			}
		}
	}
	return 0
}

func commonIds(i, j string) (bool, string) {
	foundOneDiff := false
	commonLetters := make([]uint8,0)
	for x := 0; x < len(i); x++ {
		if i[x] != j[x] {
			if foundOneDiff {
				return false, ""
			} else {
				foundOneDiff = true
			}
		} else {
			commonLetters = append(commonLetters, i[x])
		}
	}
	return foundOneDiff, string(commonLetters)
}

func main() {
	ans := process(os.Args[1])
	fmt.Println(ans)
}
