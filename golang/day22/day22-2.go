package main

import (
	"fmt"
	"github.com/albertorestifo/dijkstra"
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

	maxX := targetX+100
	maxY := targetY+100

	grid := make([][]*Element, maxY)

	// fill out zero range first
	for iy := 0; iy < maxY; iy++ {
		r := make([]*Element, maxX)
		grid[iy] = r
		for ix := 0; ix < maxX; ix++ {
			e := new(Element)
			grid[iy][ix] = e
			e.c = Coord{ix,iy}
			e.geologic = geologic(e.c, grid)
			e.erosion = erosion(e.c, depth, grid)
			e.region = region(e.c, grid)

		}
	}

	// and the rest
	for iy := 1; iy < maxY; iy++ {
		for ix := 1; ix < maxX; ix++ {
			e := new(Element)
			grid[iy][ix] = e
			e.c = Coord{ix,iy}
			e.geologic = geologic(e.c, grid)
			e.erosion = erosion(e.c, depth, grid)
			e.region = region(e.c, grid)
		}
	}

	// overwrite target to be rocky
	grid[10][10].region = ROCKY

	gr := buildGraph(grid)
	targetCoord := fmt.Sprintf("%d-%d-1",targetX,targetY)
	path, weight, err := gr.Path("0-0-1",targetCoord)
	fmt.Println(err)
	fmt.Print("Path: ")
	fmt.Println(path)
	fmt.Printf("Weight: %d\n",weight)
	fmt.Println(grid[10][10].region)


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

func getValidTools(r Region) []Tool {
	switch r {
	case WET:
		return []Tool{CLIMBING, NEITHER}
	case NARROW:
		return []Tool{TORCH, NEITHER}
	case ROCKY:
		return []Tool{CLIMBING, TORCH}
	default:
		return []Tool{}
	}
}

func getNeighbourCoords(c Coord) []Coord {
	return []Coord {
		Coord{c.x, c.y - 1},
		Coord{c.x - 1,c.y},
		Coord{c.x + 1,c.y},
		Coord{c.x, c.y + 1},
	}
}

func buildGraph(grid [][]*Element) dijkstra.Graph {

	g := dijkstra.Graph{}
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[0]); x++ {

			tools := getValidTools(grid[y][x].region)

			t1s := fmt.Sprintf("%d-%d-%d",x,y,tools[0])
			t2s := fmt.Sprintf("%d-%d-%d",x,y,tools[1])
			g[t1s] = make(map[string]int,0)
			g[t2s] = make(map[string]int,0)
			g[t1s][t2s] = 7
			g[t2s][t1s] = 7
			fmt.Printf("Added self connection %s %s\n",t1s,t2s)

			neighbours := getNeighbourCoords(Coord{x,y})
			for _, n := range neighbours {
				if n.x >= 0 && n.x < len(grid[0]) && n.y >= 0 && n.y < len(grid) {

					nexttools := getValidTools(grid[n.y][n.x].region)

					for _, ct := range tools {
						for _, nt := range nexttools {
							if ct == nt {
								// can bring current tool into next node
								source := fmt.Sprintf("%d-%d-%d",x,y,ct)
								neighbour := fmt.Sprintf("%d-%d-%d",n.x,n.y,ct)
								_, ok := g[neighbour]
								if !ok {
									g[neighbour] = make(map[string]int,0)
								}
								fmt.Printf("Added connection %s %s\n",source,neighbour)
								g[source][neighbour] = 1
							}
						}
					}
				}
			}
		}
	}
	return g
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
	//process(510, 10, 10)
}

