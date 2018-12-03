package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"regexp"
)

type Claim struct {
	id int
	startX int
	startY int
	width int
	height int
}

func process(datafile string) [][]*Claim {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	r_claim, _ := regexp.Compile(`^#(\d+) @ (\d+),(\d+): (\d+)x(\d+)$`)

	fabric_map := make([][]*Claim,1000)
	for i := range(fabric_map) {
		fabric_map[i] = make([]*Claim,1000)
	}
	for i := range(fabric_map) {
		for j := range(fabric_map[i]) {
			fabric_map[i][j] = nil
		}
	}

	multi_claim := Claim{-1,0,0,0,0}
	for s.Scan() {
		line := s.Text()
		res_claim := r_claim.FindStringSubmatch(line)

		if res_claim != nil {
			startX, _ := strconv.Atoi(res_claim[2])
			startY, _  := strconv.Atoi(res_claim[3])
			width, _  := strconv.Atoi(res_claim[4])
			height, _  := strconv.Atoi(res_claim[5])
			claim_id, _ := strconv.Atoi(res_claim[1])
			c := Claim{claim_id, startX, startY, width, height}

			// check if any element of the claim is already claimed, if so
			// flag that both claims need to be set as invalid
			foundAnotherClaim := false
			old_claims := make([]Claim,0)
			for i := startX; i < startX+width; i++ {
				for j := startY; j < startY+height; j++ {
					if fabric_map[i][j] != nil  {
						foundAnotherClaim = true
						old_claims = append(old_claims,*fabric_map[i][j])
					} else {
						fabric_map[i][j] = &c
					}
				}
			}

			if foundAnotherClaim {
				for _, old_claim := range old_claims {
					for i := old_claim.startX; i < old_claim.startX+old_claim.width; i++ {
						for j := old_claim.startY; j < old_claim.startY+old_claim.height; j++ {
							fabric_map[i][j] = &multi_claim
						}
					}
				}
				for i := c.startX; i < c.startX+c.width; i++ {
					for j := c.startY; j < c.startY+c.height; j++ {
						fabric_map[i][j] = &multi_claim
					}
				}
			}
		}
	}

	return fabric_map
}

func find_isolated_claim(claims [][]*Claim) *Claim {
	for i := range(claims) {
		for j := range(claims[i]) {
			if claims[i][j] != nil && claims[i][j].id > 0 {
				return claims[i][j]
			}
		}
	}
	return new(Claim)
}

func main() {
	claims := process(os.Args[1])
	answer := find_isolated_claim(claims)
	fmt.Printf("Answer is: %d\n", answer.id)
}
