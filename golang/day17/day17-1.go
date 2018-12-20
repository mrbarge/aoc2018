package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Coord struct {
	x int
	y int
}

type WaterState struct {
	c Coord
	state FlowState
}

type FlowState int
const (
	FLOW_DOWN = iota
	FLOW_LEFT
	FLOW_RIGHT
	FILL
	REST
)

type Element rune
const (
	CLAY = '#'
	SAND = '.'
	WATER_FLOW = '|'
	WATER_REST = '~'
)

var data [][]Element

func process(datafile string) {

	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	lines := make([]string, 0)
	for s.Scan() {
		line := s.Text()
		lines = append(lines, line)
	}

	// first find the min and max y values
	mapmin := make(map[string]int)
	mapmax := make(map[string]int)
	mapmin["x"] = 500
	mapmin["y"] = 100
	for _, line := range lines {

		r_x_val, _ := regexp.Compile(`(x)=([\d\.]+)`)
		r_y_val, _ := regexp.Compile(`(y)=([\d\.]+)`)
		res_x := r_x_val.FindStringSubmatch(line)
		res_y := r_y_val.FindStringSubmatch(line)

		if res_x != nil {
			if strings.Contains(res_x[2], "..") {
				m1, m2 := getMinMax(res_x[2])
				if m1 < mapmin[res_x[1]] {
					mapmin[res_x[1]] = m1
				}
				if m2 > mapmax[res_x[1]] {
					mapmax[res_x[1]] = m2
				}
			} else {
				m, _ := strconv.Atoi(res_x[2])
				if m < mapmin[res_x[1]] {
					mapmin[res_x[1]] = m
				} else if m > mapmax[res_x[1]] {
					mapmax[res_x[1]] = m
				}
			}
		}

		if res_y != nil {
			if strings.Contains(res_y[2], "..") {
				m1, m2 := getMinMax(res_y[2])
				if m1 < mapmin[res_y[1]] {
					mapmin[res_y[1]] = m1
				}
				if m2 > mapmax[res_y[1]] {
					mapmax[res_y[1]] = m2
				}
			} else {
				m, _ := strconv.Atoi(res_y[2])
				if m < mapmin[res_y[1]] {
					mapmin[res_y[1]] = m
				} else if m > mapmax[res_y[1]] {
					mapmax[res_y[1]] = m
				}
			}
		}
	}

	// create our data plot
	offsetx := 5

	data = make([][]Element, mapmax["y"]+1)
	for y := 0; y < mapmax["y"]+1; y++ {
		// create an offsetx buffer on either side
		data[y] = make([]Element, mapmax["x"]+(offsetx*2))
		for i, _ := range data[y] {
			data[y][i] = SAND
		}
	}

	fmt.Printf("Size of grid is %d x %d\n",len(data),len(data[0]))
	// fill it with data
	for _, line := range lines {
		r_x_val, _ := regexp.Compile(`x=([\d\.]+)`)
		r_y_val, _ := regexp.Compile(`y=([\d\.]+)`)
		res_x := r_x_val.FindStringSubmatch(line)
		res_y := r_y_val.FindStringSubmatch(line)

		if res_x != nil && res_y != nil {
			if strings.Contains(res_x[1],"..") {
				// range of x for y
				minx1, minx2 := getMinMax(res_x[1])
				yval, _ := strconv.Atoi(res_y[1])
				for i := minx1; i < minx2; i++ {
					data[yval][i+offsetx] = CLAY
				}
			} else if strings.Contains(res_y[1], "..") {
				// range of y for x
				miny1, miny2 := getMinMax(res_y[1])
				xval, _ := strconv.Atoi(res_x[1])
				for i := miny1; i <= miny2; i++ {
					data[i][xval+offsetx] = CLAY
				}
			} else {
				// got to be just a 1x1 block
				xval, _ := strconv.Atoi(res_x[1])
				yval, _ := strconv.Atoi(res_y[1])
				data[yval][xval+offsetx] = CLAY
			}
		}
	}
	data[0][500+offsetx] = WATER_FLOW

	waterstack := make([]*WaterState,0)
	w := new(WaterState)
	w.c.x = 500+offsetx
	w.c.y = 0
	w.state = FLOW_DOWN
	waterstack = append(waterstack, w)
	done := false
	for !done {

		// get the most recent flow off the stack
		if len(waterstack) == 0 {
			done = true
			continue
		}
		latestFlow := waterstack[len(waterstack)-1]

		// has this breached the y max?
		if latestFlow.c.y >= len(data)-1 {
			// just remove from the stack
			waterstack = waterstack[:len(waterstack)-1]
			continue
		}

		switch latestFlow.state {
		case FLOW_DOWN:
			if data[latestFlow.c.y+1][latestFlow.c.x] == SAND {
				data[latestFlow.c.y+1][latestFlow.c.x] = WATER_FLOW
				newflow := new(WaterState)
				newflow.c.x = latestFlow.c.x
				newflow.c.y = latestFlow.c.y+1
				newflow.state = FLOW_DOWN
				waterstack = append(waterstack, newflow)
			}
			latestFlow.state = FLOW_RIGHT
			continue
		case FLOW_RIGHT:
			if (data[latestFlow.c.y+1][latestFlow.c.x] == WATER_REST || data[latestFlow.c.y+1][latestFlow.c.x] == CLAY) &&
				data[latestFlow.c.y][latestFlow.c.x+1] == SAND {
				data[latestFlow.c.y][latestFlow.c.x+1] = WATER_FLOW
				newflow := new(WaterState)
				newflow.c.x = latestFlow.c.x+1
				newflow.c.y = latestFlow.c.y
				newflow.state = FLOW_DOWN
				waterstack = append(waterstack, newflow)
			}
			latestFlow.state = FLOW_LEFT
			continue
		case FLOW_LEFT:
			if (data[latestFlow.c.y+1][latestFlow.c.x] == WATER_REST || data[latestFlow.c.y+1][latestFlow.c.x] == CLAY) &&
				data[latestFlow.c.y][latestFlow.c.x-1] == SAND {
				data[latestFlow.c.y][latestFlow.c.x-1] = WATER_FLOW
				newflow := new(WaterState)
				newflow.c.x = latestFlow.c.x-1
				newflow.c.y = latestFlow.c.y
				newflow.state = FLOW_DOWN
				waterstack = append(waterstack, newflow)
			}
			latestFlow.state = FILL
			continue
		case FILL:
			if inHole(latestFlow.c,len(data[0]),len(data)) {
				fillHole(latestFlow.c)
			}
			latestFlow.state = REST
			continue
		case REST:
			// remove from the stack, it cannae flow any longa capn!
			waterstack = waterstack[:len(waterstack)-1]
			continue
		}
	}
	fmt.Printf("PART ONE: %d\n", countWaters(data,mapmin["y"]))
	fmt.Printf("PART TWO: %d\n", countRestWaters(data,mapmin["y"]))

}

