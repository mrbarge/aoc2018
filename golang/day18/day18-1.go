package main

import (
	"bufio"
	"fmt"
	"os"
)

type Coord struct {
	x int
	y int
}

type Acre rune
const (
	OPEN = '.'
	TREE = '|'
	LUMBER = '#'
)

func process(datafile string, size int, minutes int) {

	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	acres := make([][]Acre,0)

	for s.Scan() {
		line := s.Text()
		l := make([]Acre, len(line))
		for i, v := range line {
			l[i] = Acre(v)
		}
		acres = append(acres, l)
	}

	for i := 0; i < minutes; i += 1 {
		acres = tick(acres)
	}

	fmt.Printf("Part 1: %d\n", countType(TREE, acres) * countType(LUMBER, acres))
}

func process2() {

	// on observing the output we know that the value at minute 1000 will reoccur
	// every 28 minutes

	// so we can figure out the nth iteration of the cycle it will be on at minute
	// at minute 1000000000 as follows

	idx := 0
	for i := 1000; i <= 999999999; i++ {
		idx += 1
		idx = idx % 28
	}
	fmt.Printf("index %d\n",idx)
}

func countType(acreType Acre, acres [][]Acre) int {
	retcount := 0
	for _, r := range acres {
		for _, v := range r {
			if v == acreType {
				retcount += 1
			}
		}
	}
	return retcount
}


	func tick(acres [][]Acre) [][]Acre {
	retacre := make([][]Acre, len(acres))
	for i, v := range acres {
		retacre[i] = make([]Acre,len(v))
		for j, v2 := range v {
			retacre[i][j] = v2
		}
	}

	for i, r := range acres {
		for j, v := range r {
			switch v {
			case OPEN:
				if equalAdjacents(Coord{i,j}, TREE, acres) >= 3 {
					retacre[i][j] = TREE
				}
			case TREE:
				if equalAdjacents(Coord{i,j}, LUMBER, acres) >= 3 {
					retacre[i][j] = LUMBER
				}
			case LUMBER:
				if equalAdjacents(Coord{i,j}, TREE, acres) >= 1 &&
					equalAdjacents(Coord{i,j}, LUMBER, acres) >= 1	{
					retacre[i][j] = LUMBER
				} else {
					retacre[i][j] = OPEN
				}
			}
		}
	}
	return retacre
}

func equalAdjacents(c Coord, acreType Acre, acres [][]Acre) int {
	retcount := 0
	coords := getNeighbours(c)
	maxSize := len(acres)
	for _, coord := range coords {
		if coord.x >= 0 && coord.x < maxSize &&
			coord.y >= 0 && coord.y < maxSize &&
			acres[coord.x][coord.y] == acreType {
			retcount += 1
		}
	}
	return retcount
}

func getNeighbours(c Coord) []Coord {
	return []Coord {
		Coord{c.x - 1, c.y},
		Coord{c.x - 1, c.y + 1},
		Coord{c.x - 1, c.y - 1},
		Coord{c.x, c.y - 1},
		Coord{c.x, c.y + 1},
		Coord{c.x + 1, c.y},
		Coord{c.x + 1, c.y + 1},
		Coord{c.x + 1, c.y - 1},
	}


}

func printAcres(acres [][]Acre) {
	for _, r := range acres {
		for _, v := range r {
			fmt.Printf(string(v))
		}
		fmt.Println("")
	}
}

func main() {
	process("input.dat", 50, 10)
	process2()
}

