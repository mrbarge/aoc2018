package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"github.com/albertorestifo/dijkstra"
	"strconv"
)

type Actor struct {
	id int
	hp int
	c Coord
	isGoblin bool
}

type Coord struct {
	x int
	y int
}

type Element rune
const (
	WALL = '#'
	SPACE = '.'
	GOBLIN = 'G'
	ELF = 'E'
)

func process(datafile string) ([][]Element, []*Actor, dijkstra.Graph) {

	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	g := dijkstra.Graph{}

	gmap := make([][]Element, 0)
	actors := make([]*Actor, 0)
	y := 0
	actorId := 0
	for s.Scan() {
		line := s.Text()
		mapline := make([]Element,len(line))

		for x, c := range line {
			scoord := fmt.Sprintf("%d,%d",x,y)
			switch c {
			case '#':
				mapline[x] = WALL
			case '.':
				g[scoord] = make(map[string]int,0)
				mapline[x] = SPACE
			case 'G':
				g[scoord] = make(map[string]int,0)
				actor := new(Actor)
				actor.id = actorId
				actorId += 1
				actor.hp = 200
				actor.c.x = x
				actor.c.y = y
				actor.isGoblin = true
				mapline[x] = GOBLIN
				actors = append(actors, actor)
			case 'E':
				g[scoord] = make(map[string]int,0)
				actor := new(Actor)
				actor.id = actorId
				actorId += 1
				actor.hp = 200
				actor.c.x = x
				actor.c.y = y
				actor.isGoblin = false
				mapline[x] = ELF
				actors = append(actors, actor)
			}
		}
		gmap = append(gmap, mapline)
		y += 1
	}

	// Now add connections to all nodes
	for y, _ := range gmap {
		for x, _ := range gmap {
			if gmap[y][x] != WALL {
				sc := fmt.Sprintf("%d,%d",x,y)
				neighbours := getNeighboursForNode(x,y,gmap)
				for _, n := range neighbours {
					g[sc][n] = 1
					g[n][sc] = 1
				}
			}
		}
	}
	fmt.Println("Got here?")
	printMap(gmap,actors)
	return gmap, actors, g
}

func getNeighboursForNode(x int, y int, gmap [][]Element) []string {
	coords := make([]string, 0)
	if gmap[y][x-1] != WALL {
		coords = append(coords, fmt.Sprintf("%d,%d",x-1,y))
	}
	if gmap[y][x+1] != WALL {
		coords = append(coords, fmt.Sprintf("%d,%d",x+1,y))
	}
	if gmap[y-1][x] != WALL {
		coords = append(coords, fmt.Sprintf("%d,%d",x,y-1))
	}
	if gmap[y+1][x] != WALL {
		coords = append(coords, fmt.Sprintf("%d,%d",x,y+1))
	}
	return coords
}

func targetInRange(e Actor, gmap [][]Element) {
}

func getNeighbourCoords(a *Actor) []Coord {
	return []Coord {
		Coord{a.c.x, a.c.y - 1},
		Coord{a.c.x - 1, a.c.y},
		Coord{a.c.x + 1, a.c.y},
		Coord{a.c.x, a.c.y + 1},
	}
}

func sortActorsByHp(actors []*Actor) []*Actor{
	sort.Slice(actors, func(i, j int) bool {
		return actors[i].hp < actors[j].hp
	})
	return actors
}

func sortActorsByLoc(actors []*Actor) []*Actor{
	sort.Slice(actors, func(i, j int) bool {
		if actors[i].c.y < actors[j].c.y {
			return true
		}
		if actors[i].c.y > actors[j].c.y {
			return false
		}
		return actors[i].c.x < actors[j].c.x
	})
	return actors
}

func sortActorsByLocandHp(actors []*Actor) []*Actor{
	sort.Slice(actors, func(i, j int) bool {
		if actors[i].hp < actors[j].hp {
			return true
		}
		if actors[i].hp > actors[j].hp {
			return false
		}
		if actors[i].c.y < actors[j].c.y {
			return true
		}
		if actors[i].c.y > actors[j].c.y {
			return false
		}
		return actors[i].c.x < actors[j].c.x
	})
	return actors
}

func getNearestNeighbouringTarget(e *Actor, actors []*Actor) *Actor {
	nb := make([]*Actor, 0)
	nc := getNeighbourCoords(e)

	for _, a := range actors {
		if a.hp <= 0 {
			continue
		}
		for _, n := range nc {
			if a.c == n && a.isGoblin != e.isGoblin {
				nb = append(nb, a)
			}
		}
	}
	if len(nb) > 0 {
		return sortActorsByLocandHp(nb)[0]
	} else {
		return nil
	}
}

func findTargets(e *Actor, actors []*Actor) []*Actor {
	targets := make([]*Actor, 0)
	actors = sortActorsByLoc(actors)
	for _, a := range actors {
		if e != a && e.isGoblin != a.isGoblin && a.hp > 0 {
			targets = append(targets, a)
		}
	}
	return targets
}

func attack(attacker *Actor, target *Actor) bool {
	fmt.Printf("Actor at (%d,%d) attacking enemy at (%d,%d)\n",attacker.c.x,attacker.c.y,
		target.c.x,target.c.y)
	target.hp -= 3
	return target.hp <= 0
}

func removeActor(actor *Actor, actors []*Actor) []*Actor {
	for pos := 0; pos < len(actors); pos++ {
		if actors[pos].id == actor.id {
			actors = append(actors[:pos], actors[pos+1:]...)
			break
		}
	}
	return actors
}

