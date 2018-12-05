package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func process(datafile string) string {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)
	s.Scan()
	line := s.Text()

	stack := make([]string,0)

	stack = append(stack,string(line[0]))
	for i := 1; i < len(line); i++ {
		if len(stack) == 0 || ! match(stack[len(stack)-1], string(line[i])) {
			stack = append(stack, string(line[i]))
		} else {
			stack = stack[:len(stack)-1]
		}
	}

	return strings.Join(stack[:], "")
}

func match(a string, b string) bool {
	aIsUpper := (a == strings.ToUpper(a))
	bIsUpper := (b == strings.ToUpper(b))

	// check if they are matching case, in which case they're not a valid match
	if (aIsUpper && bIsUpper) || (!aIsUpper && !bIsUpper) {
		return false
	} else {
		// they're mixed case, so now just compare equality
		return strings.ToLower(a) == strings.ToLower(b)
	}
}

func main() {
	answer := process("input.dat")
	fmt.Println(len(answer))
	fmt.Println(answer)
}
