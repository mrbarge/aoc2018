package main

import (
	"bufio"
	"fmt"
	"os"
)

func process(datafile string) {

	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)
	s.Scan()
	line := s.Text()

	length, path := handle(line)
	fmt.Printf("PAth length: %d  PAth: %s\n",length,path)
}

func handle(s string) (int, string) {
	fmt.Printf("Handling new batch: %s\n", s)

	splits := getSplits(s)

	if len(splits) == 0 {
		return 0, ""
	} else if len(splits) > 1 {
		longest := 0
		var longestpath string
		for _, split := range splits {
			l, _ := handle(split)
			if l > longest {
				longest = l
				longestpath = split
			}
		}
		return longest, longestpath
	} else {
		length := 0
		path := ""
		for i := 0; i < len(s); i++ {
			if s[i] == '(' {
				endpos := findLastBracket(s,i)
				if endpos < 0 {
					fmt.Println("Bad input")
					break
				}
				l2, p2 := handle(s[i+1:endpos])
				length += l2
				path += p2
				i = endpos
			} else {
				length += 1
				path += string(s[i])
			}
		}
		return length, path
	}
}

func getSplits(s string) []string {
	bc := 0
	splits := make([]string,0)

	// just ignore this if it is an 'optional' path
	if s[len(s)-1] == '|' {
		return splits
	}

	lastidx := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '|' && bc == 0 {
			splits = append(splits,s[lastidx:i])
			lastidx = i+1
		} else if s[i] == '(' {
			bc += 1
		} else if s[i] == ')' {
			bc -= 1
		}
	}
	splits = append(splits,s[lastidx:])
	return splits
}

func findLastBracket(s string, pos int) int {
	bc := 0
	for i := pos; i < len(s); i++ {
		if s[i] == '(' {
			bc += 1
		} else if s[i] == ')' {
			if bc == 1 {
				return i
			} else {
				bc -= 1
			}
		}
	}
	return -1
}

func main() {
	process("input.dat")
//	process("test.dat")
}

