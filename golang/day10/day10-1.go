package main

import (
	"os"
	"bufio"
	"fmt"
	"strconv"
	"regexp"
)

type Star struct {
	posX int
	posY int
	velocityX int
	velocityY int
}

func process(datafile string) []*Star {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	r_log, _ := regexp.Compile(`^position=<\s*(-?[0-9]\d*),\s*(-?[0-9]\d*)> velocity=<\s*(-?[0-9]\d*),\s*(-?[0-9]\d*)>$`)

	stars := make([]*Star,0)
	for s.Scan() {
		line := s.Text()
		res_log := r_log.FindStringSubmatch(line)

		if res_log != nil {
			s1, _ := strconv.Atoi(res_log[1])
			s2, _ := strconv.Atoi(res_log[2])
			s3, _ := strconv.Atoi(res_log[3])
			s4, _ := strconv.Atoi(res_log[4])
			s := Star{s1,s2,s3,s4 }
			stars = append(stars, &s)
		}
	}
	return stars
}

func starToString(s Star) string {
	return fmt.Sprintf("(x:%d,y:%d,vx:%d,vy:%d)",s.posX,s.posY,s.velocityX,s.velocityY)
}

func simulate(stars []*Star, seconds int) {

	for i := 0; i < seconds; i++ {
		for _, star := range stars {
			star.posX += star.velocityX
			star.posY += star.velocityY
		}
	}
}

func simulateTilRangeMet(stars []*Star, numSimulations int, rangeX int, rangeY int) {

	done := false
	i := 0
	for !done {
		for _, star := range stars {
			star.posX += star.velocityX
			star.posY += star.velocityY
		}

		minX, minY, maxX, maxY := findRange(stars)
		fmt.Printf("Range %d Diff %d,%d\n",i,(maxX - minX),(maxY - minY))
		if i == numSimulations || (maxX - minX) < rangeX  && (maxY - minY) < rangeY {
			done = true
		}
		i += 1
	}
}

func findRange(stars []*Star) (int, int, int, int) {
	minX, minY, maxX, maxY := 0, 0, 0, 0
	for _, star := range stars {
		if star.posX < minX {
			minX = star.posX
		}
		if star.posY < minY {
			minY = star.posY
		}
		if star.posX > maxX {
			maxX = star.posX
		}
		if star.posY > maxY {
			maxY = star.posY
		}
	}
	return minX, minY, maxX, maxY
}

func findStarAtCoord(x int, y int, stars []*Star) bool {
	for _, star := range stars {
		if star.posX == x && star.posY == y {
			return true
		}
	}
	return false
}

func printStars(stars []*Star) {

	// find minX, minY,maxX and maxY
	minX, minY, maxX, maxY := findRange(stars)

	fmt.Printf("Range: %d, %d, %d, %d\n",minX,minY,maxX,maxY)
	for i := minY; i <= maxY; i++ {
		for j := minX; j <= maxX; j++ {
			if findStarAtCoord(j, i, stars) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
}


func main() {
	stars := process("input.dat")
	simulateTilRangeMet(stars, 10311, 0,0)
	printStars(stars)
}

