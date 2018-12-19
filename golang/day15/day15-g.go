package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"github.com/twmb/algoimpl/go/graph"
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

func process(datafile string) ([][]Element, []*Actor) {

	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	g := graph.New(graph.Undirected)
	nodes := make(map[string]graph.Node)
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
				nodes[scoord] = g.MakeNode()
				mapline[x] = SPACE
			case 'G':
			case 'E':
				nodes[scoord] = g.MakeNode()
				actor := new(Actor)
				actor.id = actorId
				actorId += 1
				actor.hp = 300
				actor.c.x = x
				actor.c.y = y
				if c == 'E' {
					actor.isGoblin = false
					mapline[x] = ELF
				} else {
					actor.isGoblin = true
					mapline[x] = GOBLIN
				}
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
					g.MakeEdgeWeight(nodes[sc],nodes[n],1)
				}
			}
		}
	}
	fmt.Println("Got here?")
	return gmap, actors
}

func getNeighboursForNode(x int, y int, gmap [][]Element) []string {
	coords := make([]string, 0)
	if gmap[x-1][y] != WALL {
		coords = append(coords, fmt.Sprintf("%d,%d",x-1,y))
	}
	if gmap[x+1][y] != WALL {
		coords = append(coords, fmt.Sprintf("%d,%d",x+1,y))
	}
	if gmap[x][y-1] != WALL {
		coords = append(coords, fmt.Sprintf("%d,%d",x,y-1))
	}
	if gmap[x][y+1] != WALL {
		coords = append(coords, fmt.Sprintf("%d,%d",x,y+1))
	}
	return coords
}

func targetInRange(e Actor, gmap [][]Element) {
}

func getNeighbourCoords(a *Actor) []Coord {
	return []Coord {
		Coord{a.c.x - 1, a.c.y},
		Coord{a.c.x - 1, a.c.y + 1},
		Coord{a.c.x, a.c.y - 1},
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

func getNearestNeighbouringTarget(e *Actor, actors []*Actor) *Actor {
	nb := make([]*Actor, 0)
	nc := getNeighbourCoords(e)

	for _, a := range actors {
		if a.hp <= 0 {
			continue
		}
		for _, n := range nc {
			if a.c == n {
				nb = append(nb, a)
			}
		}
	}
	if len(nb) > 0 {
		return sortActorsByLoc(sortActorsByHp(nb))[0]
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

func runGame(gmap [][]Element, actors []*Actor) {

	done := false
	for !done {
		// Get active actors in reading list order
		actors = sortActorsByLoc(actors)
		for _, actor := range actors {

			// Get nearest neighbouring targets
			target := getNearestNeighbouringTarget(actor, actors)

			// If there was a nearest neighbour, attack it
			if target != nil {
				isDead := attack(actor, target)

				// Did the attack kill the target, if so remove it
				// from the list of active actors
				if isDead {
					actors = removeActor(target, actors)
				}
			} else {
				// Otherwise find nearest target to move to

			}

			// are there no targets left?
			if remainingTargets(actors, actor.isGoblin) {
				fmt.Println("No more targets")
				done = true
			}
		}
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

func printMap(gmap [][]rune) {
	for i:=0; i<len(gmap); i+=1 {
		for j:=0; j<len(gmap[i]); j+=1 {
			fmt.Print(string(gmap[i][j]))
		}
		fmt.Println("")
	}
}

func main() {
	//	process("input.dat")
	gmap, actors := process("test.dat")
	runGame(gmap, actors)
}

