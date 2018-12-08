package main

import (
	"os"
	"bufio"
	"strings"
	"fmt"
	"strconv"
)


type Node struct {
	readChildrenCount int
	maxChildrenCount int
	readMetadataCount int
	maxMetadataCount int
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

	state := READ_NUM_CHILDREN
	nodeId := 0
	metadataSum := 0
	for i := 0; i < len(nodes); i++ {

		switch state {

		case READ_NUM_CHILDREN:
			fmt.Printf("Reading next node: %d at pos: %d (children: %d)\n",nodeId, i, nodes[i])
			n := Node{0, nodes[i], -1, 0}
			nodestack = append(nodestack, &n)
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
			fmt.Printf("Reading metadata %d at pos: %d\n",nodes[i],i)
			metadataSum += nodes[i]
			n := nodestack[len(nodestack)-1]
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

	fmt.Printf("Metadata sum: %d\n",metadataSum)
}

func main() {
	process("test.dat")
}

