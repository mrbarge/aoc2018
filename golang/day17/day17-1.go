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

	r_v_val, _ := regexp.Compile(`([xy])=([\d\.]+)`)

	lines := make([]string, 0)
	for s.Scan() {
		line := s.Text()
		lines = append(lines, line)
	}

	// first find the min and max y values
	mapmin := make(map[string]int)
	mapmax := make(map[string]int)
	for _, line := range lines {
		res_v := r_v_val.FindStringSubmatch(line)

		if strings.Contains(res_v[2], "..") {
			m1, m2 := getMinMax(line)
			if m1 < mapmin[res_v[1]] {
				mapmax[res_v[1]] = m1
			}
			if m2 > mapmax[res_v[1]] {
				mapmax[res_v[1]] = m2
			}
		} else {
			m, _ := strconv.Atoi(res_v[2])
			if m < mapmin[res_v[1]] {
				mapmin[res_v[1]] = m
			} else if m > mapmax[res_v[1]] {
				mapmax[res_v[1]] = m
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
				for i := miny1; i < miny2; i++ {
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
			fmt.Println("Creamed?")
			done = true
			continue
		}
		latestFlow := waterstack[len(waterstack)-1]
		fmt.Printf("Flowing at (%d,%d)\n",latestFlow.c.x,latestFlow.c.y)

		// has this breached the y max?
		if latestFlow.c.y >= len(data)-1 {
			// just remove from the stack
			waterstack = waterstack[:len(waterstack)-1]
			continue
		}

		switch latestFlow.state {
		case FLOW_DOWN:
			if data[latestFlow.c.y+1][latestFlow.c.x] == SAND {
				fmt.Println("Flowing down to new node")
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
				fmt.Println("Flowing right to new node")
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
				fmt.Println("Flowing left to new node")
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
				fmt.Println("Filling level")
				fillHole(latestFlow.c)
			}
			latestFlow.state = REST
			continue
		case REST:
			// remove from the stack, it cannae flow any longa capn!
			fmt.Println("No more\n")
			waterstack = waterstack[:len(waterstack)-1]
			continue
		}
	}
/*
	waters := make([]*WaterState,0)
	waters = append(waters, new(WaterState{Coord{500,0}, true}))
	done := false

	//printGrid(data)

	for !done {

		data, waters = tick(data, waters)

		// check if water has reached limits
		if checkFinished(waters, mapmax["y"]) {
			done = true
		}
	}
*/
	printGrid(data)
	fmt.Println("PART UNO: %d", countWaters(data))
}

func countWaters(data [][]Element) int {
	retcount := 0
	for _, r := range data {
		for _, v := range r {
			if v == WATER_REST || v == WATER_FLOW {
				retcount += 1
			}
		}
	}
	return retcount
}

func flow(c Coord, maxX int, maxY int) {

	fmt.Printf("Flowing with maxX %d, maxY %d and coord (%d,%d)\n",len(data[0]),len(data),c.x,c.y)
	if c.y >= maxY {
		return
	}

 	if data[c.y+1][c.x] == SAND {
		data[c.y+1][c.x] = WATER_FLOW
		flow(Coord{c.x,c.y+1}, maxX, maxY)
	}

	if data[c.y+1][c.x] == WATER_REST || data[c.y+1][c.x] == CLAY &&
		data[c.y][c.x+1] == SAND {
			data[c.y][c.x+1] = WATER_FLOW
			flow(Coord{c.x+1, c.y}, maxX, maxY)
	}
	if data[c.y+1][c.x] == WATER_REST || data[c.y+1][c.x] == CLAY &&
		data[c.y][c.x-1] == SAND {
		data[c.y][c.x-1] = WATER_FLOW
		flow(Coord{c.x-1, c.y}, maxX, maxY)
	}
	if inHole(c,maxX,maxY) {
		fillHole(c)
	}

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

/*
func tick(data [][]Element, waters []*WaterState) ([][]Element, []*WaterState) {

	retdata := make([][]Element, len(data))
	for i, v := range data {
		retdata[i] = make([]Element,len(v))
		for j, v2 := range v {
			retdata[i][j] = v2
		}
	}

	for _, w := range waters {
		if !w.flowing {
			continue
		}

		neighbours := getValidNeighbours(w.c, data, len(data[0]), len(data))
//		fmt.Printf("Returned %d neighbours\n",len(neighbours))
		// if there is clay below, can we flow to the left or right?
		dc, isdown := neighbours["down"]
		lc, isleft := neighbours["left"]
		rc, isright := neighbours["right"]
		uc, isup := neighbours["up"]
		if isdown {
			waters = append(waters, &WaterState{dc, true})
			retdata[dc.y][dc.x] = WATER_FLOW
			fmt.Printf("Adding new down water at (%d,%d)\n", dc.x, dc.y)
		} else {
			if isleft {
				waters = append(waters, &WaterState{lc, true})
				retdata[lc.y][lc.x] = WATER_FLOW
				fmt.Printf("Adding new left water at (%d,%d)\n", lc.x, lc.y)
			}
			if isright {
				waters = append(waters, &WaterState{rc, true})
				retdata[rc.y][rc.x] = WATER_FLOW
				fmt.Printf("Adding new right water at (%d,%d)\n", rc.x, rc.y)
			}

			// if we can't flow down, left or right, can we pool upwards?
			if !isleft && !isright && isup {
				waters = append(waters, &WaterState{uc, true})
				retdata[uc.y][uc.x] = WATER_FLOW
				fmt.Printf("Adding new up water at (%d,%d)\n", uc.x, uc.y)
			}
		}
	}

	return retdata, waters
}
*/
func getValidNeighbours(c Coord, data [][]Element, maxX int, maxY int) map[string]Coord {
	r := make(map[string]Coord,0)
	if c.x > 0 && data[c.y][c.x-1] == SAND {
		r["left"] = Coord{c.x-1, c.y}
	}
	if c.y < maxY && data[c.y+1][c.x] == CLAY {
		if c.x < maxX && data[c.y][c.x+1] == SAND {
			r["right"] = Coord{c.x + 1, c.y}
		}
		if c.y < maxY && data[c.y+1][c.x] == SAND {
			r["down"] = Coord{c.x, c.y + 1}
		}
		if c.y > 0 && data[c.y-1][c.x] == SAND {
			r["up"] = Coord{c.x, c.y - 1}
		}
	}
	return r
}

func checkFinished(waters []*WaterState, maxy int) bool {
	for _, w := range waters {
		if w.c.y == maxy {
			return true
		}
	}
	return false
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
	//process("input.dat")
	process("test.dat")
}

