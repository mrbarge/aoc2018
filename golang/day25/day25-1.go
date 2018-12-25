package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Coord struct {
	w float64
	x float64
	y float64
	z float64
}

// returns true if immune win
func process(datafile string)  {

	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	coords := make([]*Coord,0)
	for s.Scan() {
		line := s.Text()

		c := new(Coord)
		fmt.Sscanf(line,"%f,%f,%f,%f",&c.w,&c.x,&c.y,&c.z)
		coords = append(coords, c)

	}

	// Create a slice of known constellations
	constellations := make([][]*Coord,0)

	// Add the first coord as the first constellation-in-progress
	con1 := make([]*Coord,0)
	con1 = append(con1,coords[0])
	constellations = append(constellations, con1)
	// remove it from our list of orphan coordinates
	coords = coords[1:]

	// while we still have orphan coordinates
	for len(coords) > 0 {
		// to avoid a stalemate of no existing constellations that a coord can be added to,
		// record when we don't add any coords to a constellation, so that we can
		// create a brand new one with it
		addedConstellation := false

		// for every constellation-in-progress, see if a coordinate is within distance
		for k := 0; k < len(constellations) && !addedConstellation; k++ {
			constellation := constellations[k]
			coordToRemove := -1
			for i := 0; i < len(constellation) && !addedConstellation; i++ {
				constellationCoord := constellation[i]
				// don't fiddle with the coords 'til we're done, otherwise it gets messy
				for j := 0; j < len(coords) && !addedConstellation; j++ {
					orphanedCoord := coords[j]
					if distance(constellationCoord, orphanedCoord) <= 3 {
						constellation = append(constellation, orphanedCoord)
						constellations[k] = constellation
						coordToRemove = j
						addedConstellation = true
					}
				}
			}
			if addedConstellation {
				copy(coords[coordToRemove:], coords[coordToRemove+1:])
				coords[len(coords)-1] = nil // or the zero value of T
				coords = coords[:len(coords)-1]
			}
		}

		if !addedConstellation {
			newcon := make([]*Coord,0)
			newcon = append(newcon, coords[0])
			constellations = append(constellations, newcon)
			coords = coords[1:]
		}
	}

	fmt.Printf("There were: %d constellations\n",len(constellations))
}

func distance(a *Coord, b *Coord) float64 {
	c := math.Abs(a.w - b.w) + math.Abs(a.x - b.x) + math.Abs(a.y - b.y) + math.Abs(a.z - b.z)
	return c
}

func main() {
	process("input.dat")
}

