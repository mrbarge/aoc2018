package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
	"math"
	"sort"
)

type Pair struct {
	id int
	x int
	y int
}

func process(datafile string) {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	largest := 0
	pairs := make([]Pair,0)
	for s.Scan() {
		line := s.Text()
		elems := strings.Split(line,", ")
		n1, _ := strconv.Atoi(elems[0])
		n2, _ := strconv.Atoi(elems[1])

		if n1 > largest {
			largest = n1
		} else if n2 > largest {
			largest = n2
		}

		pairs = append(pairs,Pair{len(pairs)+1, n1, n2})
	}

	data := make([][]int,largest+1)
	for i, _ := range data {
		data[i] = make([]int, largest+1)
	}

	for _, l := range(pairs) {
		data[l.x][l.y] = l.id
	}

	infinitePairs := make([]int, 0)
	for x := 0; x < len(data); x++ {
		for y := 0; y < len(data); y++ {
			p, dupe := findShortestDistance(x,y,pairs)
			if dupe {
				data[x][y] = 0
			} else {
				if x == 0 || y == 0 || x == len(data)-1 || y == len(data) - 1 {
					infinitePairs = append(infinitePairs, p.id)
				}
				data[x][y] = p.id
			}
		}
	}

	maxArea := 0
	var maxPair Pair
	_ = maxPair
	for _, p := range pairs {
		if contains(infinitePairs, p.id) {
			continue
		} else {
			area := countArea(p, data)
			if area > maxArea {
				maxArea = area
				maxPair = p
			}
		}

	}

	fmt.Printf("The biggest is %d with area %d\n",maxPair.id,maxArea)
}

func countArea(p Pair, data [][]int) int {
	c := 0
	for x, _ := range data {
		for y, _ := range data[x] {
			if data[x][y] == p.id {
				c += 1
			}
		}
	}
	return c
}

func findShortestDistance(x int, y int, pairs []Pair) (Pair, bool) {
	smallest := float64(99999)
	p := Pair{0,0,0}

	distances := make([]float64,0)
	for _, pair := range pairs {
		if (x == pair.x && y == pair.y) {
			return pair, false
		}
		d := math.Abs(float64(x-pair.x)) + math.Abs(float64(y-pair.y))
		distances = append(distances,d)
		if d < smallest {
			smallest = d
			p = pair
		}
	}

	sort.Float64s(distances)
	return p, (distances[0] == distances[1])
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	process("input.dat")
}
