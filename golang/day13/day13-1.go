package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Cart struct {
	id int
	dir Direction
	x int
	y int
	lastTurn LastTurn
	isCrashed bool
}

type LastTurn int
const (
	TURN_LEFT = iota
	STRAIGHT
	TURN_RIGHT
)

type Direction int
const (
	UP = iota
	DOWN
	LEFT
	RIGHT
)

type Track int
const (
	VERTICAL = iota
	HORIZONTAL
	FWD_TURN
	BACK_TURN
	INTERSECTION
)

func tick(carts []*Cart, grid [][]int) {

	sort.Slice(carts, func(i, j int) bool {
		if carts[i].y < carts[j].y {
			return true
		}
		if carts[i].y > carts[j].y {
			return false
		}
		return carts[i].x < carts[j].x
	})

	for _, cart := range carts {
		switch grid[cart.x][cart.y] {
			case VERTICAL:
				switch cart.dir {
					case UP:
						cart.y -= 1
					case DOWN:
						cart.y += 1
				}
			case HORIZONTAL:
				switch cart.dir {
					case LEFT:
						cart.x -= 1
					case RIGHT:
						cart.x += 1
				}
			case FWD_TURN:
				switch cart.dir {
					case LEFT:
						cart.y += 1
						cart.dir = DOWN
					case UP:
						cart.x += 1
						cart.dir = RIGHT
					case DOWN:
						cart.x -= 1
						cart.dir = LEFT
					case RIGHT:
						cart.y -= 1
						cart.dir = UP
				}
			case BACK_TURN:
				switch cart.dir {
					case LEFT:
						cart.y -= 1
						cart.dir = UP
					case UP:
						cart.x -= 1
						cart.dir = LEFT
					case DOWN:
						cart.x += 1
						cart.dir = RIGHT
					case RIGHT:
						cart.y += 1
						cart.dir = DOWN
				}
			case INTERSECTION:
				switch cart.lastTurn {
					case TURN_LEFT:
						// going STRAIGHT
						switch cart.dir {
						case LEFT:
							cart.x -= 1
						case UP:
							cart.y -= 1
						case DOWN:
							cart.y += 1
						case RIGHT:
							cart.x += 1
						}
						cart.lastTurn = STRAIGHT
					case STRAIGHT:
						// going RIGHT
						switch cart.dir {
						case LEFT:
							cart.y -= 1
							cart.dir = UP
						case UP:
							cart.x += 1
							cart.dir = RIGHT
						case DOWN:
							cart.x -= 1
							cart.dir = LEFT
						case RIGHT:
							cart.y += 1
							cart.dir = DOWN
						}
						cart.lastTurn = TURN_RIGHT
					case TURN_RIGHT:
						// going LEFT
						switch cart.dir {
						case LEFT:
							cart.y += 1
							cart.dir = DOWN
						case UP:
							cart.x -= 1
							cart.dir = LEFT
						case DOWN:
							cart.x += 1
							cart.dir = RIGHT
						case RIGHT:
							cart.y -= 1
							cart.dir = UP
						}
						cart.lastTurn = TURN_LEFT
				}
		}
		findCollisionCarts(carts)
	}

}

func findCollisionCarts(carts []*Cart) (bool,[]*Cart) {
	isCrash := false
	crashCarts := make([]*Cart,0)
	for _, c1 := range carts {
		if c1.isCrashed {
			continue
		}
		for _, c2 := range carts {
			if c1 != c2 && !c2.isCrashed && c1.x == c2.x && c1.y == c2.y {
				fmt.Printf("Crash at %d,%d (%d,%d)\n",c1.x,c1.y,c1.id,c2.id)
				c1.isCrashed = true
				c2.isCrashed = true
				isCrash = true
				crashCarts = append(crashCarts,c1)
				crashCarts = append(crashCarts,c2)
			}
		}
	}
	return isCrash,crashCarts
}

func findCollision(carts []*Cart) (bool,int,int) {
	for _, c1 := range carts {
		for _, c2 := range carts {
			if c1 != c2 && c1.x == c2.x && c1.y == c2.y {
				return true,c1.x,c1.y
			}
		}
	}
	return false,-1,-1
}

func process(datafile string) {

	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	grid := make([][]int, 151)
	for i, _ := range grid {
		grid[i] = make([]int, 151)
	}

	y := 0
	numcarts := 0
	carts := make([]*Cart, 0)
	for s.Scan() {
		line := s.Text()

		for x, c := range line {
			switch c {
			case '|':
				grid[x][y] = VERTICAL
			case '-':
				grid[x][y] = HORIZONTAL
			case '+':
				grid[x][y] = INTERSECTION
			case '/':
				grid[x][y] = FWD_TURN
			case '\\':
				grid[x][y] = BACK_TURN
			case '>':
				grid[x][y] = HORIZONTAL
				cart := new(Cart)
				cart.id = numcarts
				cart.dir = RIGHT
				cart.x = x
				cart.y = y
				cart.lastTurn = TURN_RIGHT
				carts = append(carts, cart)
				numcarts += 1
			case '<':
				grid[x][y] = HORIZONTAL
				cart := new(Cart)
				cart.id = numcarts
				cart.dir = LEFT
				cart.x = x
				cart.y = y
				cart.lastTurn = TURN_RIGHT
				carts = append(carts, cart)
				numcarts += 1
			case 'v':
				grid[x][y] = VERTICAL
				cart := new(Cart)
				cart.id = numcarts
				cart.dir = DOWN
				cart.x = x
				cart.y = y
				cart.lastTurn = TURN_RIGHT
				carts = append(carts, cart)
				numcarts += 1
			case '^':
				grid[x][y] = VERTICAL
				cart := new(Cart)
				cart.id = numcarts
				cart.dir = UP
				cart.x = x
				cart.y = y
				cart.lastTurn = TURN_RIGHT
				carts = append(carts, cart)
				numcarts += 1
			}
		}
		y += 1
	}

//	part1(carts,grid)
	part2(carts,grid)
}

func part1(carts []*Cart, grid [][]int) {
	for true {
		tick(carts, grid)
		isCrash,x,y := findCollision(carts)
		if isCrash {
			fmt.Printf("Crash at %d,%d\n",x,y)
			break
		}
	}
}

func countUncrashedCarts(carts []*Cart) int {
	cnt := 0
	for _, c := range(carts) {
		if ! c.isCrashed {
			cnt += 1
		}
	}
	return cnt
}

func part2(carts []*Cart, grid [][]int) {
	mm := 0
	for true {

		mm += 1
		tick(carts, grid)
		if countUncrashedCarts(carts) == 1 {
			for _, vv := range carts {
				if ! vv.isCrashed {
					fmt.Printf("Last survivor is at %d,%d ID %d\n", vv.x, vv.y, vv.id)
				}
			}

			break
		}	}
}

func removeCart(cart *Cart, carts []*Cart) []*Cart {
	cartpos := -1
	for i, c := range(carts) {
		if cart == c {
			cartpos = i
		}
	}
	if cartpos >= 0 {
		carts = append(carts[:cartpos], carts[cartpos+1:]...)
	}
	return carts
}

func main() {
	process("input.dat")
}

