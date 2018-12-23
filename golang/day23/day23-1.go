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
/*		if inRadius(int64(bcomp.x),int64(bcomp.y),int64(bcomp.z),int64(bots[largestBot].region)) {
			maxBots += 1
		}*/
	}
	fmt.Printf("Max bots %d is %d\n",largestBot,maxBots)

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
	//process("test.dat")
}

