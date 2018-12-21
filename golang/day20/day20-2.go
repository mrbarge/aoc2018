package main

import (
	"bufio"
	"fmt"
	"os"
)

type Coord struct {
	x int
	y int
}

type Node struct {
	c Coord
}

type Path struct {
	n Node
	length int
}

func process(datafile string) {

	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)
	s.Scan()
	line := s.Text()

	nodes := make(map[Node]int,0)
	nodepaths := make([]Path,0)
	n := Node{Coord{0,0 }}
	length := 0
	for _, v := range line {
		if v == 'N' || v == 'E' || v == 'W' || v == 'S' {
			switch v {
			case 'N':
				n = Node{Coord{n.c.x, n.c.y - 1}}
			case 'E':
				n = Node{Coord{n.c.x-1, n.c.y}}
			case 'W':
				n = Node{Coord{n.c.x+1, n.c.y}}
			case 'S':
				n = Node{Coord{n.c.x, n.c.y + 1}}
			}
			length += 1

			_, ok := nodes[n]
			if !ok {
				nodes[n] = length
			} else {
				if length < nodes[n] {
					nodes[n] = length
				}
			}
		} else if v == '(' {
			nodepaths = append(nodepaths, Path{n, length})
		} else if v == ')' {
			last := nodepaths[len(nodepaths)-1]
			nodepaths = nodepaths[:len(nodepaths)-1]
			n, length = last.n, last.length
		} else if v == '|' {
			last := nodepaths[len(nodepaths)-1]
			n, length = last.n, last.length
		}
	}

	retcnt := 0
	for _, v := range nodes {
		if v >= 1000 {
			retcnt += 1
		}
	}
	fmt.Println(nodes)
	fmt.Printf("%d\n",retcnt)

}

func main() {
	process("input.dat")
	//	process("test.dat")
}

