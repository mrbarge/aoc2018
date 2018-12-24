package main

import (
	"bufio"
	"math"
	"os"
	"regexp"
	"strconv"
	"fmt"
)

type Bot struct {
	x int
	y int
	z int
	region int
}

func process(datafile string) {

	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	re_match, _ := regexp.Compile(`^pos=<(.+),(.+),(.+)>, r=(.+)$`)

	bots := make([]Bot,0)
	for s.Scan() {
		line := s.Text()
		res_match := re_match.FindStringSubmatch(line)

		if res_match != nil {
			pos1,_ := strconv.Atoi(res_match[1])
			pos2,_ := strconv.Atoi(res_match[2])
			pos3,_ := strconv.Atoi(res_match[3])

			rval,_ := strconv.Atoi(res_match[4])

			b := Bot{pos1,pos2,pos3, rval}
			bots = append(bots, b)
		}
	}

	maxBots := 0

	largestBot := 0
	maxRegion := 0
	for i := 0; i < len(bots); i++ {
		if bots[i].region > maxRegion {
			largestBot = i
			maxRegion = bots[i].region
		}
	}
	fmt.Printf("Largest radius is %d\n",maxRegion)


	for i := 0; i < len(bots); i++ {
		bcomp := bots[i]

		nx, ny, nz := getCoordDiffs(bots[largestBot],bcomp)
		if inRadius(nx,ny,nz,bots[largestBot].region) {
			maxBots += 1
		}
	}
	fmt.Printf("Part 1 max bots %d is %d\n",largestBot,maxBots)

}

func process2(datafile string) {

	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	re_match, _ := regexp.Compile(`^pos=<(.+),(.+),(.+)>, r=(.+)$`)

	// find ranges for all axes as that will define our search space
	minX, maxX, minY, maxY, minZ, maxZ := 0, 0, 0, 0, 0, 0

	bots := make([]Bot,0)
	first := true
	for s.Scan() {
		line := s.Text()
		res_match := re_match.FindStringSubmatch(line)

		if res_match != nil {
			pos1,_ := strconv.Atoi(res_match[1])
			pos2,_ := strconv.Atoi(res_match[2])
			pos3,_ := strconv.Atoi(res_match[3])
			rval,_ := strconv.Atoi(res_match[4])

			if first {
				minX, maxX = pos1, pos1
				minY, maxY = pos2, pos2
				minZ, maxZ = pos3, pos3
				first = false
			} else {
				if pos1 < minX {
					minX = pos1
				}
				if pos1 > maxX {
					maxX = pos1
				}
				if pos2 < minY {
					minY = pos2
				}
				if pos2 > maxY {
					maxY = pos2
				}
				if pos3 < minZ {
					minZ = pos3
				}
				if pos3 > maxZ {
					maxZ = pos3
				}
			}

			b := Bot{pos1,pos2,pos3, rval}
			bots = append(bots, b)
		}
	}

	deltaX := int(math.Abs(float64(maxX - minX)))
	deltaY := int(math.Abs(float64(maxY - minY)))
	deltaZ := int(math.Abs(float64(maxZ - minZ)))

	fmt.Printf("minx=%d,maxx=%d,miny=%d,maxy=%d,minz=%d,maxz=%d,deltaX=%d,deltaY=%d,deltaZ=%d\n",minX,maxX,
		minY,maxY,minZ,maxZ,deltaX,deltaY,deltaZ)

	riter := 1
	for riter < deltaX {
		riter *= 2
	}

	done := false
	for !done {
		bestDist := -1
		bestCount := 0
		bestX, bestY, bestZ := 0, 0, 0
		x, y, z := minX, minY, minZ
		for x <= maxX+1 {
			for y <= maxY+1 {
				for z <= maxZ+1 {
					//fmt.Printf("Repeat %d - ",z)
					bc := 0
					dummy := Bot{x,y,z,0}
					for _, bot := range bots {

						nx, ny, nz := getCoordDiffs(dummy, bot)
						newdist := int(math.Abs(float64(nx))) + int(math.Abs(float64(ny))) + int(math.Abs(float64(nz)))
						if (newdist - bot.region) / riter <= 0 {
							bc += 1
						}
					}

					if bc > bestCount {
						bestDist = int(math.Abs(float64(x))) + int(math.Abs(float64(y))) + int(math.Abs(float64(z)))
						bestCount = bc
						bestX, bestY, bestZ = x, y, z
					} else if bc == bestCount {
						calcDist := int(math.Abs(float64(x))) + int(math.Abs(float64(y))) + int(math.Abs(float64(z)))
						if bestDist < 0 || calcDist < bestDist {
							bestDist = calcDist
							bestX, bestY, bestZ = x, y, z
						}
					}
					z += riter
				}
				y += riter
			}
			x += riter
		}

		if riter == 1 {
			fmt.Printf("Best Count: %d\n",bestCount)
			fmt.Printf("Best coords: %d,%d,%d\n",bestX,bestY,bestZ)
			done = true
		} else {
			minX = bestX - riter
			maxX = bestX + riter
			minY = bestY - riter
			maxY = bestY + riter
			minZ = bestZ - riter
			maxZ = bestZ + riter
			riter /= 2
			fmt.Printf("%d,%d  %d,%d   %d,%d\n",minX,maxX,minY,maxY,minZ,maxZ)
		}
	}
}

func getCoordDiffs(largest Bot, test Bot) (int, int, int) {
	//return largest.x-test.x, largest.y-test.y,largest.z-test.z
	return int(math.Abs(float64(largest.x-test.x))),
		int(math.Abs(float64(largest.y-test.y))),
		int(math.Abs(float64(largest.z-test.z)))

}

func inRadius(x, y, z, radius int) bool {
	return (x^2 + y^2 + z^2) <= radius^2
}

func main() {
	process("input.dat")
	process2("input.dat")
}

