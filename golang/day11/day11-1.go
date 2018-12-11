package main

import (
	"fmt"
	"strconv"
)

func processProblem1(serial int, dimX int, dimY int) (int, int) {

	maxPower := 0
	maxX := 0
	maxY := 0

	grid := make([][]int, dimX)
	for x, _ := range grid {
		grid[x] = make([]int, dimY)
	}

	for x := 0; x < dimX; x++ {
		for y := 0; y < dimY; y++ {
			grid[x][y] = getPowerLevel(x, y, serial)
		}
	}

	for x := 0; x < dimX-3; x++ {
		for y := 0; y < dimY-3; y++ {
			p := getTotalPower(x, y, 3, 3, grid)
			if p > maxPower {
				maxPower = p
				maxX = x
				maxY = y
			}
		}
	}

	return maxX, maxY
}

func processProblem2(serial int, dimX int, dimY int) (int, int, int) {

	maxPower := 0
	maxX := 0
	maxY := 0
	size := 0

	grid := make([][]int, dimX)
	for x, _ := range grid {
		grid[x] = make([]int, dimY)
	}

	for x := 0; x < dimX; x++ {
		for y := 0; y < dimY; y++ {
			grid[x][y] = getPowerLevel(x, y, serial)
		}
	}

	// iterate over grid size / 2 because squares can't be bigger than that
	for i := 0; i < dimX / 2; i++ {
		for x := 0; x < dimX; x++ {
			for y := 0; y < dimY; y++ {
				p := getTotalPower(x, y, i, i, grid)
				if p > maxPower {
					maxPower = p
					maxX = x
					maxY = y
					size = i
				}
			}
		}
	}

	return maxX, maxY, size
}

func getTotalPower(x int, y int, mx int, my int, grid [][]int) int {
	if x + mx > len(grid) || y + my > len(grid) {
		return -1
	}
	p := 0
	for i := x; i < x+mx; i++ {
		for j := y; j < y+my; j++ {
			p += grid[i][j]
		}
	}
	return p
}

func getPowerLevel(x int, y int, serial int) int {
	rackId := x + 10
	powerLevel := rackId * y
	powerLevel += serial
	powerLevel = powerLevel * rackId
	// Self-acknowledgement of laziness here!
	sp := strconv.Itoa(powerLevel)
	if len(sp) >= 3 {
		powerLevel, _ = strconv.Atoi(string(sp[len(sp)-3]))
	} else {
		powerLevel = 0
	}
	powerLevel -= 5
	return powerLevel
}

func main() {
	x, y, size := processProblem2(7165, 300, 300)
//	x, y, size := processProblem2(18, 300, 300)
	fmt.Printf("%d %d %d\n",x, y, size)

}

