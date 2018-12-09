package main

import (
	"os"
	"bufio"
	"strings"
	"fmt"
	"strconv"
	"sort"
)

type Node struct {
	value int
	prev *Node
	next *Node

}

func process(datafile string) (int, int) {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)
	s.Scan()
	line := s.Text()

	numPlayers, _ := strconv.Atoi((strings.Split(line, " "))[0])
	points, _ := strconv.Atoi((strings.Split(line, " "))[6])

	return numPlayers, points
}

func game(players int, marbles int) {

	circle := make([]int, 0)
	circle = append(circle, 0)

	// set up starting circle
	currentNode := &Node{0, nil, nil }
	currentNode.next = currentNode
	currentNode.prev = currentNode

	// set up players
	currentPlayer := 1
	score := make([]int, players)

	for m := 1; m < marbles; m++ {
		if m % 23 != 0 {
			n1 := currentNode.next
			n2 := n1.next
			new := &Node{m, nil, nil }
			insertNode(n1, n2, new)

			currentNode = new
		} else {
			// add the current marble to the score
			score[currentPlayer] += m

			// remove pos-7th marble
			points, nextNode := removeFromLeft(currentNode, 7)
			score[currentPlayer] += points
			currentNode = nextNode
		}
		currentPlayer = (currentPlayer+1) % players
	}

	sort.Ints(score)
	fmt.Printf("Highest score is: %d\n",score[len(score)-1])
}

func insertNode(n1 *Node, n2 *Node, v *Node) {
	v.prev = n1
	v.next = n2
	n1.next = v
	n2.prev = v
}

func removeFromLeft(n *Node, pos int) (int,*Node) {
	i := 0
	p := n
	for i < pos {
		p = p.prev
		i += 1
	}
	t := p.next
	(p.prev).next = p.next
	t.prev = p.prev
	return p.value, t
}

func printAll(n *Node) {
	end := n
	done := false
	for n != nil && !done {
		fmt.Printf("%d ", n.value)
		n = n.next
		if n == end {
			done = true
		}
	}
	fmt.Println("")
}

func main() {
	players, marbles := process("input.dat")
	game(players, marbles)
}

