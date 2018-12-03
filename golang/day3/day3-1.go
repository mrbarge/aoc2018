package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"regexp"
)

func process(datafile string) [][]int {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	r_claim, _ := regexp.Compile(`^#(\d+) @ (\d+),(\d+): (\d+)x(\d+)$`)

	fabric_map := make([][]int,1000)
	for i := range(fabric_map) {
		fabric_map[i] = make([]int,1000)
	}

	for s.Scan() {
		line := s.Text()
		res_claim := r_claim.FindStringSubmatch(line)

		if res_claim != nil {
			startX, _ := strconv.Atoi(res_claim[2])
			startY, _  := strconv.Atoi(res_claim[3])
			width, _  := strconv.Atoi(res_claim[4])
			height, _  := strconv.Atoi(res_claim[5])

			for i := startX; i < startX+width; i++ {
				for j := startY; j < startY+height; j++ {
					fabric_map[i][j] += 1
				}
			}
		}
	}

	return fabric_map
}

func count_multi_claims(claims [][]int) int {

	multi := 0
	for i := range(claims) {
		for j := range(claims[i]) {
			if claims[i][j] > 1 {
				multi += 1
			}
		}
	}
	return multi
}

func main() {
	claims := process(os.Args[1])
	answer := count_multi_claims(claims)
	fmt.Printf("Answer is: %d\n", answer)
}