func countWaters(data [][]Element, starty int) int {
	retcount := 0
	for i := starty; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			if data[i][j]== WATER_REST || data[i][j] == WATER_FLOW {
				retcount += 1
			}
		}
	}
	return retcount
}

func countRestWaters(data [][]Element, starty int) int {
	retcount := 0
	for i := starty; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			if data[i][j]== WATER_REST {
				retcount += 1
			}
		}
	}
	return retcount
}

func inHole(c Coord, maxX int, maxY int) bool {
	// ensure we have clay along the bottom 'til we reach walls on either side
	// first move left
	inhole := true
	if c.y == maxY {
		return false
	}
	for i := c.x; i > 0 && data[c.y][i] != CLAY && inhole; i-- {
		if data[c.y+1][i] == SAND {
			inhole = false
		}
	}
	for i := c.x; i < maxX && data[c.y][i] != CLAY && inhole; i++ {
		if data[c.y+1][i] == SAND {
			inhole = false
		}
	}
	return inhole
}

func fillHole(c Coord) {
	for i := c.x; data[c.y][i] != CLAY ; i-- {
		data[c.y][i] = WATER_REST
	}

	for i := c.x; data[c.y][i] != CLAY ; i++ {
		data[c.y][i] = WATER_REST
	}
}

func printGrid(grid [][]Element) {
	for _, r := range grid {
		for _, v := range r {
			fmt.Print(string(v))
		}
		fmt.Println()
	}
}

func getMinMax(s string) (int,int) {
	v1, v2 := 0, 0
	r_val, _ := regexp.Compile(`(\d+)..(\d+)`)
	res_val := r_val.FindStringSubmatch(s)
	if res_val != nil {
		v1, _ = strconv.Atoi(res_val[1])
		v2, _ = strconv.Atoi(res_val[2])
		return v1, v2
	} else {
		return -1, -1
	}
}


func main() {
	process("input.dat")
	//process("test.dat")
}

