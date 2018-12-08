package main

import (
	"os"
	"bufio"
	"strings"
	"fmt"
	"strconv"
)


type Node struct {
	id int
	readChildrenCount int
	maxChildrenCount int
	readMetadataCount int
	maxMetadataCount int
	metadataSum int
	parent *Node
	children []*Node
	metadata []int
}

type ReadState int
const (
	READ_NUM_CHILDREN = iota
	READ_NUM_METADATA
	READ_METADATA
)


func process(datafile string) {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)
	s.Scan()
	line := s.Text()

	fmt.Println(line)
	strnodes := strings.Split(line, " ")
	nodes := make([]int,len(strnodes))
	for i, v := range strnodes {
		nodes[i], _ = strconv.Atoi(v)
	}

	nodestack := make([]*Node, 0)
	allnodes := make(map[int]*Node,0)

	state := READ_NUM_CHILDREN
	nodeId := 0
	for i := 0; i < len(nodes); i++ {

		switch state {

		case READ_NUM_CHILDREN:
			var parent *Node = nil
			if len(nodestack) > 0 {
				parent = nodestack[len(nodestack)-1]
				fmt.Printf("Reading next node: %d at pos: %d (parent: %d, children: %d)\n",nodeId, i, parent.id, nodes[i])
			} else {
				fmt.Printf("Reading next node: %d at pos: %d (parent: %d, children: %d)\n",nodeId, i, -1, nodes[i])
			}
			n := Node{nodeId,0, nodes[i], -1, 0, 0, parent, nil, nil}
			n.children = make([]*Node,0)
			n.metadata = make([]int, 0)
			if (n.parent != nil) {
				(allnodes[parent.id]).children = append((allnodes[parent.id]).children, &n)
			}
			nodestack = append(nodestack, &n)
			allnodes[nodeId] = &n
			state = READ_NUM_METADATA
			nodeId += 1

		case READ_NUM_METADATA:
			fmt.Printf("Reading max metadata %d at pos: %d\n",nodes[i], i)
			n := nodestack[len(nodestack)-1]
			n.readMetadataCount = 0
			n.maxMetadataCount = nodes[i]

			if (n.readChildrenCount < n.maxChildrenCount) {
				state = READ_NUM_CHILDREN
				n.readChildrenCount += 1
			} else {
				state = READ_METADATA
			}

		case READ_METADATA:
			n := nodestack[len(nodestack)-1]
			fmt.Printf("Reading metadata %d at pos: %d for node %d\n",nodes[i],i,n.id)

			if n.maxChildrenCount == 0 {
				n.metadataSum += nodes[i]
			} else {
				n.metadata = append(n.metadata,nodes[i]-1)
			}
			n.readMetadataCount += 1

			if n.readMetadataCount < n.maxMetadataCount {
				state = READ_METADATA
			} else {
				nodestack = nodestack[:len(nodestack)-1]
				if len(nodestack) > 0 {
					n := nodestack[len(nodestack)-1]
					if n.readChildrenCount < n.maxChildrenCount {
						state = READ_NUM_CHILDREN
						n.readChildrenCount += 1
					} else if n.readMetadataCount < n.maxMetadataCount {
						state = READ_METADATA
					}
				}
			}
		}
	}

	rootNode := allnodes[0]
	metadataSum := getSum(rootNode)
	fmt.Printf("Metadata sum: %d\n",metadataSum)
}

func getSum(n *Node) int {
	fmt.Printf("Summing node: %d (num metadata: %d)\n", n.id, len(n.metadata))
	if n.maxChildrenCount == 0 {
		fmt.Printf("Summing node %d returns %d\n", n.id, n.metadataSum)
		return n.metadataSum
	} else {
		ms := 0
		for i:=0; i < len(n.metadata); i++ {
			if (n.metadata[i] >= n.maxChildrenCount) {
				fmt.Printf("Ignoring sum of child %d\n",n.metadata[i])
				ms += 0
			} else {
				fmt.Printf("Going to sum child %d\n", n.metadata[i])
				ms += getSum(n.children[n.metadata[i]])
			}
		}
		return ms
	}
}

func main() {
	process("input.dat")
}

