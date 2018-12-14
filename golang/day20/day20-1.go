package main

import (
	"bufio"
	"os"
	"regexp"
)

func process(datafile string) {

	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	r_stub, _ := regexp.Compile(`^stub$`)

	for s.Scan() {
		line := s.Text()
		res_stub := r_stub.FindStringSubmatch(line)

	}

}


func main() {
//	process("input.dat")
	process("test.dat")
}