func runGame(gmap [][]Element, actors []*Actor, graph dijkstra.Graph) {

	done := false
	counter := 0
	maxCount := 23
	for !done && counter < maxCount {
		// Get active actors in reading list order
		actors = sortActorsByLoc(actors)
		for _, actor := range actors {
			fmt.Printf("Handling actor (%d,%d)\n",actor.c.x,actor.c.y)
			// Get nearest neighbouring targets
			target := getNearestNeighbouringTarget(actor, actors)

			// If there was a nearest neighbour, attack it
			if target != nil {
				isDead := attack(actor, target)

				// Did the attack kill the target, if so remove it
				// from the list of active actors
				if isDead {
					fmt.Printf("Target at (%d,%d) died\n",target.c.x,target.c.y)
					actors = removeActor(target, actors)
					// clear the map
					gmap[target.c.y][target.c.x] = SPACE
				}
			} else {
				// Otherwise find nearest target to move to

				nextCoord := move(actor,actors,gmap,graph)
				if nextCoord.x == actor.c.x && nextCoord.y == actor.c.y {
					fmt.Printf("Actor at (%d,%d) doing nothing\n",actor.c.x,actor.c.y)
					continue
				}
				gmap[actor.c.y][actor.c.x] = SPACE
				if actor.isGoblin {
					gmap[nextCoord.y][nextCoord.x] = GOBLIN
				} else {
					gmap[nextCoord.y][nextCoord.x] = ELF
				}
				actor.c.x = nextCoord.x
				actor.c.y = nextCoord.y

				// and try to attack
				target := getNearestNeighbouringTarget(actor, actors)
				if target != nil {
					isDead := attack(actor, target)
					if isDead {
						fmt.Printf("Target at (%d,%d) died\n", target.c.x, target.c.y)
						actors = removeActor(target, actors)
						// clear the map
						gmap[target.c.y][target.c.x] = SPACE
					}
				}
			}

			// are there no targets left?
			if !remainingTargets(actors, actor.isGoblin) {
				fmt.Println("No more targets")
				done = true
			}
		}
		counter += 1
		printMap(gmap,actors)
	}
}

func move(source *Actor, actors []*Actor, gmap [][]Element, graph dijkstra.Graph) Coord {

	actors = sortActorsByLoc(actors)

	sourcestr := fmt.Sprintf("%d,%d",source.c.x,source.c.y)
	// set the transition arcs out of occupied nodes to be stupidly huge so they're ignored
	for _, actor := range actors {
		ac := fmt.Sprintf("%d,%d",actor.c.x,actor.c.y)
		if actor.id == source.id {
			// don't do it for our source
			continue
		} else {
			for k, _ := range graph[ac] {
				graph[ac][k] = 1000
			}
		}
	}

	// build up a list of preferred targets in order of reading, then distance
	preferredTargets := make([]*Actor,0)
	minCost := 0
	minPath := make([]string, 0)
	for _, actor := range actors {
		if actor.id == source.id {
			// don't do it for our source
			continue
		} else if actor.isGoblin == source.isGoblin {
			// don't kill our own kind
			continue
		}
		ac := fmt.Sprintf("%d,%d",actor.c.x,actor.c.y)
		fmt.Printf("Finding path to %s\n", ac)
		path, cost, err := graph.Path(sourcestr,ac)


		if err != nil {
			fmt.Print("ERror occurred: ")
			fmt.Print(err)
			return source.c
		} else {
			if len(preferredTargets) == 0 {
				minCost = cost
				preferredTargets = append(preferredTargets, actor)
				minPath = path
			} else if cost < minCost {
				preferredTargets = append([]*Actor{actor}, preferredTargets...)
				minCost = cost
				minPath = path
			} else {
				preferredTargets = append(preferredTargets, actor)
			}
		}
	}

	// reset the graph?
	for _, actor := range actors {
		ac := fmt.Sprintf("%d,%d",actor.c.x,actor.c.y)
		if actor.id == source.id {
			// don't do it for our source
			continue
		} else {
			for k, _ := range graph[ac] {
				graph[ac][k] = 1
			}
		}
	}

	retc := keyToCoord(minPath[1])
	fmt.Printf("Actor at %s moving to closest enemy (%s)\n",sourcestr,minPath[1])
	return retc
}

func keyToCoord(k string) Coord {
	r_coord_op, _ := regexp.Compile(`^(\d+),(\d+)$`)
	res_reg_op := r_coord_op.FindStringSubmatch(k)
	if res_reg_op != nil {
		d1, _ := strconv.Atoi(res_reg_op[1])
		d2, _ := strconv.Atoi(res_reg_op[2])
		return Coord{d1,d2}
	} else {
		fmt.Println("BAD BADBAD BADBAD BADBAD BADBAD BADBAD BADBAD BADBAD BAD")
		return Coord{-1,-1}
	}
}

func remainingTargets(actors []*Actor, isGoblin bool) bool {
	for _, a := range actors {
		if a.isGoblin != isGoblin && a.hp > 0 {
			return true
		}
	}
	return false
}

func printMap(gmap [][]Element, actors []*Actor) {
	for i:=0; i<len(gmap); i+=1 {
		for j:=0; j<len(gmap[i]); j+=1 {
			fmt.Print(string(gmap[i][j]))
		}
		fmt.Println("")
	}
	for _, a := range actors {
		if a.isGoblin {
			fmt.Print("G ")
		} else {
			fmt.Print("E ")
		}
		fmt.Printf("(%d,%d) HP: %d\n",a.c.x,a.c.y,a.hp)
	}
}

func main() {
	//	process("input.dat")
	gmap, actors, graph := process("test.dat")
	runGame(gmap, actors, graph)
}

