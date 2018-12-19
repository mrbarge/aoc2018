package main

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Element rune
const (
	CLAY = '#'
	SAND = '.'
	WATER_FLOW = '|'
	WATER_REST = '~'
)

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
	offsetx := mapmin["x"]
	offsety := mapmin["y"]
	data := make([][]Element, mapmax["y"]-offsety)
	for y := 0; y < mapmax["y"]-offsety; y++ {
		data[y] = make([]Element, mapmax["x"]-offsetx)
		for i, _ := range data[y] {
			data[y][i] = SAND
		}
	}

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
					data[yval][i] = CLAY
				}
			} else if strings.Contains(res_y[1], "..") {
				// range of y for x
				miny1, miny2 := getMinMax(res_y[1])
				xval, _ := strconv.Atoi(res_x[1])
				for i := miny1; i < miny2; i++ {
					data[i][xval] = CLAY
				}
			} else {
				// got to be just a 1x1 block
				xval, _ := strconv.Atoi(res_x[1])
				yval, _ := strconv.Atoi(res_y[1])
				data[yval][xval] = CLAY
			}
		}
	}
}

func getMinMax(s string) (int,int) {
	v1, v2 := 0, 0
	r_val, _ := regexp.Compile(`.=(\d+)..(\d+)`)
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
	process("test.dat")
}

