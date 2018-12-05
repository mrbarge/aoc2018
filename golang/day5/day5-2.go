package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func process(datafile string) int {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)
	s.Scan()
	origline := s.Text()

	bestLen := len(origline)
	for _, chr := range "abcdefghijklmnopqrstuvwxyz" {
		line := strings.Replace(origline,strings.ToLower(string(chr)),"",-1)
		line = strings.Replace(line,strings.ToUpper(string(chr)),"",-1)

		stack := make([]string,0)
		stack = append(stack,string(line[0]))
		for i := 1; i < len(line); i++ {
			if len(stack) == 0 ||
				(! match(stack[len(stack)-1], string(line[i]))) {
				stack = append(stack, string(line[i]))
			} else {
				stack = stack[:len(stack)-1]
			}
		}
		if len(stack) < bestLen {
			bestLen = len(stack)
		}
	}

	return bestLen
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
	fmt.Println(answer)
}
