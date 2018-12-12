package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type Node struct {
	id int
	plant bool
	prev *Node
	next *Node
}

type Rule struct {
	state []bool
	plant bool
}

func process(datafile string, generations int) {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	r_initial, _ := regexp.Compile(`^initial state: (.+)$`)
	r_rule, _ := regexp.Compile(`^(.+) => (.)$`)

	rules := make([]*Rule,0)
	var node *Node
	for s.Scan() {
		line := s.Text()

		res_initial := r_initial.FindStringSubmatch(line)
		res_rule := r_rule.FindStringSubmatch(line)

		if res_initial != nil {
			node = readInitialState(res_initial[1])
		} else if res_rule != nil {
			r := readRule(res_rule[1], res_rule[2])
			rules = append(rules, r)
		}
	}

	for i := 0; i < generations; i++ {
		node = generation2(node, rules)
	}
	fmt.Println(countPlants(node))
}

func process2(datafile string, generations int) {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	r_initial, _ := regexp.Compile(`^initial state: (.+)$`)
	r_rule, _ := regexp.Compile(`^(.+) => (.)$`)

	rules := make([]*Rule,0)
	var node *Node
	for s.Scan() {
		line := s.Text()

		res_initial := r_initial.FindStringSubmatch(line)
		res_rule := r_rule.FindStringSubmatch(line)

		if res_initial != nil {
			node = readInitialState(res_initial[1])
		} else if res_rule != nil {
			r := readRule(res_rule[1], res_rule[2])
			rules = append(rules, r)
		}
	}

	last := 0
	lastDiff := 0
	total := 0
	for i := 0; i < generations; i++ {
		node = generation2(node, rules)
		n := countPlants(node)
		diff := n-last
		//fmt.Printf("%d %d (%d)\n",i,n,diff)

		// once distance has stabilized, we can calculate the rest
		if diff == lastDiff {
			total = (generations-i-1)*diff + n
			break
		}

		last = n
		lastDiff = diff
	}
	fmt.Println(total)
}

func countPlants(n *Node) int {
	c := 0
	for n != nil {
		if n.plant {
			//fmt.Printf("Match for %d\n",n.id)
			c += n.id
		}
		n = n.next
	}
	return c
}

func generation2(x *Node, rules []*Rule) *Node {

	// our jumping off point
	currNode := x

	// add two new non-plants to the left, then make that our jump-off point
	// but only if there's a plant node in slots 3,4,5..


	if x.plant || x.next.plant || x.next.next.plant {
		n1 := new(Node)
		n2 := new(Node)

		n1.plant = false
		n1.next = x
		n1.prev = n2
		n1.id = x.id - 1

		x.prev = n1

		n2.plant = false
		n2.next = n1
		n2.prev = nil
		n2.id = x.id - 2

		// adjust our jumping off point
		currNode = n2
	}

	// and two on the right..
	// again, only if slots 0,1,2 have a plant
	endNode := x
	for endNode.next != nil {
		endNode = endNode.next
	}
	if endNode.plant || endNode.prev.plant || endNode.prev.prev.plant {
		o1 := new(Node)
		o2 := new(Node)
		o1.prev = endNode
		o1.next = o2
		o1.id = endNode.id + 1
		endNode.next = o1
		o2.id = o1.id + 1
		o2.prev = o1
		o2.next = nil
	}

	//fmt.Println("Start of generation")
	//printState(currNode)
	currRet := new(Node)
	currRet.id = currNode.id
	currRet.next = currNode.next
	currRet.prev = currNode.prev

	startRet := currRet

	// then GENERATE!!!
	for currNode != nil {
		//fmt.Printf("Index %d State=%t\n",i,currNode.plant)
		if isPlantNext(currNode, rules) {
			currRet.plant = true
		} else {
			currRet.plant = false
		}
		tmp := new(Node)
		tmp.plant = false
		tmp.prev = currRet
		tmp.next = nil
		tmp.id = currRet.id + 1
		currRet.next = tmp
		currRet = tmp

		// iterate to next in the list
		currNode = currNode.next
	}

	return startRet
}

func isPlantNext(n *Node, rules []*Rule) bool {
	if n.prev == nil || n.next == nil {
		return false
	}
	for _, rule := range rules {
		// first don't do anything thats on an edge
		if n.plant != rule.state[2] {
			continue
		}
		if n.prev != nil && n.prev.plant != rule.state[1] {
			continue
		}
		if n.prev != nil && n.prev.prev != nil && n.prev.prev.plant != rule.state[0] {
			continue
		}
		if n.next != nil && n.next.plant != rule.state[3] {
			continue
		}
		if n.next != nil && n.next.next != nil && n.next.next.plant != rule.state[4] {
			continue
		}
		/*
		fmt.Print("This ")
		fmt.Print(rule.state)
		fmt.Print( " matched\n")*/
		return rule.plant
	}
	return false
}

func printState(zeroNode *Node) {
	// Find beginning node
	begin := zeroNode
	for begin.prev != nil {
		begin = begin.prev
	}

	for begin != nil {
		if begin.plant {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}
		begin = begin.next
	}
	fmt.Println("")
}

func readRule(rule string, plantstate string) *Rule {
	r := new(Rule)
	if plantstate == "." {
		r.plant = false
	} else {
		r.plant = true
	}

	r.state = make([]bool, len(rule))
	for i, v := range(rule) {
		if string(v) == "." {
			r.state[i] = false
		} else {
			r.state[i] = true
		}
	}
	return r
}

func readInitialState(state string) *Node {

	var lastNode *Node
	currentNode := new(Node)
	zeroNode := currentNode
	for i, v := range state {
		if string(v) == "." {
			currentNode.plant = false
		} else {
			currentNode.plant = true
		}
		currentNode.id = i
		lastNode = currentNode
		currentNode = new(Node)
		lastNode.next = currentNode
		currentNode.prev = lastNode
	}
	return zeroNode
}

func main() {
	process("input.dat", 20)
	process2("input.dat", 50000000000)
}

