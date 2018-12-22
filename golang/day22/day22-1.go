package main

import (
	"fmt"
)

type Region int
const (
	ROCKY = iota
	WET
	NARROW
)

type Tool int
const (
	CLIMBING = iota
	TORCH
	NEITHER
)

type Element struct {
	region Region
	c Coord
	geologic int
	erosion int
}

type Coord struct {
	x int
	y int
}

func process(depth int, targetX int, targetY int) {

	maxX := targetX+1
	maxY := targetY+1

	grid := make([][]*Element, maxY)

	// fill out zero range first
	for iy := 0; iy <= targetY; iy++ {
		r := make([]*Element, maxX)
		grid[iy] = r
		for ix := 0; ix <= targetX; ix++ {
			e := new(Element)
			grid[iy][ix] = e
			e.c = Coord{ix,iy}
			e.geologic = geologic(e.c, grid)
			e.erosion = erosion(e.c, depth, grid)
			e.region = region(e.c, grid)

		}
	}

	// and the rest
	for iy := 1; iy <= targetY; iy++ {
		for ix := 1; ix <= targetX; ix++ {
			e := new(Element)
			grid[iy][ix] = e
			e.c = Coord{ix,iy}
			e.geologic = geologic(e.c, grid)
			e.erosion = erosion(e.c, depth, grid)
			e.region = region(e.c, grid)
		}
	}

	fmt.Println("Part 1: %d", risk(grid,targetX,targetY))
	printGrid(grid,targetX,targetY)
}

func printGrid(grid [][]*Element, targetX int, targetY int) {
	for iy := 0; iy <= targetY; iy++ {
		for ix := 0; ix <= targetX; ix++ {
			if (ix == 0 && iy == 0) {
				fmt.Print("M")
			} else {
				if grid[iy][ix].region == WET {
					fmt.Print("=")
				} else if grid[iy][ix].region == NARROW {
					fmt.Print("|")
				} else if grid[iy][ix].region == ROCKY {
					fmt.Print(".")
				}
			}
		}
		fmt.Println()
	}

}

func risk(grid [][]*Element, targetX int, targetY int) int {

	retcnt := 0
	for iy := 0; iy <= targetY; iy++ {
		for ix := 0; ix <= targetX; ix++ {
			if (ix == 0 && iy == 0) {
				continue
			}
			if (ix == targetX && iy == targetY) {
				continue
			}
			if grid[iy][ix].region == WET {
				retcnt += 1
			} else if grid[iy][ix].region == NARROW {
				retcnt += 2
			}
		}
	}
	return retcnt
}

func geologic(c Coord, grid [][]*Element) int {
	if c.x == 0 && c.y == 0 {
		return 0
	}

	if c.x == 14 && c.y == 796 {
		return 0
	}

	if c.y == 0 {
		return c.x * 16807
	}

	if c.x == 0 {
		return c.y * 48271
	}

	return grid[c.y-1][c.x].erosion * grid[c.y][c.x-1].erosion
}

func erosion(c Coord, depth int, grid [][]*Element) int {
	return (grid[c.y][c.x].geologic + depth) % 20183
}

func region(c Coord, grid [][]*Element) Region {
	if grid[c.y][c.x].erosion % 3 == 0 {
		return ROCKY
	} else if grid[c.y][c.x].erosion % 3 == 1 {
		return WET
	} else {
		return NARROW
	}
}

func main() {
	process(5355, 14, 796)
}

